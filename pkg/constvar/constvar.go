package constvar

const (
	DefaultLimit = 50

	// 角色权限
	Normal     = 1 // 普通用户
	Admin      = 3 // 管理员
	SuperAdmin = 7 // 超管

	// 权限限制等级
	AuthLevelNormal     = 1 // 普通用户级别
	AuthLevelAdmin      = 2 // 管理员级别
	AuthLevelSuperAdmin = 4 // 超管级别

	// user role
	NormalRole      = "Normal"
	NormalAdminRole = "NormalAdmin"
	MuxiRole        = "Muxi"
	MuxiAdminRole   = "MuxiAdmin"
	SuperAdminRole  = "SuperAdmin"

	// item.TypeName
	Post    = "post"
	Comment = "comment"

	Collection        = "collection"
	CollectionAndLike = "collection-like"
	Feed              = "feed"

	// domain
	NormalDomain = "normal"
	MuxiDomain   = "muxi"
	AllDomain    = "all"

	// casbin
	Write = "write"
	Read  = "read"

	// comment
	SubPost            = "sub-post"
	FirstLevelComment  = "first-level"
	SecondLevelComment = "second-level"

	// score
	LikeScore       = 2
	CollectionScore = 5
	CommentScore    = 3

	// report
	BanNumber     = 5
	ValidReport   = "valid"
	InvalidReport = "invalid"
)
