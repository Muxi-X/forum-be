package dao

import (
	"time"

	"github.com/ShiinaOrez/GoSecurity/security"
)

type UserModel struct {
	ID                        uint32 `gorm:"column:id;primaryKey;not null"`
	Email                     string `gorm:"column:email;default:null"`
	Name                      string `gorm:"column:name;"`
	Avatar                    string `gorm:"column:avatar;"`
	HashPassword              string `gorm:"column:hash_password;"`
	Role                      string `gorm:"column:role;"`
	Signature                 string `gorm:"column:signature;"`
	StudentId                 string `gorm:"column:student_id;default:null"`
	Re                        bool   `gorm:"column:re;"`
	IsPublicFeed              bool
	IsPublicCollectionAndLike bool
}

func (UserModel) TableName() string {
	return "users"
}

// Create ...
func (u *UserModel) Create() error {
	return dao.DB.Create(u).Error
}

// Save ...
func (u *UserModel) Save() error {
	return dao.DB.Save(u).Error
}

func (u *UserModel) Update() error {
	err := dao.DB.Table("users").Where("id = ?", u.Id).Updates(map[string]interface{}{
		"name":                          u.Name,
		"avatar":                        u.Avatar,
		"signature":                     u.Signature,
		"is_public_collection_and_like": u.IsPublicCollectionAndLike,
		"is_public_feed":                u.IsPublicFeed,
	}).Error

	return err
}

// generatePasswordHash pwd -> hashPwd
func generatePasswordHash(password string) string {
	return security.GeneratePasswordHash(password)
}

func (u *UserModel) CheckPassword(password string) bool {
	return security.CheckPasswordHash(password, u.HashPassword)
}

type UserFollowModel struct {
	Id         uint32    `gorm:"column:id;primaryKey"`
	FollowerID uint32    `gorm:"column:follower_id;uniqueIndex:idx_user_follow,priority:1;index"`
	FolloweeID uint32    `gorm:"column:followee_id;uniqueIndex:idx_user_follow,priority:2;index"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (UserFollowModel) TableName() string {
	return "user_follows"
}
