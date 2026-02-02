package service

import (
	"context"
	"forum-feed/dao"
	pb "forum-feed/proto"
	"forum-user/pkg/role"
	pbu "forum-user/proto"
	"forum/client"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

// List ... feed列表
func (s *FeedService) List(_ context.Context, req *pb.ListRequest, res *pb.ListResponse) error {
	logger.Info("FeedService List")

	var filter = &dao.FeedModel{}

	getResp, err := client.UserClient.GetProfile(context.TODO(), &pbu.GetRequest{Id: req.UserId})
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}

	if getResp.Role == constvar.NormalRole || getResp.Role == constvar.NormalAdminRole {
		filter.Domain = role.Role2Domain(getResp.Role)
	}

	feeds, err := s.Dao.List(filter, req.Offset, req.Limit, req.LastId, req.Pagination, req.TargetUserId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 数据格式化
	res.List = FormatListData(feeds)

	return nil
}

func FormatListData(list []*dao.FeedModel) []*pb.FeedItem {
	var result []*pb.FeedItem
	var time string
	var sourceId uint32

	for index, feed := range list {

		data := &pb.FeedItem{
			Id:          feed.Id,
			Action:      feed.Action,
			ShowDivider: false,
			CreateTime:  feed.CreateTime,
			User: &pb.User{
				Name:      feed.UserName,
				Id:        feed.UserId,
				AvatarUrl: feed.UserAvatar,
			},
			Source: &pb.Source{
				TypeName: feed.SourceTypeName,
				Id:       feed.SourceObjectId,
				Name:     feed.SourceObjectName,
			},
		}

		// showDivider --> 分割线
		// 需要分割的情况
		// 1.第一条数据 2.不同日期 3.不同项目
		if index == 0 {
			time = data.CreateTime
			sourceId = data.Source.Id
			data.ShowDivider = true
		} else if time != data.CreateTime {
			time = data.CreateTime
			data.ShowDivider = true
		} else if sourceId != data.Source.Id {
			sourceId = data.Source.Id
			data.ShowDivider = true
		}

		result = append(result, data)
	}

	return result
}
