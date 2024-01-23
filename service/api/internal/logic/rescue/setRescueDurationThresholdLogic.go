package rescue

import (
	"context"
	"encoding/json"
	"strconv"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetRescueDurationThresholdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetRescueDurationThresholdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetRescueDurationThresholdLogic {
	return &SetRescueDurationThresholdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetRescueDurationThresholdLogic) SetRescueDurationThreshold(req *types.SetRescueDurationThresholdReq) error {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 判断权限
	if userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}

	key := "rescue_duration_threshold"
	value := strconv.FormatInt(req.Threshold, 10)
	// 设置缓存（不设置过期时间）
	err = l.svcCtx.RedisClient.Set(key, value)
	if err != nil {
		return response.Error(500, "redis error:"+err.Error())
	}
	return nil
}
