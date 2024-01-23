package rescue

import (
	"context"
	"encoding/json"
	"strconv"

	"management_system/common/response"
	"management_system/service/api/internal/svc"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type DestroyRescueTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDestroyRescueTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DestroyRescueTokenLogic {
	return &DestroyRescueTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DestroyRescueTokenLogic) DestroyRescueToken() error {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 先在缓存里找，有就删除
	key := "management_system_rescue_token_" + strconv.Itoa(int(userId))
	rescueToken, err := l.svcCtx.RedisClient.Get(key)
	if err != nil && err != red.Nil {
		return  response.Error(500, "redis error:"+err.Error())
	}
	if rescueToken != "" {
		_, err = l.svcCtx.RedisClient.Del(key)
		if err != nil {
			return  response.Error(500, "redis error:"+err.Error())
		}
	}
	return nil
}
