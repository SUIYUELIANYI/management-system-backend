package user

import (
	"context"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.UserInfoResp, err error) {
	userInfoResp, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "该用户不存在！")
	default:
		return nil, response.Error(500, err.Error())
	}

	var userInfo types.User

	err = copier.Copy(&userInfo, userInfoResp)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	return &types.UserInfoResp{
		UserInfo: userInfo,
	}, nil
}
