package dao

import (
	"errors"
	pb "forum-post/proto"
	"forum/pkg/constvar"
)

type ReportModel struct {
	Id         uint32
	TargetId   uint32
	UserId     uint32
	CreateTime string
	TypeName   string
	Cause      string
	Category   string
}

func (ReportModel) TableName() string {
	return "reports"
}

// Create ...
func (r *ReportModel) Create() error {
	return dao.DB.Create(r).Error
}

func (r *ReportModel) Delete() error {
	return dao.DB.Delete(r).Error
}

func (Dao) CreateReport(report *ReportModel) error {
	return report.Create()
}

func (d Dao) GetReport(id uint32) (*ReportModel, error) {
	var report ReportModel
	err := d.DB.Table("reports").Where("id = ?", id).First(&report).Error

	return &report, err
}

func (d Dao) ValidReport(typeName string, targetId uint32) error {
	tx := d.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 只要有一条有效举报就会删post,comment，故其余举报不再需要，一并删除
	if err := tx.Table("reports").Where("type_name = ? AND target_id = ?", typeName, targetId).Delete(&ReportModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var err error

	if typeName == constvar.Post {
		err = d.DeletePost(targetId, tx)
	} else if typeName == constvar.Comment {
		err = d.DeleteComment(targetId, tx)
	} else {
		err = errors.New("wrong TypeName")
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d Dao) InValidReport(id uint32, typeName string, targetId uint32) error {
	if err := d.DB.Table("reports").Where("id = ?", id).Delete(&ReportModel{}).Error; err != nil {
		return err
	}

	count, err := d.GetReportNumByTypeNameAndId(typeName, targetId)
	if err != nil {
		return err
	}

	if count == constvar.BanNumber-1 { // cancel auto ban
		var item BeReporter

		if typeName == constvar.Post {
			item, err = d.GetPost(targetId)
		} else if typeName == constvar.Comment {
			item, err = d.GetComment(targetId)
		} else {
			err = errors.New("wrong TypeName")
		}

		if err != nil {
			return err
		}

		return item.CancelReported()
	}

	return nil
}

func (d *Dao) ListReport(offset, limit, lastId uint32, pagination bool) ([]*pb.Report, error) {
	var reports []*pb.Report

	query := d.DB.Table("reports").Select("reports.id id, reports.target_id, reports.user_id, reports.create_time, reports.category, cause, reports.type_name, u.name user_name, u.avatar user_avatar").Joins("join users u on u.id = reports.user_id")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("reports.id < ?", lastId)
		}
	}

	if err := query.Scan(&reports).Error; err != nil {
		return nil, err
	}

	for i, report := range reports {
		if report.TypeName == constvar.Post {
			if err := d.DB.Select("title be_reported_content, u.id be_reported_user_id, u.name be_reported_user_name").Joins("join posts p on p.id = reports.target_id").Joins("join users u on u.id = p.creator_id").Where("reports.id = ?", report.Id).First(&reports[i]).Error; err != nil {
				return nil, err
			}
		} else if report.TypeName == constvar.Comment {
			if err := d.DB.Select("content be_reported_content, u.id be_reported_user_id, u.name be_reported_user_name").Joins("join comments c on c.id = reports.target_id").Joins("join users u on u.id = c.creator_id").Where("reports.id = ?", report.Id).First(&reports[i]).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("wrong TypeName")
		}
	}

	return reports, nil
}

func (d *Dao) GetReportNumByTypeNameAndId(typeName string, id uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&ReportModel{}).Where("type_name = ? AND target_id = ?", typeName, id).Count(&count).Error
	return uint32(count), err
}

func (d *Dao) IsUserHadReportTarget(userId uint32, typeName string, id uint32) (bool, error) {
	var count int64
	if err := d.DB.Table("reports").Where("user_id = ? AND type_name = ? AND id = ?", userId, typeName, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count != 0, nil
}

type BeReporter interface {
	BeReported() error
	CancelReported() error
}
