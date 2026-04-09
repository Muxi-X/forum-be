package dao

import (
	"fmt"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TagContentKey = "tag:content:%s" // tag:content:xxx -> tagID
	TagIDKey      = "tag:id:%s"      // tag:id:xxx -> tagContent

	TagCacheTTL = 10 * 24 * time.Hour
)

// todo content 字段加唯一索引

type TagModel struct {
	Id      uint32
	Content string
}

func (TagModel) TableName() string {
	return "tags"
}

// Create ...
func (t *TagModel) Create() error {
	return dao.DB.Create(t).Error
}

func (d *Dao) GetTagById(id uint32) (*TagModel, error) {
	tag := &TagModel{
		Id: id,
	}
	content, err := d.getTagContentById(strconv.Itoa(int(id)))
	if err != nil {
		return tag, err
	}
	if content != "" {
		tag.Content = content
		return tag, nil
	}

	// 从redis缓存中未命中则在数据库找
	if err := dao.DB.Model(tag).First(tag).Error; err != nil {
		return nil, err
	}

	if err := dao.addTag(tag.Id, tag.Content); err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrRedis))
	}

	return tag, nil
}

func (d *Dao) GetTagByContent(content string) (*TagModel, error) {
	tag := &TagModel{
		Content: content,
	}

	if content == "" {
		return tag, nil
	}

	id, err := d.getTagIdByContent(content)
	if err != nil {
		return tag, err
	}
	if id != 0 {
		tag.Id = uint32(id)
		return tag, nil
	}

	// 从redis缓存中未命中则在数据库找
	err = dao.DB.Model(tag).Where("content = ?", content).First(tag).Error
	if err == gorm.ErrRecordNotFound {
		// 在数据库未找到则新建
		err = tag.Create()
	}

	if err := dao.addTag(tag.Id, tag.Content); err != nil {
		logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
	}

	return tag, err
}

func (d *Dao) getTagContentById(id string) (string, error) {
	content, err := d.Redis.Get("tag:id:" + id).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	err = d.Redis.Expire("tag:id:"+id, 10*24*time.Hour).Err()
	return content, err
}

func (d *Dao) getTagIdByContent(content string) (int, error) {
	id, err := d.Redis.Get("tag:content:" + content).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	err = d.Redis.Expire("tag:content:"+content, 10*24*time.Hour).Err()
	return id, err
}

func (d *Dao) addTag(id uint32, content string) error {
	pipe := d.Redis.TxPipeline()

	pipe.Set("tag:id:"+strconv.Itoa(int(id)), content, 10*24*time.Hour)
	pipe.Set("tag:content:"+content, id, 10*24*time.Hour)

	_, err := pipe.Exec()
	return err
}

// BatchGetOrCreateTags
// NOTE:
// 1. redis miss -> mysql query -> mysql miss -> mysql insert -> mysql query(防止 toctou) -> redis set
// 可以进一步优化为：
// 2. redis miss -> mysql insert -> mysql query -> redis set
// 但是第二种方案在后期——如果tag基本mysql中都有的话依然，会有两次db操作，而第一种方案在tag基本mysql中都有的情况下就只有一次db操作了
// 交给后人优化了
func (d *Dao) BatchGetOrCreateTags(tags []string) ([]*TagModel, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	// 批量 redis 查询
	redisHit, redisMiss, err := d.batchGetTagIDsFromCache(tags)
	if err != nil {
		return nil, err
	}

	// 未命中部分 批量查 Mysql
	dbHit, dbMiss, err := d.batchGetTagsFromDB(redisMiss)
	if err != nil {
		return nil, err
	}

	// 将未命中部分插入数据库
	dbTags, err := d.batchInsertTagsIgnoreConflict(dbMiss)
	if err != nil {
		return nil, err
	}

	// 将数据库命中部分和新建部分的 tagID 插入 redis
	var tagsToAdd []*TagModel
	tagsToAdd = append(tagsToAdd, dbHit...)
	tagsToAdd = append(tagsToAdd, dbTags...)

	err = d.batchSetTagsToCache(tagsToAdd)
	if err != nil {
		return nil, err
	}

	// 确保返回顺序一致
	resMap := make(map[string]*TagModel, len(tags))

	// redis hit
	for content, id := range redisHit {
		resMap[content] = &TagModel{
			Id:      id,
			Content: content,
		}
	}

	// db hit
	for _, tag := range dbHit {
		resMap[tag.Content] = tag
	}

	// db insert
	for _, tag := range dbTags {
		resMap[tag.Content] = tag
	}

	res := make([]*TagModel, len(tags))
	for i, content := range tags {
		if tag, ok := resMap[content]; ok {
			res[i] = tag
		}
	}
	return res, nil
}

func (d *Dao) batchGetTagIDsFromCache(contents []string) (map[string]uint32, []string, error) {
	hit := make(map[string]uint32)
	miss := make([]string, 0)
	if len(contents) == 0 {
		return hit, miss, nil
	}

	keys := make([]string, len(contents))
	for i, content := range contents {
		keys[i] = buildTagContentKey(content)
	}

	// 批量查询
	val, err := d.Redis.MGet(keys...).Result()
	if err != nil {
		return nil, nil, err
	}

	pipe := d.Redis.TxPipeline()

	for i, v := range val {
		if v == nil {
			miss = append(miss, contents[i])
			continue
		}

		vStr, ok := v.(string)
		if !ok {
			// 这里转换失败是数据异常了，记录日志后认为未命中
			logger.Error("Failed to assert redis value as string", zap.Any("value", v))
			miss = append(miss, contents[i])
			continue
		}
		id, err := strconv.Atoi(vStr)
		if err != nil {
			// 同上
			logger.Error("Failed to convert tag ID to int", zap.String("value", vStr), zap.Error(err))
			miss = append(miss, contents[i])
			continue
		}

		hit[contents[i]] = uint32(id)

		// 刷新 TTL
		pipe.Expire(keys[i], TagCacheTTL)
	}

	_, err = pipe.Exec()
	return hit, miss, err
}

func (d *Dao) batchGetTagsFromDB(contents []string) ([]*TagModel, []string, error) {
	var tags []*TagModel
	miss := make([]string, 0)
	if len(contents) == 0 {
		return tags, miss, nil
	}

	err := d.DB.Where("content IN ?", contents).Find(&tags).Error
	if err != nil {
		return nil, nil, err
	}

	// 命中
	hit := make(map[string]*TagModel, len(tags))
	for _, tag := range tags {
		hit[tag.Content] = tag
	}

	// 未命中
	for _, content := range contents {
		if _, ok := hit[content]; !ok {
			miss = append(miss, content)
		}
	}

	return tags, miss, nil
}

func (d *Dao) batchInsertTagsIgnoreConflict(contents []string) ([]*TagModel, error) {
	tags := make([]*TagModel, 0, len(contents))
	if len(contents) == 0 {
		return tags, nil
	}
	for _, content := range contents {
		tags = append(tags, &TagModel{Content: content})
	}

	if len(tags) == 0 {
		return tags, nil
	}

	// 如果有唯一索引
	err := d.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "content"}},
		DoNothing: true,
	}).Create(&tags).Error

	// 重新查询，确保拿到 ID
	var res []*TagModel
	err = d.DB.Where("content IN ?", contents).Find(&res).Error
	return res, err
}

func (d *Dao) batchSetTagsToCache(tags []*TagModel) error {
	if len(tags) == 0 {
		return nil
	}
	pipe := d.Redis.TxPipeline()

	for _, tag := range tags {
		pipe.Set(buildTagIDKey(tag.Id), tag.Content, TagCacheTTL)
		pipe.Set(buildTagContentKey(tag.Content), tag.Id, TagCacheTTL)
	}

	_, err := pipe.Exec()
	return err
}

func buildTagContentKey(content string) string {
	return fmt.Sprintf(TagContentKey, content)
}
func buildTagIDKey(id uint32) string {
	return fmt.Sprintf(TagIDKey, strconv.Itoa(int(id)))
}

type Post2TagModel struct {
	Id     uint32
	PostId uint32
	TagId  uint32
}

func (Post2TagModel) TableName() string {
	return "post2tags"
}

// Create ...
func (p *Post2TagModel) Create() error {
	return dao.DB.Create(p).Error
}

func (Dao) CreatePost2Tag(item Post2TagModel) error {
	return item.Create()
}

func (d *Dao) ListTagsByPostId(postId uint32) ([]string, []uint32, error) {
	var tagIds []uint32
	err := d.DB.Table("post2tags").Where("post_id = ?", postId).Pluck("tag_id", &tagIds).Error
	if err != nil {
		return nil, nil, err
	}

	contents := make([]string, len(tagIds))
	for i, item := range tagIds {
		tag, err := d.GetTagById(item)
		if err != nil {
			return nil, nil, err
		}
		contents[i] = tag.Content
	}

	return contents, tagIds, nil
}

func (d *Dao) AddTagToSortedSet(tagId uint32, category string) error {
	pipe := d.Redis.TxPipeline()

	pipe.ZIncrBy("tags:", 1, strconv.Itoa(int(tagId)))
	pipe.ZIncrBy("tags:"+category, 1, strconv.Itoa(int(tagId)))

	_, err := pipe.Exec()
	return err
}

func (d *Dao) BatchAddTagsToSortedSet(tagIDs []uint32, category string) error {
	pipe := d.Redis.TxPipeline()

	for _, tagID := range tagIDs {
		member := strconv.Itoa(int(tagID))
		pipe.ZIncrBy("tags:", 1, member)
		pipe.ZIncrBy("tags:"+category, 1, member)
	}

	_, err := pipe.Exec()
	return err
}

func (d *Dao) ListPopularTags(category string) ([]string, error) {
	// 降序
	ids, err := d.Redis.ZRevRange("tags:"+category, 0, 9).Result()
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(ids))
	for i, id := range ids {
		tags[i], err = d.getTagContentById(id)
		if err != nil {
			return nil, err
		}
	}

	return tags, nil
}

func (d Dao) deletePost2TagByPostId(postId uint32, tx ...*gorm.DB) error {
	db := d.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Table("post2tags").Where("post_id = ?", postId).Delete(&Post2TagModel{}).Error
}

func (d Dao) isExistPostWithTagId(tagId int) (bool, error) {
	var count int64
	if err := d.DB.Table("post2tags").Where("tag_id = ?", tagId).Count(&count).Error; err != nil {
		return false, err
	}

	return count != 0, nil
}

func (d Dao) isExistPostWithTagIdAndCategory(tagId int, category string) (bool, error) {
	var count int64
	if err := d.DB.Table("post2tags").Joins("join posts p on p.id = post2tags.post_id").Where("tag_id = ? AND category = ?", tagId, category).Count(&count).Error; err != nil {
		return false, err
	}

	return count != 0, nil
}

type SipScoreTagModel struct {
	ID         uint32 `gorm:"primaryKey"`
	SipScoreID uint32 `gorm:"uniqueIndex:idx_rank_tag;index"`
	TagID      uint32 `gorm:"uniqueIndex:idx_rank_tag;index"`
}

func (SipScoreTagModel) TableName() string {
	return "sip_score_tags"
}

// Create ...
func (s *SipScoreTagModel) Create() error {
	return dao.DB.Create(s).Error
}

func (d *Dao) CreateSipScoreTag(item *SipScoreTagModel) error {
	return item.Create()
}

func (d *Dao) BatchCreateSipScoreTags(items []*SipScoreTagModel) error {
	return d.DB.Create(&items).Error
}
