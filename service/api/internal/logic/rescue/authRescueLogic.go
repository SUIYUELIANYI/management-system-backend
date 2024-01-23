package rescue

import (
	"context"
	"strconv"
	"strings"
	"time"

	"management_system/common/jwtx"
	"management_system/common/response"
	"management_system/common/tools"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	red "github.com/go-redis/redis/v8"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthRescueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthRescueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthRescueLogic {
	return &AuthRescueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthRescueLogic) AuthRescue(req *types.AuthReq) (resp *types.AuthResp, err error) {
	// 身份认证：发布救援信息需要管理员单独通过该接口用手机密码登录以获取专门用于发布救援信息的token。
	// 通过微信注册的账号默认密码为123456。且该token只能用一次，导入救援信息后立即失效。目前暂定组委会成员（role为43和44）才能通过认证。
	if len(strings.TrimSpace(req.Mobile)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, response.Error(100, "参数错误")
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "电话号码不存在")
	default:
		return nil, err
	}

	if userInfo.Password != tools.Md5ByString(req.Password) {
		return nil, response.Error(100, "用户密码不正确")
	}

	if userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}

	// 先在缓存中查找是否存在rescue_token
	key := "management_system_rescue_token_" + strconv.Itoa(int(userInfo.Id))
	value, err := l.svcCtx.RedisClient.Get(key)

	if err != nil && err != red.Nil {
		return nil, response.Error(500, "redis error:"+err.Error())
	}
	if value != "" {
		now := time.Now().Unix()
		expireTime, err := l.svcCtx.RedisClient.Ttl(key)
		if err != nil {
			return nil, response.Error(500, "redis error:"+err.Error())
		}

		return &types.AuthResp{
			AccessToken:  value,
			AccessExpire: now + int64(expireTime),
			RefreshAfter: now + int64(expireTime)/2,
		}, nil
	}
	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuthForRescue.AccessExpire
	jwtToken, err := jwtx.GetJwtToken(l.svcCtx.Config.JwtAuthForRescue.AccessSecret, now, l.svcCtx.Config.JwtAuthForRescue.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, response.Error(400, "Token生成失败:"+err.Error())
	}
	// ---end---

	// 设置缓存，用于后续删除，Set()不带时间，默认永久，Setex()带时间
	err = l.svcCtx.RedisClient.Setex(key, jwtToken, int(l.svcCtx.Config.JwtAuthForRescue.AccessExpire))
	if err != nil {
		return nil, response.Error(500, "redis error:"+err.Error())
	}

	return &types.AuthResp{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}
