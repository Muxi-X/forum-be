package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func (s *PostService) HandleReport(ctx context.Context, req *pb.HandleReportRequest, _ *pb.Response) error {
	logger.Info("PostService HandleReport")

	report, err := s.Dao.GetReport(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.NotFoundErr(errno.ErrItemNotFound, "report-"+strconv.Itoa(int(req.Id)))
		}

		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if req.Result == constvar.ValidReport {
		if err := s.Dao.ValidReport(report.TypeName, report.TargetId); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		s.CreateMessage(ctx, report.UserId, "One of your "+report.TypeName+"s has been deleted")

	} else if req.Result == constvar.InvalidReport {
		if err := s.Dao.InValidReport(req.Id, report.TypeName, report.TargetId); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
	} else {
		return errno.ServerErr(errno.ErrBadRequest, "result not legal")
	}

	return nil
}
