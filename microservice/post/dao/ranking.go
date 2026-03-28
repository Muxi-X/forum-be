package dao

import (
	"time"

	"gorm.io/gorm"
)

// 茶评 榜单 对象

// todo 创建索引
// 1. tag - content 二级索引

type RankingListModel struct {
	ID        uint32    `gorm:"primarykey;index:idx_rank,priority:3;index:idx_latest,priority:2;index:idx_creator,priority:2"`
	CreatedAt time.Time `gorm:"index:idx_latest,priority:1"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 用户可编辑的字段
	Name        string `gorm:"type:varchar(100);not null;index:,class:FULLTEXT,option:WITH PARSER ngram"`
	Description string `gorm:"type:varchar(500)"`
	CoverImg    string `gorm:"type:varchar(255)"`

	CreatorID uint32 `gorm:"index:idx_creator,priority:1"`

	// 是否被举报过多而被 ban 了
	IsReport bool `gorm:"index"`

	Domain   string `gorm:"type:varchar(20);not null;index"`
	Category string `gorm:"type:varchar(20);not null;index"`

	// 榜单内的茶评数量 todo 暂时不用，没发现相关的需求
	ItemCount uint32 `gorm:"type:int unsigned;default:0"`

	// 冗余字段，用于排序
	CollectCount     uint32 `gorm:"type:int unsigned;default:0;index:idx_rank,priority:1"`
	ParticipantCount uint32 `gorm:"type:int unsigned;default:0;index:idx_rank,priority:2"`
}

func (RankingListModel) TableName() string {
	return "ranking_lists"
}

// Create ...
func (r *RankingListModel) Create() error {
	return dao.DB.Create(r).Error
}

// Save ...
func (r *RankingListModel) Save(tx ...*gorm.DB) error {
	db := dao.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Save(r).Error
}

func (r *RankingListModel) Update() error {
	return dao.DB.Model(r).Where("id = ?", r.ID).Updates(map[string]interface{}{
		"name":        r.Name,
		"description": r.Description,
		"cover_img":   r.CoverImg,
	}).Error
}

func (r *RankingListModel) Delete(tx *gorm.DB) error {
	return tx.Delete(r).Error
}

func (r *RankingListModel) Get(id uint32) error {
	return dao.DB.First(r, id).Error
}

func (r *RankingListModel) BeReported() error {
	r.IsReport = true
	return r.Save()
}

type RankingListInfo struct {
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

func (d *Dao) CreateRankingList(rankingList *RankingListModel) (uint32, error) {
	err := rankingList.Create()
	return rankingList.ID, err
}
