package dao

import (
	pb "forum-post/proto"
	"forum/pkg/constvar"
	"strconv"

	"gorm.io/gorm"
)

type PostModel struct {
	Id              uint32
	Domain          string
	Content         string
	Title           string
	CreateTime      string
	Category        string
	Re              bool
	CreatorId       uint32
	LastEditTime    string
	LikeNum         uint32
	ContentType     string
	CompiledContent string
	Summary         string
	Score           uint32
	IsReport        bool
	Quality         bool
}

func (PostModel) TableName() string {
	return "posts"
}

// Create ...
func (p *PostModel) Create() error {
	return dao.DB.Create(p).Error
}

// Save ...
func (p *PostModel) Save(tx ...*gorm.DB) error {
	db := dao.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Save(p).Error
}

func (p *PostModel) Update() error {
	err := dao.DB.Table("posts").Where("id = ?", p.Id).Updates(map[string]interface{}{
		"title":            p.Title,
		"content":          p.Content,
		"compiled_content": p.CompiledContent,
		"last_edit_time":   p.LastEditTime,
		"domain":           p.Domain,
		"category":         p.Category,
		"summary":          p.Summary,
	}).Error

	return err
}

func (p *PostModel) Delete(tx *gorm.DB) error {
	p.Re = true
	return p.Save(tx)
}

func (p *PostModel) Get(id uint32) error {
	return dao.DB.Model(p).Where("id = ? AND re = 0", id).First(&p).Error
}

func (p *PostModel) BeReported() error {
	p.IsReport = true
	return p.Save()
}

type PostInfo struct {
	Id              uint32
	Content         string
	Title           string
	Category        string
	CreatorId       uint32
	LastEditTime    string
	CreatorName     string
	CreatorAvatar   string
	CommentNum      uint32
	LikeNum         uint32
	ContentType     string
	CompiledContent string
	Summary         string
	Domain          string
}

func (Dao) CreatePost(post *PostModel) (uint32, error) {
	err := post.Create()
	return post.Id, err
}

func (d *Dao) ListMainPost(filter *PostModel, typeName string, offset, limit, lastId uint32, pagination bool, searchContent string, tagId uint32) ([]*PostInfo, error) {
	var posts []*PostInfo
	query := d.DB.Table("posts").Select("posts.id id, title, category, compiled_content, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, content_type, summary, domain").Joins("join users u on u.id = posts.creator_id").Where(filter).Where("posts.re = 0 AND posts.is_report = 0")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("posts.id < ?", lastId)
		}
	}

	if tagId != 0 {
		var postIds []uint32
		if err := d.DB.Table("post2tags").Select("post_id").Distinct("post_id").Where("tag_id = ?", tagId).Find(&postIds).Error; err != nil {
			return nil, err
		}

		query = query.Where("posts.id IN ?", postIds)
	}

	if searchContent != "" {
		// query = query.Where("MATCH (content, title) AGAINST (?)", searchContent) // MySQL 5.7.6 才支持中文全文索引
		key := "%" + searchContent + "%"
		query = query.Where("posts.content LIKE ? OR posts.title LIKE ? OR posts.summary LIKE ?", key, key, key)
	}

	if typeName == "hot" {
		query = query.Order("posts.score DESC")
	} else if typeName == "quality" {
		query = query.Where("quality = ?", true).Order("posts.id DESC")
	} else {
		query = query.Order("posts.id DESC")
	}

	err := query.Scan(&posts).Error

	return posts, err
}

func (d *Dao) ListUserCreatedPost(creatorId uint32) ([]uint32, error) {
	var postIds []uint32
	err := d.DB.Table("posts").Select("id").Where("creator_id = ? AND re = 0", creatorId).Find(&postIds).Error

	return postIds, err
}

func (d *Dao) GetPostInfo(postId uint32) (*PostInfo, error) {
	var post PostInfo
	err := d.DB.Table("posts").Select("posts.id id, title, category, compiled_content, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num, content_type, summary").Joins("join users u on u.id = posts.creator_id").Where("posts.id = ? AND posts.re = 0", postId).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

func (Dao) GetPost(id uint32) (*PostModel, error) {
	var item PostModel
	err := item.Get(id)
	return &item, err
}

func (d Dao) ChangePostScore(postId uint32, score int) error {
	return d.Redis.ZIncrBy("posts:", float64(score), strconv.Itoa(int(postId))).Err()
}

func (d Dao) AddChangeRecord(postId uint32) error {
	return d.Redis.SAdd("changed_posts", strconv.Itoa(int(postId))).Err()
}

// func (d Dao) ListHotPost(domain, category string, offset, limit uint32, pagination bool) ([]*PostInfo, error) {
// 	key := "hot:" + domain
// 	if category != "" {
// 		key += ":" + category
// 	}
//
// 	var start int64
// 	var end int64 = -1
// 	if pagination {
// 		if limit == 0 {
// 			limit = constvar.DefaultLimit
// 		}
//
// 		start = int64(offset)
// 		end = int64(offset + limit)
// 	}
//
// 	result, err := d.Redis.ZRevRange(key, start, end).Result() // 降序
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	list := make([]*PostInfo, len(result))
// 	for i, r := range result {
// 		data, err := strconv.Atoi(r)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		list[i], err = d.GetPostInfo(uint32(data))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return list, nil
// }

func (d Dao) ListPostInfoByPostIds(postIds []uint32, filter *PostModel, offset, limit, lastId uint32, pagination bool) ([]*pb.PostPartInfo, error) {
	var posts []*pb.PostPartInfo
	query := d.DB.Table("posts").Select("posts.id id, title, category, summary, content, last_edit_time time, creator_id, u.name creator_name, u.avatar creator_avatar, content_type").Joins("join users u on u.id = posts.creator_id").Where("posts.re = 0").Where("posts.id IN ?", postIds).Order("posts.id DESC")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("posts.id < ?", lastId)
		}
	}

	if err := query.Scan(&posts).Error; err != nil {
		return nil, err
	}

	for _, post := range posts {
		likeNum, err := d.GetLikedNum(Item{
			Id:       post.Id,
			TypeName: constvar.Post,
		})
		if err != nil {
			return nil, err
		}
		post.LikeNum = uint32(likeNum)

		post.CommentNum, err = d.GetCommentNumByPostId(post.Id)
		if err != nil {
			return nil, err
		}

		post.CollectionNum, err = d.GetCollectionNumByPostId(post.Id)
		if err != nil {
			return nil, err
		}

		post.Tags, err = d.ListTagsByPostId(post.Id)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (d Dao) syncPostScore() error {
	result, err := d.Redis.ZRevRangeWithScores("posts:", 0, -1).Result()
	if err != nil {
		return err
	}

	for _, r := range result {
		err := d.DB.Table("posts").Where("id = ?", r.Member).Update("score", r.Score).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (d Dao) syncItemLike() error {
	postIds, err := d.Redis.SMembers("changed_posts").Result()
	if err != nil {
		return err
	}

	for _, id := range postIds {
		num, err := d.Redis.SCard("like:" + constvar.Post + "_list:" + id).Result()
		if err != nil {
			return err
		}

		if err := d.DB.Table("posts").Where("id = ?", id).Update("like_num", num).Error; err != nil {
			return err
		}
	}

	commentsIds, err := d.Redis.SMembers("changed_comments").Result()
	if err != nil {
		return err
	}

	for _, id := range commentsIds {
		num, err := d.Redis.SCard("like:" + constvar.Comment + "_list:" + id).Result()
		if err != nil {
			return err
		}

		if err := d.DB.Table("comments").Where("id = ?", id).Update("like_num", num).Error; err != nil {
			return err
		}
	}

	if err := d.Redis.SRem("changed_posts").Err(); err != nil {
		return err
	}
	if err := d.Redis.SRem("changed_comments").Err(); err != nil {
		return err
	}

	return nil
}

func (d Dao) ChangeQualityPost(postId uint32, quality bool) error {
	return dao.DB.Table("posts").Where("id = ?", postId).Update("quality", quality).Error
}

type LastClickModel struct {
	UserId      uint32
	Category    string
	LastClickAt string
}

func (c *LastClickModel) GetOrCreateLastRead() error {
	return dao.DB.Table("last_clicks").FirstOrCreate(c, LastClickModel{UserId: c.UserId, Category: c.Category}).Error
}

func (d Dao) UpdateLastRead(userId uint32, category, time string) error {
	return d.DB.Table("last_clicks").Where("user_id = ? AND category = ?", userId, category).Update("last_click_at", time).Error
}

func (d Dao) CountPostByTime(time, category string) (int64, error) {
	var count int64
	err := d.DB.Table("posts").Where("create_time >= ? AND re = 0 AND is_report = 0 AND category = ?", time, category).Count(&count).Error
	return count, err
}
