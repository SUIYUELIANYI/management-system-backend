package file

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitViewingRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitViewingRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitViewingRecordLogic {
	return &SubmitViewingRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitViewingRecordLogic) SubmitViewingRecord(req *types.SubmitViewingRecordReq) error {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	_, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 查看视频是否存在
	_, err = l.svcCtx.FileModel.FindOne(l.ctx, req.FileId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "视频不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 查看该用户是否有过观看本视频的记录
	record, err := l.svcCtx.ViewingRecordModel.FindOneByUserIdFileId(l.ctx, userId, req.FileId)
	switch err {
	case nil: // 说明有观看记录，则记录较长的时间
		if req.Duration > record.Duration {
			record.Duration = req.Duration
			if err = l.svcCtx.ViewingRecordModel.Update(l.ctx, record); err != nil {
				return response.Error(500, err.Error())
			}
		}
	case models.ErrNotFound: // 说明没有观看记录
		// 插入观看记录
		viewingRecord := new(models.ViewingRecord)
		viewingRecord.FileId = req.FileId
		viewingRecord.UserId = userId
		viewingRecord.Duration = req.Duration
		if _, err := l.svcCtx.ViewingRecordModel.Insert(l.ctx, viewingRecord); err != nil {
			return response.Error(500, err.Error())
		}
	default:
		return response.Error(500, err.Error())
	}

	return nil
}
