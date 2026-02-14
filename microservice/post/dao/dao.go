package dao

import (
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// Interface dao
type Interface interface {
	CreatePost(*PostModel) (uint32, error)
	ListUserCreatedPost(uint32) ([]uint32, error)
	ListMainPost(*PostModel, string, uint32, uint32, uint32, bool, string, uint32) ([]*PostInfo, error)
	GetPostInfo(uint32) (*PostInfo, error)
	GetPost(uint32) (*PostModel, error)
	IsUserCollectionPost(uint32, uint32) (bool, error)
	ListPostInfoByPostIds([]uint32, *PostModel, uint32, uint32, uint32, bool) ([]*pb.PostPartInfo, error)
	DeletePost(uint32, ...*gorm.DB) error
	ChangeQualityPost(uint32, bool) error
	CountPostByTime(string, string) (int64, error)

	CreateComment(*CommentModel) (uint32, error)
	GetCommentInfo(uint32) (*CommentInfo, error)
	GetComment(uint32) (*CommentModel, error)
	ListCommentByPostId(uint32) ([]*CommentInfo, error)
	GetCommentNumByPostId(uint32) (uint32, error)
	DeleteComment(uint32) error

	AddLike(uint32, Item) error
	RemoveLike(uint32, Item) error
	GetLikedNum(Item) (int64, error)
	IsUserHadLike(uint32, Item) (bool, error)
	ListUserLike(uint32) ([]*Item, error)

	CreatePost2Tag(Post2TagModel) error
	GetTagById(uint32) (*TagModel, error)
	GetTagByContent(string) (*TagModel, error)
	ListTagsByPostId(uint32) ([]string, error)

	AddTagToSortedSet(uint32, string) error
	ListPopularTags(string) ([]string, error)

	CreateCollection(*CollectionModel) (uint32, error)
	DeleteCollection(*CollectionModel) error
	GetCollectionNumByPostId(uint32) (uint32, error)
	ListCollectionByUserId(uint32) ([]uint32, error)

	ChangePostScore(uint32, int) error
	AddChangeRecord(uint32) error

	GetReport(uint32) (*ReportModel, error)
	CreateReport(*ReportModel) error
	ListReport(uint32, uint32, uint32, bool) ([]*pb.Report, error)
	GetReportNumByTypeNameAndId(string, uint32) (uint32, error)
	ValidReport(string, uint32) error
	InValidReport(uint32, string, uint32) error
	IsUserHadReportTarget(uint32, string, uint32) (bool, error)

	UpdateLastRead(uint32, string, string) error
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// init db
	model.DB.Init()

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}

	// 每小时同步一次post score 和  点赞
	go func() {
		for {
			time.Sleep(time.Hour)

			if err := dao.syncPostScore(); err != nil {
				log.Error(errno.ErrSyncPostScore.Error(), log.String(err.Error()))
			}

			if err := dao.syncItemLike(); err != nil {
				log.Error(errno.ErrSyncItemLike.Error(), log.String(err.Error()))
			}
		}
	}()
}

func GetDao() *Dao {
	return dao
}

func (d Dao) DeletePost(id uint32, tx ...*gorm.DB) error {
	db := d.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	post := &PostModel{}
	if err := post.Get(id); err != nil {
		return err
	}
	if err := post.Delete(db); err != nil {
		return err
	}

	ch := make(chan struct{})

	go func() {
		tags, err := d.ListTagsByPostId(id)
		ch <- struct{}{}
		if err != nil {
			log.Error(err.Error())
			return
		}

		for _, tag := range tags {
			tagId, err := d.getTagIdByContent(tag)
			if err != nil {
				log.Error(err.Error())
				return
			}

			isExist, err := d.isExistPostWithTagId(tagId)
			if err != nil {
				log.Error(err.Error())
				return
			}

			if !isExist {
				// TODO 添加redis事务保证数据一致性
				if err := d.Redis.ZRem("tags:", tagId).Err(); err != nil {
					log.Error(err.Error())
					return
				}

				isExist, err := d.isExistPostWithTagIdAndCategory(tagId, post.Category)
				if err != nil {
					log.Error(err.Error())
					return
				}

				if !isExist {
					if err := d.Redis.ZRem("tags:"+post.Category, tagId).Err(); err != nil {
						log.Error(err.Error())
						return
					}
				}
			}

		}
	}()

	<-ch // 上面获取成功后再删除
	if err := d.deletePost2TagByPostId(id); err != nil {
		return err
	}

	return d.Redis.ZRem("posts:", id).Err()
}

func (d Dao) DeleteComment(id uint32) error {
	comment := &CommentModel{}
	if err := comment.Get(id); err != nil {
		return err
	}
	if err := comment.Delete(); err != nil {
		return err
	}

	return d.ChangePostScore(comment.PostId, -constvar.CommentScore)
}
