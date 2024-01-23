package file

import (
	"context"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetViewingRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetViewingRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetViewingRecordLogic {
	return &GetViewingRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetViewingRecordLogic) GetViewingRecord(req *types.GetViewingRecordReq) (resp *types.GetViewingRecordResp, err error) {
	ViewingRecords, err := l.svcCtx.ViewingRecordModel.FindAllWithPage(l.ctx, l.svcCtx.ViewingRecordModel.RowBuilder(), req.Page, req.PageSize)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	var list []types.ViewingRecord
	if len(ViewingRecords) > 0 {
		for _, viewingrecordmodel := range ViewingRecords {
			var typeViewingRecord types.ViewingRecord
			_ = copier.Copy(&typeViewingRecord, viewingrecordmodel)
			typeViewingRecord.CreateTime = viewingrecordmodel.CreateTime.Format("2006-01-02 15:04:05")
			typeViewingRecord.UpdateTime = viewingrecordmodel.UpdateTime.Format("2006-01-02 15:04:05")
			if typeViewingRecord.UserId != 0 {
				typeViewingRecord.Username, err = l.svcCtx.UserModel.FindUsernameById(l.ctx, typeViewingRecord.UserId)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
			}
			if typeViewingRecord.FileId != 0 {
				typeViewingRecord.FileName, err = l.svcCtx.FileModel.FindFileNameById(l.ctx, typeViewingRecord.FileId)
			}
			list = append(list, typeViewingRecord)
		}
	}
	return &types.GetViewingRecordResp{
		List: list,
	}, nil
}
