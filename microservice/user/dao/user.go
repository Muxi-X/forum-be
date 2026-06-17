package dao

import (
	"errors"
	"forum/model"
	"forum/pkg/constvar"
	"strconv"

	"gorm.io/gorm"
)

type RegisterInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	StudentId string `json:"student_id"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

// GetUser get a single user by id
func (d *Dao) GetUser(id uint32) (*UserModel, error) {
	user := &UserModel{}
	err := d.DB.Where("id = ? AND re = 0", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func (d *Dao) BatchGetUser(ids []uint32) (map[uint32]*UserModel, error) {
	var users []*UserModel
	if err := d.DB.Where("id IN (?) AND re = 0", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	userMap := make(map[uint32]*UserModel, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	return userMap, nil
}

// GetUserByIds get user by id array
func (d *Dao) GetUserByIds(ids []uint32) ([]*UserModel, error) {
	var list []*UserModel
	if err := d.DB.Where("id IN (?) AND re = 0", ids).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// GetUserByEmail get a user by email.
func (d *Dao) GetUserByEmail(email string) (*UserModel, error) {
	u := &UserModel{}
	err := d.DB.Where("email = ? AND re = 0", email).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (d *Dao) UpdatePassword(userID uint32, newPassword string) error {
	hashedPassword := generatePasswordHash(newPassword)

	// 只更新 hash_password 字段
	return d.DB.Model(&UserModel{}).
		Where("id = ?", userID).
		Update("hash_password", hashedPassword).Error
}

// ListUser list users
func (d *Dao) ListUser(offset, limit, lastId uint32, filter *UserModel) ([]*UserModel, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	var list []*UserModel

	query := d.DB.Model(&UserModel{}).Where(filter).Offset(int(offset)).Limit(int(limit))

	if lastId != 0 {
		query = query.Where("id < ?", lastId).Order("id DESC")
	}

	if err := query.Scan(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// GetUserByStudentId get a user by studentId.
func (d *Dao) GetUserByStudentId(studentId string) (*UserModel, error) {
	u := &UserModel{}
	err := d.DB.Where("student_id = ? AND re = 0", studentId).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (Dao) RegisterUser(info *RegisterInfo) error {
	user := &UserModel{
		Name:                      info.Name,
		Email:                     info.Email,
		StudentId:                 info.StudentId,
		HashPassword:              generatePasswordHash(info.Password),
		Role:                      info.Role,
		Re:                        false,
		IsPublicCollectionAndLike: true,
		IsPublicFeed:              true,
	}

	return user.Create()
}

func (Dao) AddPublicPolicy(role string, userId uint32) error {
	ok, err := model.CB.Self.AddPolicy(role, constvar.CollectionAndLike+":"+strconv.Itoa(int(userId)), constvar.Read)
	if err != nil || !ok {
		return errors.New("AddPolicy CollectionAndLike fail")
	}

	ok, err = model.CB.Self.AddPolicy(role, constvar.Feed+":"+strconv.Itoa(int(userId)), constvar.Read)
	if err != nil || !ok {
		return errors.New("AddPolicy Feed fail")
	}
	return nil
}

func (d *Dao) ToggleFollow(followerID, followeeID uint32) (bool, error) {
	var isFollowing bool
	err := d.DB.Transaction(func(tx *gorm.DB) error {
		var follow UserFollowModel
		err := tx.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).First(&follow).Error
		if err == nil {
			if err = tx.Delete(&follow).Error; err != nil {
				return err
			}
			isFollowing = false
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err = tx.Create(&UserFollowModel{
			FollowerID: followerID,
			FolloweeID: followeeID,
		}).Error; err != nil {
			return err
		}
		isFollowing = true
		return nil
	})
	return isFollowing, err
}

func (d *Dao) CountFollowing(userID uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&UserFollowModel{}).Where("follower_id = ?", userID).Count(&count).Error
	return uint32(count), err
}

func (d *Dao) CountFollowers(userID uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&UserFollowModel{}).Where("followee_id = ?", userID).Count(&count).Error
	return uint32(count), err
}

func (d *Dao) IsFollowing(followerID, followeeID uint32) (bool, error) {
	var follow UserFollowModel
	err := d.DB.Select("id").Where("follower_id = ? AND followee_id = ?", followerID, followeeID).First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

type FollowListUser struct {
	Id        uint32
	Name      string
	Avatar    string
	Role      string
	Signature string
}

type followCountRow struct {
	Id    uint32
	Count uint32
}

func (d *Dao) ListFollowing(userID, limit, offset uint32) ([]*FollowListUser, error) {
	var users []*FollowListUser
	err := d.DB.Table("user_follows f").
		Select("u.id, u.name, COALESCE(u.avatar, '') AS avatar, u.role, COALESCE(u.signature, '') AS signature").
		Joins("JOIN users u ON u.id = f.followee_id AND COALESCE(u.re, 0) = 0").
		Where("f.follower_id = ?", userID).
		Order("f.created_at DESC, f.id DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&users).Error
	return users, err
}

func (d *Dao) ListFollowers(userID, limit, offset uint32) ([]*FollowListUser, error) {
	var users []*FollowListUser
	err := d.DB.Table("user_follows f").
		Select("u.id, u.name, COALESCE(u.avatar, '') AS avatar, u.role, COALESCE(u.signature, '') AS signature").
		Joins("JOIN users u ON u.id = f.follower_id AND COALESCE(u.re, 0) = 0").
		Where("f.followee_id = ?", userID).
		Order("f.created_at DESC, f.id DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&users).Error
	return users, err
}

func (d *Dao) BatchCountFollowing(userIDs []uint32) (map[uint32]uint32, error) {
	counts := make(map[uint32]uint32, len(userIDs))
	if len(userIDs) == 0 {
		return counts, nil
	}

	var rows []followCountRow
	err := d.DB.Table("user_follows").
		Select("follower_id AS id, COUNT(*) AS count").
		Where("follower_id IN ?", userIDs).
		Group("follower_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.Id] = row.Count
	}
	return counts, nil
}

func (d *Dao) BatchCountFollowers(userIDs []uint32) (map[uint32]uint32, error) {
	counts := make(map[uint32]uint32, len(userIDs))
	if len(userIDs) == 0 {
		return counts, nil
	}

	var rows []followCountRow
	err := d.DB.Table("user_follows").
		Select("followee_id AS id, COUNT(*) AS count").
		Where("followee_id IN ?", userIDs).
		Group("followee_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.Id] = row.Count
	}
	return counts, nil
}

func (d *Dao) BatchIsFollowing(followerID uint32, followeeIDs []uint32) (map[uint32]bool, error) {
	following := make(map[uint32]bool, len(followeeIDs))
	if followerID == 0 || len(followeeIDs) == 0 {
		return following, nil
	}

	var rows []followCountRow
	err := d.DB.Table("user_follows").
		Select("followee_id AS id, 1 AS count").
		Where("follower_id = ? AND followee_id IN ?", followerID, followeeIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		following[row.Id] = true
	}
	return following, nil
}
