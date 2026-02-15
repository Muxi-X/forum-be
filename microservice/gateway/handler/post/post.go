package post

import (
	"forum-gateway/dao"
	"forum-gateway/handler/comment"
	pb "forum-post/proto"
	"reflect"
)

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

type UpdateInfoRequest struct {
	Id       uint32   `json:"id" binding:"required"`
	Domain   string   `json:"domain" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags" binding:"required"`
}

type Post struct {
	Id            uint32             `json:"id"`
	Title         string             `json:"title"`
	Time          string             `json:"time"`
	Category      string             `json:"category"`
	CreatorId     uint32             `json:"creator_id"`
	CreatorName   string             `json:"creator_name"`
	CreatorAvatar string             `json:"creator_avatar"`
	CommentNum    uint32             `json:"comment_num"`
	LikeNum       uint32             `json:"like_num"`
	IsLiked       bool               `json:"is_liked"`
	IsCollection  bool               `json:"is_collection"`
	Comments      []*comment.Comment `json:"comments"`
	Tags          []string           `json:"tags"`
	ContentType   string             `json:"content_type"` // md or rtf
	Summary       string             `json:"summary"`
	CollectionNum uint32             `json:"collection_num"`
}

type TrimHtmlRequest struct {
	Data string `json:"data" binding:"required"`
}

type CreateRequest struct {
	Domain          string   `json:"domain" binding:"required"` // normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal)
	Content         string   `json:"content" binding:"required"`
	CompiledContent string   `json:"compiled_content"`
	Title           string   `json:"title,omitempty" binding:"required"`
	Category        string   `json:"category,omitempty" binding:"required"`
	ContentType     string   `json:"content_type" binding:"required"` // md or rtf
	Summary         string   `json:"summary" binding:"required"`
	Tags            []string `json:"tags" binding:"required"`
}

type GetPostResponse struct {
	info
	CommentNum      uint32     `json:"comment_num"`
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	IsCollection    bool       `json:"is_collection"`
	SubPosts        []*SubPost `json:"sub_posts"`
	Tags            []string   `json:"tags"`
	ContentType     string     `json:"content_type"` // md or rtf
	CompiledContent string     `json:"compiled_content"`
	Summary         string     `json:"summary"`
	CollectionNum   uint32     `json:"collection_num"`
}

type SubPost struct {
	info
	ImgUrl     string     `json:"img_url"`
	CommentNum uint32     `json:"comment_num"`
	Comments   []*Comment `json:"comments"`
}

type Comment struct {
	info
	BeRepliedId       uint32 `json:"be_replied_id"`
	BeRepliedUserId   uint32 `json:"be_replied_user_id"`
	BeRepliedContent  string `json:"be_replied_content"`
	BeRepliedUserName string `json:"be_replied_user_name"`
}

type info struct {
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	Time          string `json:"time"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	LikeNum       uint32 `json:"like_num"`
	IsLiked       bool   `json:"is_liked"`
}

type PostPartInfo struct {
	Id            uint32   `json:"id"`
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Category      string   `json:"category"`
	Time          string   `json:"time"`
	CreatorId     uint32   `json:"creator_id"`
	CreatorName   string   `json:"creator_name"`
	CreatorAvatar string   `json:"creator_avatar"`
	CommentNum    uint32   `json:"comment_num"`
	CollectionNum uint32   `json:"collection_num"`
	LikeNum       uint32   `json:"like_num"`
	Tags          []string `json:"tags"`
}

type UnReadNum struct {
	Num      uint32 `json:"num"`
	Category string `json:"category"`
}

type PostPartInfoResponse struct {
	Posts []*PostPartInfo `json:"posts"`
}

type ListMainPostResponse struct {
	Posts []*Post `json:"posts"`
}

type UnReadNumResponse struct {
	UnReadNum []*UnReadNum `json:"un_read_num"`
}

// 这里用了 generics and reflect, 更一般的写法应该是用interface
func setInfo[T pb.CommentInfo | pb.Post](info *info, comment *T) {
	typeT := reflect.TypeOf(*comment)
	value := reflect.ValueOf(comment).Elem()

	for i := 0; i < typeT.NumField(); i++ {
		v := value.Field(i)
		field := typeT.Field(i)

		switch field.Name {
		case "Id":
			info.Id = uint32(v.Uint())
		case "Content":
			info.Content = v.String()
		case "IsLiked":
			info.IsLiked = v.Bool()
		case "CreatorName":
			info.CreatorName = v.String()
		case "CreatorId":
			info.CreatorId = uint32(v.Uint())
		case "LikeNum":
			info.LikeNum = uint32(v.Uint())
		case "CreatorAvatar":
			info.CreatorAvatar = v.String()
		case "Time", "CreateTime":
			info.Time = v.String()
		}
	}
}

type QiNiuToken struct {
	Token string `json:"token"`
}

type IdResponse struct {
	Id uint32 `json:"id"`
}
