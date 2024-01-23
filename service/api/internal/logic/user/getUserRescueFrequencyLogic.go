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

type GetUserRescueFrequencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRescueFrequencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRescueFrequencyLogic {
	return &GetUserRescueFrequencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRescueFrequencyLogic) GetUserRescueFrequency(req *types.GetUserRescueFrequencyReq) (resp *types.GetUserRescueFrequencyResp, err error) {
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

	_, err = l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "查询用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	// 查看救援过程次数（rescue_process表）
	rescueFrequency, err := l.svcCtx.RescueProcessModel.CountRescueFrequency(l.ctx, req.UserId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	return &types.GetUserRescueFrequencyResp{
		RescueFrequency: rescueFrequency,
	}, nil
}
