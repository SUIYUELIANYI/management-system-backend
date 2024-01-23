package user

import (
	"context"
	"fmt"
	"time"

	"management_system/common/jwtx"
	"management_system/common/response"
	"management_system/common/tools"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type WXMiniAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWXMiniAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WXMiniAuthLogic {
	return &WXMiniAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WXMiniAuthLogic) WXMiniAuth(req *types.WXMiniAuthReq) (resp *types.WXMiniAuthResp, err error) {
	// 为了保证微信账号和系统内部注册账号一致性，先检查用户授权表中是否有系统内部使用同一电话注册
	ifExitSamePhone, err := l.svcCtx.UserAuthModel.FindOneByAuthTypeAuthKey(l.ctx, models.UserAuthTypeSystem, req.PhoneNumber)
	if err != nil && err != models.ErrNotFound {
		return nil, response.Error(500, err.Error())
	}
	if ifExitSamePhone != nil {
		return nil, response.Error(100, "该手机号已注册")
	}
	// 1、Wechat-Mini
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	})
	authResult, err := miniprogram.GetAuth().Code2Session(req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, response.Error(100, fmt.Sprintf("发起授权请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult))
	}
	// 2、Parsing WeChat-Mini return data
	// GetEncryptor为小程序加解密函数，但由于微信官方的限制，该接口已弃用，已经无法解析到有效数据
	/* userData, err := miniprogram.GetEncryptor().Decrypt(authResult.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil, response.Error(100, fmt.Sprintf("解析数据失败 req : %+v , err: %v , authResult:%+v ", req, err, authResult))
	} */
	fmt.Println("OpenID", authResult.OpenID)
	fmt.Println("session_key", authResult.SessionKey)
	fmt.Println("unionid", authResult.UnionID)
	// fmt.Println("头像地址", userData.AvatarURL)
	// fmt.Println("OpenID", userData.OpenID)
	// fmt.Println("电话", userData.PhoneNumber)
	// fmt.Println("access_token", authResult.SessionKey)
	// 3、bind user or login.
	var userId int64
	userAuthInfo, err := l.svcCtx.UserAuthModel.FindOneByAuthTypeAuthKey(l.ctx, models.UserAuthTypeSmallWX, authResult.OpenID)
	if err != nil && err != models.ErrNotFound {
		return nil, response.Error(500, err.Error())
	}

	if userAuthInfo == nil || userAuthInfo.Id == 0 {
		fmt.Println("测试点1")
		// 如果根据AuthType和Authkey查出的用户数据信息为空，说明是第一次用这种方式登录，需要将数据插入数据库
		// 创建users和usersauth数据
		user := new(models.User)
		user.Sex = 0
		user.Role = 0
		user.Mobile = req.PhoneNumber
		user.Password = tools.Md5ByString("123456") // 微信登录的账户设置默认密码为123456
		sqlResult, err := l.svcCtx.UserModel.Insert(l.ctx, user)
		if err != nil {
			return nil, response.Error(500, err.Error())
		}
		userId, err = sqlResult.LastInsertId()
		if err != nil {
			return nil, response.Error(500, err.Error())
		}
		userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
		switch err {
		case nil:
		case models.ErrNotFound:
			return nil, response.Error(100, "该手机号未注册")
		default:
			return nil, response.Error(500, err.Error())
		}

		usersAuth := models.UserAuth{
			AuthKey:  authResult.OpenID,
			AuthType: models.UserAuthTypeSmallWX,
			UserId:   userInfo.Id,
		}

		if _, err := l.svcCtx.UserAuthModel.Insert(l.ctx, &usersAuth); err != nil {
			return nil, response.Error(500, err.Error())
		}

		// 生成token
		now := time.Now().Unix()
		accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
		jwtToken, err := jwtx.GetJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userInfo.Id)
		if err != nil {
			return nil, response.Error(400, "Token生成失败:"+err.Error())
		}

		return &types.WXMiniAuthResp{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
			Role:         userInfo.Role,
		}, nil
	} else {
		userId = userAuthInfo.UserId

		// 生成token
		now := time.Now().Unix()
		accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
		jwtToken, err := jwtx.GetJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userId)
		if err != nil {
			return nil, response.Error(400, "Token生成失败:"+err.Error())
		}

		userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
		switch err {
		case nil:
		case models.ErrNotFound:
			return nil, response.Error(100, "该手机号未注册")
		default:
			return nil, response.Error(500, err.Error())
		}

		return &types.WXMiniAuthResp{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
			Role:         userInfo.Role,
		}, nil
	}
}
