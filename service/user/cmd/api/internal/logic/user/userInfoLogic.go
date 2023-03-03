package user

import (
	"context"
	"encoding/json"

	"management-system/service/user/cmd/api/internal/svc"
	"management-system/service/user/cmd/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {

	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致
	var userId int64
	if jsonUid, ok := l.ctx.Value("userId").(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			userId = int64Uid
		} else {
			logx.WithContext(l.ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}

	userInfoResp, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	var userInfo types.User
	_ = copier.Copy(&userInfo, userInfoResp)

	return &types.UserInfoResp{
		UserInfo: userInfo,
	}, nil
}
