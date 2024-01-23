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

type GetTrainingUserInfoByMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTrainingUserInfoByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTrainingUserInfoByMobileLogic {
	return &GetTrainingUserInfoByMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTrainingUserInfoByMobileLogic) GetTrainingUserInfoByMobile(req *types.GetTrainingUserInfoByMobileReq) (resp *types.GetTrainingUserInfoByMobileResp, err error) {
	userInfoResp, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "该用户不存在！")
	default:
		return nil, response.Error(500, err.Error())
	}

	if userInfoResp.Role != 2 {
		return nil, response.Error(100, "岗前培训成员不包含该用户！")
	}

	var userInfo types.User

	err = copier.Copy(&userInfo, userInfoResp)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	return &types.GetTrainingUserInfoByMobileResp{
		UserInfo: userInfo,
	}, nil
}
