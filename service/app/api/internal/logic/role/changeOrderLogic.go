package role

import (
	"context"

	"management-system/service/app/api/internal/svc"
	"management-system/service/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOrderLogic {
	return &ChangeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeOrderLogic) ChangeOrder(req *types.ChangeRoleReq) (resp *types.ChangeRoleResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致
	return
}
