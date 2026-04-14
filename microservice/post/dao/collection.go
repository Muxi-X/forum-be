package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// todo user_id + content_type + content_id 唯一索引
// todo content_type content_id 组合索引

type CollectionModel struct {
	ID uint32 `gorm:"primaryKey"`

	UserID      uint32 `gorm:"uniqueIndex:idx_user_target,priority:1"`
	ContentType uint8  `gorm:"uniqueIndex:idx_user_target,priority:2"`
	ContentID   uint32 `gorm:"uniqueIndex:idx_user_target,priority:3"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CollectionModel) TableName() string {
	return "collections"
}

// Create ...
func (c *CollectionModel) Create() error {
	return dao.DB.Create(c).Error
}

func (c *CollectionModel) Delete() error {
	return dao.DB.Delete(c).Error
}

func (d *Dao) CreateCollection(collection *CollectionModel) (uint32, error) {
	err := collection.Create()
	return collection.ID, err
}

func (d *Dao) DeleteCollection(collection *CollectionModel) error {
	return d.DB.
		Where("user_id = ? AND content_type = ? AND content_id = ?",
			collection.UserID, collection.ContentType, collection.ContentID,
		).Delete(&CollectionModel{}).Error
}

func (d *Dao) ListCollectionByUserId(userId uint32, contentType uint8) ([]uint32, error) {
	var ids []uint32
	err := d.DB.Model(&CollectionModel{}).
		Select("content_id").
		Where("user_id = ? AND content_type = ?", userId, contentType).
		Find(&ids).
		Error

	return ids, err
}

func (d *Dao) IsUserCollected(userId uint32, contentType uint8, contentId uint32) (bool, error) {
	var c CollectionModel
	err := d.DB.
		Select("id").
		Where("user_id = ? AND content_type = ? AND content_id = ?",
			userId, contentType, contentId,
		).Take(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (d *Dao) GetCollectionNum(contentType uint8, contentId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&CollectionModel{}).
		Where("content_type = ? AND content_id = ?",
			contentType, contentId,
		).Count(&count).Error

	return uint32(count), err
}
