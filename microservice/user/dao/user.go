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
		userMap[user.Id] = user
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
