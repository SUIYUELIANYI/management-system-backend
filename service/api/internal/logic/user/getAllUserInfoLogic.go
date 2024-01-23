package user

import (
	"context"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserInfoLogic {
	return &GetAllUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllUserInfoLogic) GetAllUserInfo(req *types.GetAllUserInfoReq) (resp *types.GetAllUserInfoResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))

	userInfos, err := l.svcCtx.UserModel.FindAllWithPage(l.ctx, l.svcCtx.UserModel.RowBuilder(), req.Page, req.PageSize)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.User
	if len(userInfos) > 0 {
		for _, usermodel := range userInfos {
			var typeUserInfo types.User
			typeUserInfo.Id = usermodel.Id
			typeUserInfo.Mobile = usermodel.Mobile
			typeUserInfo.Username = usermodel.Username
			typeUserInfo.Sex = usermodel.Sex
			typeUserInfo.Email = usermodel.Email
			typeUserInfo.Role = usermodel.Role
			typeUserInfo.Avatar = usermodel.Avatar
			typeUserInfo.Address = usermodel.Address
			typeUserInfo.Birthday = usermodel.Birthday
			list = append(list, typeUserInfo)
		}
	}
	return &types.GetAllUserInfoResp{
		List: list,
	}, nil
}
