package dao

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// 茶评 榜单 对象

const (
	orderDirDesc = "DESC"
	orderDirAsc  = "ASC"
)

// todo 创建索引
// todo 1. tag - content 二级索引
// todo 2.

type SipScoreModel struct {
	ID             uint32    `gorm:"primarykey;index:idx_rank,priority:3;index:idx_latest,priority:2;index:idx_creator,priority:2"`
	CreatedAt      time.Time `gorm:"index:idx_latest,priority:1"`
	UpdatedAt      time.Time
	DeletedAt      soft_delete.DeletedAt `gorm:"index;softDelete:nano"`
	LastModifiedBy uint32
	CreatorID      uint32 `gorm:"index:idx_creator,priority:1"`
	EntryCount     uint32 `gorm:"type:int unsigned;default:0"`

	// 冗余字段，用于排序
	CollectCount     uint32 `gorm:"type:int unsigned;default:0;index:idx_rank,priority:1"`
	ParticipantCount uint32 `gorm:"type:int unsigned;default:0;index:idx_rank,priority:2"`

	// 是否被举报过多而被 ban 了
	IsReport bool `gorm:"index"`

	// 用户可编辑的字段
	Name        string `gorm:"type:varchar(100);not null;index:,class:FULLTEXT,option:WITH PARSER ngram"`
	Description string `gorm:"type:varchar(500)"`
	CoverImg    string `gorm:"type:varchar(255)"`
	Domain      string `gorm:"type:varchar(20);not null;index"`
	Category    string `gorm:"type:varchar(20);not null;index"`
}

func (SipScoreModel) TableName() string {
	return "sip_scores"
}

// Create ...
func (s *SipScoreModel) Create() error {
	return dao.DB.Create(s).Error
}

// Save ...
func (s *SipScoreModel) Save(tx ...*gorm.DB) error {
	db := dao.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Save(s).Error
}

func (s *SipScoreModel) Delete(tx *gorm.DB) error {
	return tx.Delete(s).Error
}

func (s *SipScoreModel) Get(id uint32) error {
	return dao.DB.First(s, id).Error
}

func (s *SipScoreModel) BeReported() error {
	s.IsReport = true
	return s.Save()
}

type SipScoreInfo struct {
	ID        uint32
	CreatedAt string
	UpdatedAt string

	CreatorID uint32

	Name        string
	Description string
	CoverImg    string

	CollectCount     uint32
	ParticipantCount uint32
}

func (d *Dao) CreateSipScore(sipScore *SipScoreModel) (uint32, error) {
	err := sipScore.Create()
	return sipScore.ID, err
}

func (d *Dao) UpdateSipScore(id uint32, update map[string]interface{}, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	result := db.Model(&SipScoreModel{}).Where("id = ?", id).Updates(update)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (d *Dao) GetSipScore(id uint32, tx ...*gorm.DB) (*SipScoreModel, error) {
	db := d.getDB(tx...)
	var sipScore SipScoreModel
	err := db.First(&sipScore, id).Error
	return &sipScore, err
}

func (d *Dao) IncrSipScoreEntryCount(id uint32, incr int64, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	result := db.Model(&SipScoreModel{}).Where("id = ?", id).UpdateColumn("entry_count", gorm.Expr("entry_count + ?", incr))

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (d *Dao) IncrSipScoreCollectCount(sipScoreID uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Model(&SipScoreModel{}).Where("id = ?", sipScoreID).
		UpdateColumn("collect_count", gorm.Expr("collect_count + 1")).Error
}

func (d *Dao) DecrSipScoreCollectCount(sipScoreID uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Model(&SipScoreModel{}).Where("id = ? AND collect_count > 0", sipScoreID).
		UpdateColumn("collect_count", gorm.Expr("collect_count - 1")).Error
}

func (d *Dao) IncrSipScoreParticipantCount(sipScoreID uint32, incr int64, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Model(&SipScoreModel{}).Where("id = ?", sipScoreID).
		UpdateColumn("participant_count", gorm.Expr("participant_count + ?", incr)).Error
}

func (d *Dao) DeleteSipScore(id uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Delete(&SipScoreModel{}, id).Error
}

func (d *Dao) ListSipScoreNewest(limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	return d.listSipScoreByTimeField("updated_at", orderDirDesc, limit, tx...)
}

func (d *Dao) ListSipScoreNewestWithCursor(lastID uint32, lastUpdatedAt time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	return d.listSipScoreByTimeFieldWithCursor("updated_at", orderDirDesc, lastID, lastUpdatedAt, limit, tx...)
}

func (d *Dao) ListSipScoreHottest(limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	return d.listSipScoreByUintField("participant_count", orderDirDesc, limit, tx...)
}

func (d *Dao) ListSipScoreHottestWithCursor(lastID uint32, lastCount uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	return d.listSipScoreByUintFieldWithCursor("participant_count", orderDirDesc, lastID, lastCount, limit, tx...)
}

func (d *Dao) listSipScoreByTimeFieldWithCursor(field string, orderDir string, lastID uint32, lastTime time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	db := d.getDB(tx...)

	whereOp := ">"
	idOp := ">"
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		whereOp = "<"
		idOp = "<"
		order = field + " DESC, id DESC"
	}

	var sipScores []*SipScoreModel
	err := db.Where(field+" "+whereOp+" ? OR ("+field+" = ? AND id "+idOp+" ?)", lastTime, lastTime, lastID).
		Order(order).Limit(int(limit)).Find(&sipScores).Error

	return sipScores, err
}

func (d *Dao) listSipScoreByTimeField(field string, orderDir string, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	db := d.getDB(tx...)
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		order = field + " DESC, id DESC"
	}

	var sipScores []*SipScoreModel
	err := db.Order(order).Limit(int(limit)).Find(&sipScores).Error

	return sipScores, err
}

func (d *Dao) listSipScoreByUintFieldWithCursor(field string, orderDir string, lastID uint32, lastValue uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	db := d.getDB(tx...)
	whereOp := ">"
	idOp := ">"
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		whereOp = "<"
		idOp = "<"
		order = field + " DESC, id DESC"
	}

	var sipScores []*SipScoreModel
	err := db.Where(field+" "+whereOp+" ? OR ("+field+" = ? AND id "+idOp+" ?)", lastValue, lastValue, lastID).
		Order(order).Limit(int(limit)).Find(&sipScores).Error

	return sipScores, err
}

func (d *Dao) listSipScoreByUintField(field string, orderDir string, limit uint32, tx ...*gorm.DB) ([]*SipScoreModel, error) {
	db := d.getDB(tx...)
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		order = field + " DESC, id DESC"
	}

	var sipScores []*SipScoreModel
	err := db.Order(order).Limit(int(limit)).Find(&sipScores).Error

	return sipScores, err
}

// DecrSipScoreStats 多字段递减，减少 DB 操作
func (d *Dao) DecrSipScoreStats(sipScoreID uint32, entryCount uint32, participantCount uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)

	return db.Model(&SipScoreModel{}).
		Where("id = ?", sipScoreID).
		UpdateColumns(map[string]interface{}{
			"entry_count":       gorm.Expr("GREATEST(entry_count - ?, 0)", entryCount),
			"participant_count": gorm.Expr("GREATEST(participant_count - ?, 0)", participantCount),
		}).Error
}

// todo sipScoreID + name deletedAt 唯一索引，确保同一榜单内条目名称唯一
// todo sipScoreID ID 联合索引
// todo sipScoreID + UpdatedAt id 联合索引
// todo sipScoreID + participantCount id 联合索引
// todo sipScoreID + scoreAvg id 联合索引
// todo deletedAt 纳秒级别
// todo sipScoreID id participantCount 联合索引

type SipScoreEntryModel struct {
	ID             uint32 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      soft_delete.DeletedAt `gorm:"index;softDelete:nano"`
	SipScoreID     uint32
	LastModifiedBy uint32
	CreatorID      uint32
	IsReport       bool

	// 冗余统计字段 - 用于排序
	ParticipantCount uint32
	CommentCount     uint32
	ScoreTotal       uint32
	ScoreAvg         uint32 // 实际是 ScoreAvg * 100 保留两位小数

	// 用户可编辑字段
	Name        string `gorm:"type:varchar(100);not null;index:,class:FULLTEXT,option:WITH PARSER ngram"`
	Description string `gorm:"type:varchar(500)"`
	CoverImg    string `gorm:"type:varchar(255)"`
}

func (SipScoreEntryModel) TableName() string {
	return "sip_score_entries"
}

func (s *SipScoreEntryModel) Create() error {
	return dao.DB.Create(s).Error
}

func (s *SipScoreEntryModel) Save(tx ...*gorm.DB) error {
	db := dao.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Save(s).Error
}

func (s *SipScoreEntryModel) Delete(tx *gorm.DB) error {
	return tx.Delete(s).Error
}

func (s *SipScoreEntryModel) Get(id uint32) error {
	return dao.DB.First(s, id).Error
}

func (s *SipScoreEntryModel) BeReported() error {
	s.IsReport = true
	return s.Save()
}

func (d *Dao) GetSipScoreEntry(sipScoreID, entryID uint32, tx ...*gorm.DB) (*SipScoreEntryModel, error) {
	db := d.getDB(tx...)
	var entry SipScoreEntryModel
	err := db.Where("id = ? AND sip_score_id = ?", entryID, sipScoreID).First(&entry).Error
	return &entry, err
}

func (d *Dao) BatchCreateSipScoreEntries(entries []*SipScoreEntryModel, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Create(entries).Error
}

func (d *Dao) IncrSipScoreEntryScore(sipScoreID, entryID uint32, scoreIncr uint32, participantIncr uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)

	return db.Model(&SipScoreEntryModel{}).
		Where("id = ? AND sip_score_id = ?", entryID, sipScoreID).
		UpdateColumns(map[string]interface{}{
			"score_total":       gorm.Expr("score_total + ?", scoreIncr),
			"participant_count": gorm.Expr("participant_count + ?", participantIncr),
			"score_avg": gorm.Expr(
				"((score_total + ?) * 100) / (participant_count + ?)",
				scoreIncr,
				participantIncr,
			),
		}).Error
}

func (d *Dao) UpdateSipScoreEntry(sipScoreID, entryID uint32, update map[string]interface{}, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	result := db.Model(&SipScoreEntryModel{}).Where("id = ? AND sip_score_id = ?", entryID, sipScoreID).Updates(update)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// GetSipScoreEntryIDs 用于筛选 未被删除的 entryID
func (d *Dao) GetSipScoreEntryIDs(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) ([]uint32, error) {
	db := d.getDB(tx...)

	var result []uint32
	err := db.Model(&SipScoreEntryModel{}).
		Where("sip_score_id = ? AND id IN ?", sipScoreID, entryIDs).
		Pluck("id", &result).Error

	return result, err
}

func (d *Dao) DeleteSipScoreEntries(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Where("sip_score_id = ? AND id IN ?", sipScoreID, entryIDs).Delete(&SipScoreEntryModel{}).Error
}

func (d *Dao) GetSipScoreEntryStats(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) (entryCount uint32, participantCount uint32, err error) {
	db := d.getDB(tx...)

	type result struct {
		EntryCount       uint32
		ParticipantCount uint32
	}

	var r result

	err = db.Model(&SipScoreEntryModel{}).
		Select("COUNT(*) as entry_count, COALESCE(SUM(participant_count),0) as participant_count").
		Where("sip_score_id = ? AND id IN ?", sipScoreID, entryIDs).Scan(&r).Error

	if err != nil {
		return
	}

	return r.EntryCount, r.ParticipantCount, nil
}

func (d *Dao) ListSipScoreEntriesNewest(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByTimeField(sipScoreID, "updated_at", orderDirDesc, limit, tx...)
}

func (d *Dao) ListSipScoreEntriesNewestWithCursor(sipScoreID, lastID uint32, lastUpdatedAt time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByTimeFieldWithCursor(sipScoreID, "updated_at", "DESC", lastID, lastUpdatedAt, limit, tx...)
}

// 热门
func (d *Dao) ListSipScoreEntriesHottest(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintField(sipScoreID, "participant_count", "DESC", limit, tx...)
}

func (d *Dao) ListSipScoreEntriesHottestWithCursor(sipScoreID, lastID uint32, lastCount uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintFieldWithCursor(sipScoreID, "participant_count", "DESC", lastID, lastCount, limit, tx...)
}

// 高分 / 低分
func (d *Dao) ListSipScoreEntriesHighestScore(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintField(sipScoreID, "score_avg", "DESC", limit, tx...)
}

func (d *Dao) ListSipScoreEntriesHighestScoreWithCursor(sipScoreID, lastID uint32, lastScore uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintFieldWithCursor(sipScoreID, "score_avg", "DESC", lastID, lastScore, limit, tx...)
}

func (d *Dao) ListSipScoreEntriesLowestScore(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintField(sipScoreID, "score_avg", "ASC", limit, tx...)
}

func (d *Dao) ListSipScoreEntriesLowestScoreWithCursor(sipScoreID, lastID uint32, lastScore uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	return d.listSipScoreEntriesByUintFieldWithCursor(sipScoreID, "score_avg", "ASC", lastID, lastScore, limit, tx...)
}

// listSipScoreEntriesByTimeFieldWithCursor
// field: e.g. "updated_at", "created_at"
// orderDir: "ASC" or "DESC"
// 按时间字段排序的通用函数，支持 cursor 分页
func (d *Dao) listSipScoreEntriesByTimeFieldWithCursor(sipScoreID uint32, field string, orderDir string, lastID uint32, lastTime time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	db := d.getDB(tx...)

	whereOp := ">"
	idOp := ">"
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		whereOp = "<"
		idOp = "<"
		order = field + " DESC, id DESC"
	}

	// e.g. sip_score_id = ? AND (updated_at < ? OR (updated_at = ? AND id < ?))
	var entries []*SipScoreEntryModel
	err := db.Where("sip_score_id = ? AND ("+field+" "+whereOp+" ? OR ("+field+" = ? AND id "+idOp+" ?))",
		sipScoreID, lastTime, lastTime, lastID,
	).Order(order).Limit(int(limit)).Find(&entries).Error

	return entries, err
}

// listSipScoreEntriesByTimeField
// field: 同上
// orderDir: 同上
// 按时间字段排序的通用函数，适用于第一页请求（没有 cursor）
func (d *Dao) listSipScoreEntriesByTimeField(sipScoreID uint32, field string, orderDir string, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	db := d.getDB(tx...)
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		order = field + " DESC, id DESC"
	}

	var entries []*SipScoreEntryModel
	err := db.Where("sip_score_id = ?", sipScoreID).
		Order(order).Limit(int(limit)).Find(&entries).Error

	return entries, err
}

// listSipScoreEntriesByUintFieldWithCursor
// field: e.g. "participant_count" or "score_avg"
// orderDir: "ASC" or "DESC"
// 按 uint 字段排序的通用函数，支持 cursor 分页
func (d *Dao) listSipScoreEntriesByUintFieldWithCursor(sipScoreID uint32, field string, orderDir string, lastID uint32, lastValue uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	db := d.getDB(tx...)

	whereOp := ">"
	idOp := ">"
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		whereOp = "<"
		idOp = "<"
		order = field + " DESC, id DESC"
	}

	var entries []*SipScoreEntryModel
	err := db.
		Where("sip_score_id = ? AND ("+field+" "+whereOp+" ? OR ("+field+" = ? AND id "+idOp+" ?))",
			sipScoreID, lastValue, lastValue, lastID,
		).
		Order(order).Limit(int(limit)).Find(&entries).Error

	return entries, err
}

// listSipScoreEntriesByUintField
// field: e.g. 同上
// orderDir: 同上
// 按 uint 字段排序的通用函数，适用于第一页请求（没有 cursor）
func (d *Dao) listSipScoreEntriesByUintField(sipScoreID uint32, field string, orderDir string, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error) {
	db := d.getDB(tx...)
	order := field + " ASC, id ASC"
	if orderDir == orderDirDesc {
		order = field + " DESC, id DESC"
	}

	var entries []*SipScoreEntryModel
	err := db.
		Where("sip_score_id = ?", sipScoreID).
		Order(order).Limit(int(limit)).Find(&entries).Error

	return entries, err
}

// BatchListSipScoreEntriesHottest 批量获取多个榜单的热门条目，返回结果按榜单 ID 分组
// 不如窗口函数
func (d *Dao) BatchListSipScoreEntriesHottest(sipScoreIDs []uint32, limit uint32, tx ...*gorm.DB) (map[uint32][]*SipScoreEntryModel, error) {
	if len(sipScoreIDs) == 0 {
		return map[uint32][]*SipScoreEntryModel{}, nil
	}

	db := d.getDB(tx...)

	var entries []*SipScoreEntryModel

	// 按热门排序查全部
	err := db.
		Where("sip_score_id IN ?", sipScoreIDs).
		Order("sip_score_id ASC, participant_count DESC, id DESC").
		Find(&entries).Error

	if err != nil {
		return nil, err
	}

	// 分组并截断
	result := make(map[uint32][]*SipScoreEntryModel)

	for _, e := range entries {
		list := result[e.SipScoreID]
		if uint32(len(list)) >= limit {
			continue
		}
		result[e.SipScoreID] = append(list, e)
	}

	return result, nil
}

//// BatchListSipScoreEntriesHottest 批量获取多个 sipScoreID 的热门条目，返回结果按 sipScoreID 分组
//// 学习 SQL 语句是对的（虽然我数据库作业也是抄的）
//// 虽然因为 mysql-5.7 不支持，但保留给予提示——sql语句🐮
//func (d *Dao) BatchListSipScoreEntriesHottest(sipScoreIDs []uint32, limit uint32, tx ...*gorm.DB) (map[uint32][]*SipScoreEntryModel, error) {
//	if len(sipScoreIDs) == 0 {
//		return map[uint32][]*SipScoreEntryModel{}, nil
//	}
//
//	db := d.getDB(tx...)
//
//	var entries []*SipScoreEntryModel
//
//	err := db.Raw(`
//		SELECT * FROM (
//			SELECT *,
//			       ROW_NUMBER() OVER (
//			           PARTITION BY sip_score_id
//			           ORDER BY participant_count DESC, id DESC
//			       ) AS rn
//			FROM sip_score_entry_models
//			WHERE sip_score_id IN ?
//		) t
//		WHERE t.rn <= ?
//	`, sipScoreIDs, limit).Scan(&entries).Error
//
//	if err != nil {
//		return nil, err
//	}
//
//	// 转换成 map
//	result := make(map[uint32][]*SipScoreEntryModel)
//	for _, e := range entries {
//		result[e.SipScoreID] = append(result[e.SipScoreID], e)
//	}
//
//	return result, nil
//}

type SipScoreEntryCommentRating struct {
	ID             uint32 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      soft_delete.DeletedAt `gorm:"index;softDelete:nano"`
	CreatorID      uint32
	LastModifiedBy uint32
	SipScoreID     uint32
	EntryID        uint32
	CommentID      uint32
	LikeNum        uint32
}

func (SipScoreEntryCommentRating) TableName() string {
	return "sip_score_entry_comment_ratings"
}

func (d *Dao) CreateSipScoreEntryCommentRating(rating *SipScoreEntryCommentRating, tx ...*gorm.DB) (uint32, error) {
	db := d.getDB(tx...)
	err := db.Create(rating).Error
	return rating.ID, err
}
