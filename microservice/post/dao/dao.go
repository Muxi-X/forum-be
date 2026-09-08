package dao

import (
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
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
	Transaction(fc func(tx *gorm.DB) error) error

	CreatePost(*PostModel) (uint32, error)
	ListUserCreatedPost(uint32, uint32, uint32, uint32, bool) ([]uint32, error)
	ListCollectionPostIDsByUser(uint32, uint32, uint32, uint32, bool) ([]uint32, error)
	ListMainPost(*PostModel, string, uint32, uint32, uint32, bool, string, uint32) ([]*PostInfo, error)
	GetPostInfo(uint32) (*PostInfo, error)
	GetPost(uint32) (*PostModel, error)
	ListPostInfoByPostIds([]uint32, *PostModel, uint32, uint32, uint32, bool) ([]*pb.PostPartInfo, error)
	DeletePost(uint32, ...*gorm.DB) error
	ChangeQualityPost(uint32, bool) error
	CountPostByTime(string, string) (int, error)

	CreateSipScore(sipScore *SipScoreModel, tx ...*gorm.DB) (uint32, error)
	BatchGetOrCreateTags(tags []string) ([]*TagModel, error)
	BatchCreateSipScoreTags(items []*SipScoreTagModel, tx ...*gorm.DB) error
	BatchAddTagsToSortedSet(tagIDs []uint32, category string) error
	BatchRemoveTagsFromSortedSet(tagIDs []uint32, category string) error
	UpdateSipScore(id uint32, update map[string]interface{}, tx ...*gorm.DB) error
	ListTagIDsBySipScoreId(sipScoreId uint32, tx ...*gorm.DB) ([]uint32, error)
	ListTagsBySipScoreId(sipScoreId uint32, tx ...*gorm.DB) ([]string, []uint32, error)
	DeleteSipScoreTagsBySipScoreId(sipScoreId uint32, tx ...*gorm.DB) error
	GetSipScore(id uint32, tx ...*gorm.DB) (*SipScoreModel, error)
	BatchCreateSipScoreEntries(entries []*SipScoreEntryModel, tx ...*gorm.DB) error
	IncrSipScoreEntryCount(id uint32, incr int64, tx ...*gorm.DB) error
	UpdateSipScoreEntry(sipScoreID, entryID uint32, update map[string]interface{}, tx ...*gorm.DB) error
	IncrSipScoreCollectCount(sipScoreID uint32, tx ...*gorm.DB) error
	DecrSipScoreCollectCount(sipScoreID uint32, tx ...*gorm.DB) error
	SearchSipScoreEntry(sipScoreID uint32, keyword string, limit uint32) ([]*SipScoreEntryModel, error)
	SearchSipScoreEntryWithCursor(sipScoreID uint32, keyword string, lastID uint32, lastCreatedAt time.Time, limit uint32) ([]*SipScoreEntryModel, error)

	ListSipScoreEntriesNewest(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesNewestWithCursor(sipScoreID, lastEntryID uint32, lastUpdatedAt time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesHottest(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesHottestWithCursor(sipScoreID, lastID uint32, lastCount uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesHighestScore(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesHighestScoreWithCursor(sipScoreID, lastID uint32, lastScore uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesLowestScore(sipScoreID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	ListSipScoreEntriesLowestScoreWithCursor(sipScoreID, lastID uint32, lastScore uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryModel, error)
	DeleteSipScore(id uint32, tx ...*gorm.DB) error
	GetSipScoreEntryIDs(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) ([]uint32, error)
	DeleteSipScoreEntries(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) error
	GetSipScoreEntryStats(sipScoreID uint32, entryIDs []uint32, tx ...*gorm.DB) (entryCount uint32, participantCount uint32, err error)
	DecrSipScoreStats(sipScoreID uint32, entryCount uint32, participantCount uint32, tx ...*gorm.DB) error
	ListSipScoreNewest(limit uint32, domain string, tx ...*gorm.DB) ([]*SipScoreModel, error)
	ListSipScoreNewestWithCursor(lastID uint32, lastUpdatedAt time.Time, limit uint32, domain string, tx ...*gorm.DB) ([]*SipScoreModel, error)
	ListSipScoreHottest(limit uint32, domain string, tx ...*gorm.DB) ([]*SipScoreModel, error)
	ListSipScoreHottestWithCursor(lastID uint32, lastCount uint32, limit uint32, domain string, tx ...*gorm.DB) ([]*SipScoreModel, error)
	SearchSipScore(keyword string, limit uint32, domain string) ([]*SipScoreModel, error)
	SearchSipScoreWithCursor(keyword string, lastID uint32, lastCreatedAt time.Time, limit uint32, domain string) ([]*SipScoreModel, error)
	ListSipScoreByCreator(userID uint32, offset uint32, limit uint32, lastID uint32, pagination bool) ([]*SipScoreModel, error)
	ListCollectedSipScoreByUser(userID uint32, offset uint32, limit uint32, lastID uint32, pagination bool) ([]*SipScoreModel, error)
	BatchListSipScoreEntriesHottest(sipScoreIDs []uint32, limit uint32, tx ...*gorm.DB) (map[uint32][]*SipScoreEntryModel, error)
	GetSipScoreEntry(sipScoreID, entryID uint32, tx ...*gorm.DB) (*SipScoreEntryModel, error)
	CreateSipScoreEntryCommentRating(rating *SipScoreEntryCommentRating, tx ...*gorm.DB) (uint32, error)
	GetSipScoreEntryCommentRating(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) (*SipScoreEntryCommentRating, error)
	GetSipScoreEntryCommentRatingByID(ratingID uint32, tx ...*gorm.DB) (*SipScoreEntryCommentRating, error)
	GetSipScoreEntryCommentRatingForUpdate(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) (*SipScoreEntryCommentRating, error)
	GetSipScoreEntryCommentRatingByUser(sipScoreID, entryID, userID uint32, tx ...*gorm.DB) (*SipScoreEntryCommentRating, error)
	GetSipScoreEntryCommentRatingByUserForUpdate(sipScoreID, entryID, userID uint32, tx ...*gorm.DB) (*SipScoreEntryCommentRating, error)
	ListSipScoreEntryCommentRatings(sipScoreID, entryID, offset, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryCommentRating, error)
	ListSipScoreEntryCommentRatingsNewest(sipScoreID, entryID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryCommentRating, error)
	ListSipScoreEntryCommentRatingsNewestWithCursor(sipScoreID, entryID, lastID uint32, lastCreatedAt time.Time, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryCommentRating, error)
	ListSipScoreEntryCommentRatingsHottest(sipScoreID, entryID, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryCommentRating, error)
	ListSipScoreEntryCommentRatingsHottestWithCursor(sipScoreID, entryID, lastID uint32, lastLikeNum uint32, limit uint32, tx ...*gorm.DB) ([]*SipScoreEntryCommentRating, error)
	UpdateSipScoreEntryScoreByRatingDelta(sipScoreID, entryID uint32, delta int, tx ...*gorm.DB) error
	UpdateSipScoreEntryCommentRating(sipScoreID, entryID, ratingID uint32, update map[string]interface{}, tx ...*gorm.DB) error
	IncrSipScoreEntryCommentRatingCommentNum(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) error
	DecrSipScoreEntryCommentRatingCommentNum(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) error
	IncrSipScoreEntryCommentRatingLikeNum(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) error
	DecrSipScoreEntryCommentRatingLikeNum(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) error
	IncrSipScoreEntryCommentCount(sipScoreID, entryID uint32, tx ...*gorm.DB) error
	DecrSipScoreEntryCommentCount(sipScoreID, entryID uint32, tx ...*gorm.DB) error
	DeleteSipScoreEntryCommentRating(sipScoreID, entryID, ratingID uint32, tx ...*gorm.DB) error
	IncrSipScoreParticipantCount(sipScoreID uint32, incr int64, tx ...*gorm.DB) error
	IncrSipScoreEntryScore(sipScoreID, entryID uint32, scoreIncr uint32, participantIncr uint32, tx ...*gorm.DB) error
	DecrSipScoreEntryScore(sipScoreID, entryID uint32, rating uint32, tx ...*gorm.DB) error
	LockSipScoreEntryForUpdate(sipScoreID, entryID uint32, tx ...*gorm.DB) error

	CreateComment(comment *CommentModel, tx ...*gorm.DB) (uint32, error)
	GetCommentInfo(uint32) (*CommentInfo, error)
	GetComment(uint32, ...*gorm.DB) (*CommentModel, error)
	ListCommentByTarget(targetID uint32, targetType string) ([]*CommentInfo, error)
	ListCommentByPostId(postId uint32) ([]*CommentInfo, error)
	GetCommentNumByTarget(targetID uint32, targetType string) (uint32, error)
	GetCommentNumByPostId(uint32) (uint32, error)
	BatchListCommentsByTargets(targetIDs []uint32, targetType string, limit int) (map[uint32][]*CommentInfo, error)
	BatchGetCommentNumByTargets(targetIDs []uint32, targetType string) (map[uint32]uint32, error)
	ListPrimaryCommentsNewest(targetID uint32, targetType string, limit uint32) ([]*CommentInfo, error)
	ListPrimaryCommentsNewestWithCursor(targetID uint32, targetType string, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error)
	ListPrimaryCommentsHottest(targetID uint32, targetType string, limit uint32) ([]*CommentInfo, error)
	ListPrimaryCommentsHottestWithCursor(targetID uint32, targetType string, lastID uint32, lastLikeNum uint32, limit uint32) ([]*CommentInfo, error)
	BatchListCommentsByFatherIDs(fatherIDs []uint32, limit int) (map[uint32][]*CommentInfo, error)
	BatchGetCommentNumByFatherIDs(fatherIDs []uint32) (map[uint32]uint32, error)
	ListSubCommentsNewest(fatherID, limit uint32) ([]*CommentInfo, error)
	ListSubCommentsNewestWithCursor(fatherID, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error)
	ListSubCommentsHottest(fatherID, limit uint32) ([]*CommentInfo, error)
	ListSubCommentsHottestWithCursor(fatherID, lastID uint32, lastLikeNum, limit uint32) ([]*CommentInfo, error)
	IncrCommentSubNum(commentID uint32, tx ...*gorm.DB) error
	DecrCommentSubNum(commentID uint32, tx ...*gorm.DB) error
	DeleteComment(uint32, ...*gorm.DB) error

	AddLike(uint32, Item) error
	RemoveLike(uint32, Item) error
	GetLikedNum(Item) (int64, error)
	IsUserHadLike(uint32, Item) (bool, error)
	ListUserLike(uint32) ([]*Item, error)

	CreatePost2Tag(Post2TagModel) error
	GetTagById(uint32) (*TagModel, error)
	GetTagByContent(string) (*TagModel, error)
	ListTagsByPostId(uint32) ([]string, []uint32, error)

	AddTagToSortedSet(uint32, string) error
	ListPopularTags(string) ([]string, error)

	CreateCollection(collection *CollectionModel, tx ...*gorm.DB) (uint32, error)
	TryCreateCollection(collection *CollectionModel, tx ...*gorm.DB) (bool, error)
	DeleteCollection(collection *CollectionModel, tx ...*gorm.DB) error
	TryDeleteCollection(collection *CollectionModel, tx ...*gorm.DB) (bool, error)
	GetCollectionNum(contentType uint32, contentId uint32) (uint32, error)
	ListCollectionByUserId(userId uint32, contentType uint32) ([]uint32, error)
	IsUserCollected(userID uint32, contentType uint32, contentID uint32, tx ...*gorm.DB) (bool, error)
	ListIsUserCollected(userID, contentType uint32, contentIDs []uint32, tx ...*gorm.DB) (map[uint32]bool, error)

	ChangePostScore(uint32, int) error
	AddChangeRecord(uint32) error

	GetReport(uint32) (*ReportModel, error)
	CreateReport(*ReportModel) error
	ListReport(uint32, uint32, uint32, bool) ([]*pb.Report, error)
	GetReportNumByTypeNameAndId(string, uint32) (uint32, error)
	ValidReport(string, uint32) error
	InValidReport(uint32, string, uint32) error
	IsUserHadReportTarget(uint32, string, uint32) (bool, error)
	CreateFeedback(*FeedbackModel) error

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

func (d *Dao) DeletePost(id uint32, tx ...*gorm.DB) error {
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

	_, tagIds, err := d.ListTagsByPostId(id)
	if err != nil {
		return err
	}

	if err := d.deletePost2TagByPostId(id); err != nil {
		return err
	}

	go func() {
		for _, tagId := range tagIds {
			isExist, err := d.isExistPostWithTagId(int(tagId))
			if err != nil {
				log.Error(err.Error())
				continue
			}

			if !isExist {
				pipe := d.Redis.Pipeline()
				pipe.ZRem("tags:", tagId)

				isExist, err := d.isExistPostWithTagIdAndCategory(int(tagId), post.Category)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				if !isExist {
					pipe.ZRem("tags:"+post.Category, tagId)
				}

				if _, err := pipe.Exec(); err != nil {
					log.Error(err.Error())
				}
			}
		}
	}()

	return d.Redis.ZRem("posts:", id).Err()
}

func (d *Dao) DeleteComment(id uint32, tx ...*gorm.DB) error {
	db := d.getDB(tx...)

	var comment CommentModel
	if err := db.Where("id = ? AND deleted_at = 0", id).First(&comment).Error; err != nil {
		return err
	}
	return comment.Delete(db)
}

func (d *Dao) Transaction(fc func(tx *gorm.DB) error) error {
	return d.DB.Transaction(fc)
}

func (d *Dao) getDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return d.DB
}

func (d *Dao) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return d.Redis.SetNX(key, value, expiration).Result()
}

func (d *Dao) Exists(key string) (int64, error) {
	return d.Redis.Exists(key).Result()
}
