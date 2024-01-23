package user

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RescueFrequencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRescueFrequencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RescueFrequencyLogic {
	return &RescueFrequencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RescueFrequencyLogic) RescueFrequency() (resp *types.RescueFrequencyResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	// 判断权限（只有见习队员以上的人才能有救援时长）
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 查看救援过程次数（rescue_process表）
	rescueFrequency, err := l.svcCtx.RescueProcessModel.CountRescueFrequency(l.ctx, userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	return &types.RescueFrequencyResp{
		RescueFrequency: rescueFrequency,
	}, nil
}
