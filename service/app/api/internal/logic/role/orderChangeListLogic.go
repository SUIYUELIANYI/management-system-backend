package role

import (
	"context"

	"management-system/service/app/api/internal/svc"
	"management-system/service/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderChangeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderChangeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderChangeListLogic {
	return &OrderChangeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderChangeListLogic) OrderChangeList(req *types.RoleChangeInfoListReq) (resp *types.RoleChangeInfoListResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致

	return
}
