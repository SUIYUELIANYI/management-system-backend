package user

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserInfoLogic {
	return &EditUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserInfoLogic) EditUserInfo(req *types.EditUserInfoReq) error {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return response.Error(500, err.Error())
	}

	userInfo.Mobile = req.Mobile
	userInfo.Username = req.Username
	userInfo.Sex = req.Sex
	userInfo.Address = req.Address
	userInfo.Birthday = req.Birthday
	userInfo.Email = req.Email

	err = l.svcCtx.UserModel.Update(l.ctx, userInfo)
	if err != nil {
		return response.Error(500, err.Error())
	}
	return response.Success("修改个人信息成功！")
}
