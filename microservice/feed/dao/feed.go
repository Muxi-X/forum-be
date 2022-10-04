package dao

import (
	"forum/pkg/constvar"
)

type FeedModel struct {
	Id               uint32
	UserId           uint32
	UserName         string
	UserAvatar       string
	Action           string
	SourceTypeName   string
	SourceObjectName string
	SourceObjectId   uint32
	TargetUserId     uint32
	TypeName         string
	CreateTime       string
	Re               bool
}

func (*FeedModel) TableName() string {
	return "feeds"
}

// Create a new feed
func (f *FeedModel) Create() error {
	return dao.DB.Create(f).Error
}

func (f *FeedModel) Delete() error {
	f.Re = true
	return dao.DB.Save(f).Error
}

func (Dao) Create(feed *FeedModel) (uint32, error) {
	err := feed.Create()
	return feed.Id, err
}

func (d Dao) Delete(id uint32) error {
	var f *FeedModel
	err := d.DB.Model(f).Where("id = ? AND re = 0", id).First(f).Error
	if err != nil {
		return err
	}
	return f.Delete()
}

// // FilterParams provide filter's params.
// type FilterParams struct {
// 	UserId  uint32
// 	GroupId uint32
// }

// List ...
func (d *Dao) List(filter *FeedModel, offset, limit, lastId uint32, pagination bool) ([]*FeedModel, error) {
	var feeds []*FeedModel

	query := d.DB.Table("feeds").Select("feeds.*").Order("feeds.id desc").Where(filter).Where("re = 0")

	// 查找用户的feed
	if filter.UserId != 0 {
		query = query.Where("feeds.user_id = ?", filter.UserId)
	}

	// // 组别筛选
	// if filter.GroupId != 0 {
	// 	query = query.Where("users.group_id = ?", filter.GroupId).Joins("left join users on users.id = feeds.userid")
	// }

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("feeds.id < ?", lastId)
		}
	}

	err := query.Scan(&feeds).Error

	return feeds, err
}
