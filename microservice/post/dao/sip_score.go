package dao

import (
	"time"

	"gorm.io/gorm"
)

// 茶评 榜单 对象

// todo 创建索引
// todo 1. tag - content 二级索引
// todo 2.

type SipScoreModel struct {
	ID             uint32    `gorm:"primarykey;index:idx_rank,priority:3;index:idx_latest,priority:2;index:idx_creator,priority:2"`
	CreatedAt      time.Time `gorm:"index:idx_latest,priority:1"`
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
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

// todo sipScoreID + name 唯一索引，确保同一榜单内条目名称唯一
// todo sipScoreID ID 联合索引

type SipScoreEntryModel struct {
	ID             uint32 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
	SipScoreID     uint32
	LastModifiedBy uint32
	CreatorID      uint32
	IsReport       bool

	// 冗余统计字段 - 用于排序
	ParticipantCount uint32
	CommentCount     uint32
	ScoreTotal       uint32

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

func (d *Dao) BatchCreateSipScoreEntries(entries []*SipScoreEntryModel, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Create(entries).Error
}

func (d *Dao) UpdateSipScoreEntry(sipScoreID, entryID uint32, update map[string]interface{}, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	result := db.Model(&SipScoreEntryModel{}).Where("id = ? AND sip_score_id = ?", entryID, sipScoreID).Updates(update)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
