package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"management-system/service/user/cmd/api/internal/svc"
	"management-system/service/user/cmd/api/internal/types"
	"management-system/service/user/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if len(strings.TrimSpace(req.Mobile)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errors.New("参数错误")
	}

	userInfo, err := l.svcCtx.UsersModel.FindOneByMobile(l.ctx, req.Mobile)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("电话未注册")
	default:
		return nil, err
	}

	if userInfo.Password != req.Password {
		return nil, errors.New("用户密码不正确")
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userInfo.Id)
	fmt.Println(userInfo.Id)
	if err != nil {
		return nil, err
	}
	// ---end---

	return &types.LoginResp{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
