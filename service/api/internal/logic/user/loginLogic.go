package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"management_system/common/jwtx"
	"management_system/common/response"
	"management_system/common/tools"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

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
		return nil, response.Error(100, "参数错误")
	}

	userInfo, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "电话号码不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	if userInfo.Password != tools.Md5ByString(req.Password) {
		return nil, response.Error(100, "用户密码不正确")
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	jwtToken, err := jwtx.GetJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userInfo.Id)
	fmt.Println(userInfo.Id)
	if err != nil {
		return nil, response.Error(400, "Token生成失败:"+err.Error())
	}
	// ---end---

	return &types.LoginResp{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
		Role:         userInfo.Role,
	}, nil
}
