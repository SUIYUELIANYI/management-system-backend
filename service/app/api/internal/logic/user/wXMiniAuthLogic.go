package user

import (
	"context"
	"fmt"
	"time"

	"management-system/common/xerr"
	"management-system/service/app/api/internal/svc"
	"management-system/service/app/api/internal/types"
	"management-system/service/app/models"
	"management-system/service/user/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/cache"
	wechat "github.com/silenceper/wechat/v2"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ErrWxMiniAuthFailError error
var ErrWxMiniAuthFailError = xerr.NewErrMsg("wechat mini auth fail")

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
	//1、Wechat-Mini
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	})
	fmt.Println(l.svcCtx.Config.WxMiniConf.AppId)
	fmt.Println(l.svcCtx.Config.WxMiniConf.AppSecret)

	authResult, err := miniprogram.GetAuth().Code2Session(req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "发起授权请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult)
	}
	//2、Parsing WeChat-Mini return data
	userData, err := miniprogram.GetEncryptor().Decrypt(authResult.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "解析数据失败 req : %+v , err: %v , authResult:%+v ", req, err, authResult)
	}
	fmt.Println("头像地址", userData.AvatarURL)
	fmt.Println("OpenID", userData.OpenID)
	fmt.Println("电话", userData.PhoneNumber)
	fmt.Println("access_token", authResult.SessionKey)

	GetPhoneNumberResponse, err := miniprogram.GetAuth().GetPhoneNumber(req.Code)
	if err != nil {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "解析电话请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult)
	}

	userData.PhoneNumber = GetPhoneNumberResponse.PhoneInfo.PhoneNumber

	//3、bind user or login.
	var userId int64
	// 如果是下面这行报错，请查看你的UsersAuthModel有没有添加到svc中的函数中去
	userAuthInfo, err := l.svcCtx.UsersAuthModel.FindOneByAuthTypeAuthKey(l.ctx, authResult.OpenID, model.UserAuthTypeSmallWX)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "rpc call userAuthByAuthKey err : %v , authResult : %+v", err, authResult)
	}

	if userAuthInfo == nil || userAuthInfo.Id == 0 {
		//if ptr == nil {
		// 如果根据AuthType和Authkey查出的用户数据信息为空，说明是第一次用这种方式登录，需要将数据插入数据库
		//bind user.

		//Wechat-Mini Decrypted data
		mobile := userData.PhoneNumber
		userName := fmt.Sprintf("用户%s", mobile[7:])

		// 创建users和usersauth数据
		users := new(models.Users)
		users.Mobile = mobile
		users.Username = userName

		if _, err := l.svcCtx.UsersModel.Insert(l.ctx, users); err != nil {
			return nil, nil
		}

		usersAuth := models.UsersAuth{
			AuthKey:  authResult.OpenID,
			AuthType: models.UserAuthTypeSmallWX,
		}

		if _, err := l.svcCtx.UsersAuthModel.Insert(l.ctx, &usersAuth); err != nil {
			return nil, nil
		}
		userInfo, err := l.svcCtx.UsersModel.FindOneByMobile(l.ctx, mobile)
		switch err {
		case nil:
		case models.ErrNotFound:
			return nil, errors.New("电话未注册")
		default:
			return nil, err
		}

		// 生成token
		now := time.Now().Unix()
		accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
		jwtToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userInfo.Id)
		fmt.Println(userInfo.Id)
		if err != nil {
			return nil, err
		}

		return &types.WXMiniAuthResp{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		}, nil
	} else {
		userId = userAuthInfo.UserId
		// 生成token
		now := time.Now().Unix()
		accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
		jwtToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userId)
		if err != nil {
			return nil, err
		}

		return &types.WXMiniAuthResp{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		}, nil
	}
}

func (l *WXMiniAuthLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
