package exam

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"management_system/service/api/internal/svc"
)

type ImportObjectiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportObjectiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportObjectiveLogic {
	return &ImportObjectiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportObjectiveLogic) ImportObjective() error {
	logx.Infof("userId: %v", l.ctx.Value("userId"))

	return nil
}
