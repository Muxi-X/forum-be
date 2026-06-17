package user

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name LoginRequest

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name RegisterRequest

// GetInfoRequest 获取 info 请求
type GetInfoRequest struct {
	Ids []uint32 `json:"ids" binding:"required"`
} // @name GetInfoRequest

type userInfo struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}

// GetInfoResponse 获取 info 响应
type GetInfoResponse struct {
	List []userInfo `json:"list"`
} // @name GetInfoResponse

// GetProfileRequest 获取 profile 请求
type GetProfileRequest struct {
	Id uint32 `json:"id"`
} // @name GetProfileRequest

// UserProfile 获取 profile 响应
type UserProfile struct {
	Id                        uint32 `json:"id"`
	Name                      string `json:"name"`
	Avatar                    string `json:"avatar"`
	Email                     string `json:"email"`
	Role                      string `json:"role"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
	FollowingCount            uint32 `json:"following_count"`
	FollowerCount             uint32 `json:"follower_count"`
	IsFollowing               bool   `json:"is_following"`
} // @name UserProfile

// MyProfile 获取 my profile 响应
type MyProfile struct {
	Id                        uint32 `json:"id"`
	Name                      string `json:"name"`
	Avatar                    string `json:"avatar"`
	Email                     string `json:"email"`
	StudentId                 string `json:"student_id"`
	Role                      string `json:"role"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
	FollowingCount            uint32 `json:"following_count"`
	FollowerCount             uint32 `json:"follower_count"`
	IsFollowing               bool   `json:"is_following"`
} // @name MyProfile

type FollowRequest struct {
	TargetUserID uint32 `json:"target_user_id" binding:"required"`
} // @name FollowRequest

type FollowResponse struct {
	IsFollowing    bool   `json:"is_following"`
	FollowingCount uint32 `json:"following_count"`
	FollowerCount  uint32 `json:"follower_count"`
} // @name FollowResponse

type FollowListUser struct {
	Id             uint32 `json:"id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Role           string `json:"role"`
	Signature      string `json:"signature"`
	FollowingCount uint32 `json:"following_count"`
	FollowerCount  uint32 `json:"follower_count"`
	IsFollowing    bool   `json:"is_following"`
} // @name FollowListUser

type FollowListResponse struct {
	Users []FollowListUser `json:"users"`
} // @name FollowListResponse

// ListRequest 获取 userList 请求
type ListRequest struct {
	Team  uint32 `json:"team"`
	Group uint32 `json:"group"`
} // @name ListRequest

type user struct {
	Id     uint32 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
} // @name user

// ListResponse 获取 userList 响应
type ListResponse struct {
	Count uint32 `json:"count"`
	List  []user `json:"list"`
} // @name ListResponse

// UpdateInfoRequest 更新 userInfo 请求
type UpdateInfoRequest struct {
	Name                      string `json:"name"`
	AvatarURL                 string `json:"avatar_url"`
	Signature                 string `json:"signature"`
	IsPublicFeed              bool   `json:"is_public_feed"`
	IsPublicCollectionAndLike bool   `json:"is_public_collection_and_like"`
} // @name UpdateInfoRequest

type ListMessageResponse struct {
	Messages []string `json:"messages"`
}

type PrivateMessage struct {
	Id             string `json:"id"`
	SendUserId     string `json:"send_user_id"`
	PostId         string `json:"post_id"`
	CommentId      string `json:"comment_id"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	PostTitle      string `json:"post_title"`
	CommentContent string `json:"comment_content"`
	Avatar         string `json:"avatar"`
	SenderName     string `json:"sender_name"`
	Read           bool   `json:"read"`
	CreatedAt      string `json:"created_at"`
}

type ListPrivateMessageResponse struct {
	Messages []PrivateMessage `json:"messages"`
}

type CreateMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type CreatePrivateMessageRequest struct {
	ReceiveUserid  uint32 `json:"receive_userid" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Content        string `json:"content"`
	PostId         uint32 `json:"post_id" binding:"required"`
	CommentId      uint32 `json:"comment_id"`
	PostTitle      string `json:"post_title" binding:"required"`
	CommentContent string `json:"comment_content"`
}
