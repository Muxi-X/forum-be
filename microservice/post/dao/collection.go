package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type CollectionModel struct {
	ID          uint32 `gorm:"primaryKey"`
	UserID      uint32 `gorm:"uniqueIndex:idx_user_target,priority:1"`
	ContentType uint32 `gorm:"uniqueIndex:idx_user_target,priority:2;index:idx_target,priority:1"`
	ContentID   uint32 `gorm:"uniqueIndex:idx_user_target,priority:3;index:idx_target,priority:2"`
	CreatedAt   time.Time
	// 使用这个字段来实现软删除与唯一索引的兼容
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_user_target,priority:4;index;softDelete:nano"`
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

func (d *Dao) CreateCollection(collection *CollectionModel, tx ...*gorm.DB) (uint32, error) {
	db := d.getDB(tx...)
	if err := db.Create(collection).Error; err != nil {
		return 0, err
	}
	return collection.ID, nil
}

func (d *Dao) TryCreateCollection(collection *CollectionModel, tx ...*gorm.DB) (bool, error) {
	db := d.getDB(tx...)

	result := db.Create(collection)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false, nil
		}
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (d *Dao) DeleteCollection(collection *CollectionModel, tx ...*gorm.DB) error {
	db := d.getDB(tx...)
	return db.Where("user_id = ? AND content_type = ? AND content_id = ?",
		collection.UserID, collection.ContentType, collection.ContentID).
		Delete(&CollectionModel{}).Error
}

func (d *Dao) TryDeleteCollection(collection *CollectionModel, tx ...*gorm.DB) (bool, error) {
	db := d.getDB(tx...)

	result := db.Where(
		"user_id = ? AND content_type = ? AND content_id = ?",
		collection.UserID, collection.ContentType, collection.ContentID,
	).Delete(&CollectionModel{})

	return result.RowsAffected > 0, result.Error
}

func (d *Dao) ListCollectionByUserId(userId uint32, contentType uint32) ([]uint32, error) {
	var ids []uint32
	err := d.DB.Model(&CollectionModel{}).
		Select("content_id").
		Where("user_id = ? AND content_type = ?", userId, contentType).
		Find(&ids).
		Error

	return ids, err
}

func (d *Dao) IsUserCollected(userID uint32, contentType uint32, contentID uint32, tx ...*gorm.DB) (bool, error) {
	db := d.getDB(tx...)
	var c CollectionModel
	err := db.
		Select("id").
		Where("user_id = ? AND content_type = ? AND content_id = ?",
			userID, contentType, contentID,
		).Take(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (d *Dao) GetCollectionNum(contentType uint32, contentId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&CollectionModel{}).
		Where("content_type = ? AND content_id = ?",
			contentType, contentId,
		).Count(&count).Error

	return uint32(count), err
}
