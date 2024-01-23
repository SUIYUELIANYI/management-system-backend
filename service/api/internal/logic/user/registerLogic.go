package user

import (
	"context"
	"time"

	"management_system/common/jwtx"
	"management_system/common/response"
	"management_system/common/tools"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 先根据手机号查找用户是否注册
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil && err != models.ErrNotFound {
		return nil, response.Error(500, err.Error())
	}
	if user != nil {
		return nil, response.Error(100, "该手机号已注册")
	}
	// 用户表创建数据
	user = new(models.User)
	user.Mobile = req.Mobile
	user.Password = tools.Md5ByString(req.Password)
	user.Sex = 0
	user.Role = 0
	if _, err := l.svcCtx.UserModel.Insert(l.ctx, user); err != nil {
		return nil, response.Error(500, err.Error())
	}
	// 获取user_id
	userInfo, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	// 用户授权表创建数据
	userAuth := new(models.UserAuth)
	userAuth.AuthKey = req.Mobile
	userAuth.AuthType = models.UserAuthTypeSystem
	userAuth.UserId = userInfo.Id
	if _, err := l.svcCtx.UserAuthModel.Insert(l.ctx, userAuth); err != nil {
		return nil, response.Error(500, err.Error())
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	jwtToken, err := jwtx.GetJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userInfo.Id)
	if err != nil {
		return nil, response.Error(400, "Token生成失败:"+err.Error())
	}
	// ---end---

	return &types.RegisterResp{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
		Role:         0,
	}, nil
}
