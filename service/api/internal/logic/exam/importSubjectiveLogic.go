package exam

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"management_system/service/api/internal/svc"
)

type ImportSubjectiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportSubjectiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportSubjectiveLogic {
	return &ImportSubjectiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportSubjectiveLogic) ImportSubjective() error {
	logx.Infof("userId: %v", l.ctx.Value("userId"))

	return nil
}
