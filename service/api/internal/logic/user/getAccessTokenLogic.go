package user

import (
	"context"
	"fmt"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	red "github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccessTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessTokenLogic {
	return &GetAccessTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccessTokenLogic) GetAccessToken() (resp *types.GetAccessTokenResp, err error) {
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	})
	// 先找有没有access_token缓存
	key := "access_token"
	value, err := l.svcCtx.RedisClient.Get(key)
	if err != nil && err != red.Nil {
		return nil, response.Error(500, "redis error:"+err.Error())
	}
	if value != "" {
		return &types.GetAccessTokenResp{
			AccessToken: value,
		}, nil
	}
	authResult, err := miniprogram.GetAuth().GetAccessToken()
	if err != nil {
		return nil, response.Error(100, fmt.Sprintf("获取接口调用凭证失败 err : %v , authResult: %v", err, authResult))
	}
	// 设置缓存,access_token过期时间为7200s（2小时）
	err = l.svcCtx.RedisClient.Setex(key, authResult, 7200)
	if err != nil {
		return nil, response.Error(500, "redis error:"+err.Error())
	}

	return &types.GetAccessTokenResp{
		AccessToken: authResult,
	}, nil
}
