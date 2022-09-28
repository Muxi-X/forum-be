package dao

import (
	pb "forum-post/proto"
	"gorm.io/gorm"
)

type CommentModel struct {
	Id         uint32
	TypeName   string // constvar.FirstLevelComment or constvar.SecondLevelComment
	Content    string
	FatherId   uint32
	CreateTime string
	Re         bool
	CreatorId  uint32
	PostId     uint32
	LikeNum    uint32
}

func (CommentModel) TableName() string {
	return "comments"
}

// Create ...
func (c *CommentModel) Create() error {
	return dao.DB.Create(c).Error
}

// Save ...
func (c *CommentModel) Save() error {
	return dao.DB.Save(c).Error
}

func (c *CommentModel) Get(id uint32) error {
	return dao.DB.Model(c).Where("id = ? AND re = 0", id).First(c).Error
}

func (c *CommentModel) Delete() error {
	c.Re = true
	return c.Save()
}

type CommentInfo struct {
	Id            uint32
	TypeName      string
	Content       string
	FatherId      uint32
	CreateTime    string
	CreatorId     uint32
	PostId        uint32
	CreatorName   string
	CreatorAvatar string
	LikeNum       uint32
}

func (Dao) CreateComment(comment *CommentModel) (uint32, error) {
	err := comment.Create()
	return comment.Id, err
}

func (d *Dao) ListCommentByPostId(postId uint32) ([]*pb.CommentInfo, error) {
	var comments []*pb.CommentInfo
	err := d.DB.Table("comments").Select("comments.id id, type_name, content, father_id, create_time time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = comments.creator_id").Where("post_id = ? AND comments.re = 0", postId).Find(&comments).Error
	return comments, err
}

func (d *Dao) GetComment(id uint32) (*CommentModel, error) {
	var comment CommentModel
	err := d.DB.Where("id = ? AND re = 0", id).First(&comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentInfo(commentId uint32) (*CommentInfo, error) {
	var comment CommentInfo
	err := d.DB.Table("comments").Select("comments.id id, type_name, content, father_id, create_time, creator_id, post_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = posts.creator_id").Where("comments.id = ? AND comments.re = 0", commentId).First(&comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentNumByPostId(postId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&CommentModel{}).Where("post_id = ? AND re = 0", postId).Count(&count).Error
	return uint32(count), err
}
