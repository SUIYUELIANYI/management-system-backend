package user

import (
	"context"

	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportExcelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportExcelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportExcelLogic {
	return &ImportExcelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportExcelLogic) ImportExcel() (resp *types.ImportExcelResp, err error) {
	// todo: add your logic here and delete this line

	return
}
