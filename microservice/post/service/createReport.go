package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
	"strconv"

	"gorm.io/gorm"
)

func (s *PostService) CreateReport(_ context.Context, req *pb.CreateReportRequest, _ *pb.Response) error {
	logger.Info("PostService CreateReport")

	var item dao.BeReporter
	var err error

	if req.TypeName == constvar.Post {
		item, err = s.Dao.GetPost(req.Id)
	} else if req.TypeName == constvar.Comment {
		item, err = s.Dao.GetComment(req.Id)
	} else {
		return errno.ServerErr(errno.ErrBadRequest, "wrong TypeName")
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.NotFoundErr(errno.ErrItemNotFound, req.TypeName+"-"+strconv.Itoa(int(req.Id)))
		}

		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 避免一个人重复举报
	isReported, err := s.Dao.IsUserHadReportTarget(req.UserId, req.TypeName, req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isReported {
		return errno.ServerErr(errno.ErrRepeatReport, "")
	}

	report := &dao.ReportModel{
		UserId:     req.UserId,
		CreateTime: util.GetCurrentTime(),
		TypeName:   req.TypeName,
		Cause:      req.Cause,
		Category:   req.Category,
		TargetId:   req.Id,
	}

	if err := s.Dao.CreateReport(report); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// auto ban when count >= constvar.BanNumber
	go func() {
		count, err := s.Dao.GetReportNumByTypeNameAndId(req.TypeName, req.Id)
		if err != nil {
			logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
		}

		if count >= constvar.BanNumber {
			if err := item.BeReported(); err != nil {
				logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
			}
		}
	}()

	return nil
}
