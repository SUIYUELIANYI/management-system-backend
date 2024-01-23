package user

import (
	"context"

	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

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
	// todo: add your logic here and delete this line

	return nil
}
