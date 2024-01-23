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

type GetUserRescueDurationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRescueDurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRescueDurationLogic {
	return &GetUserRescueDurationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRescueDurationLogic) GetUserRescueDuration(req *types.GetUserRescueDurationReq) (resp *types.GetUserRescueDurationResp, err error) {
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
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 判断查询的userId是否存在
	_, err = l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "查询用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	rescueProcesses, err := l.svcCtx.RescueProcessModel.FindAllByRescueTeacherId(l.ctx, l.svcCtx.RescueProcessModel.RowBuilder(), req.UserId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	durations := []string{}
	if len(rescueProcesses) > 0 {
		for _, rescueProcessModel := range rescueProcesses {
			durations = append(durations, rescueProcessModel.Duration)
		}
	}

	return &types.GetUserRescueDurationResp{
		RescueDuration: sumDurations(durations),
	}, nil
}
