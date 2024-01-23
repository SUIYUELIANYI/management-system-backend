package user

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

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

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId")) // 这里的key和生成jwt token时传入的key一致
	var userId int64
	if jsonUid, ok := l.ctx.Value("userId").(json.Number); ok { // 通过 l.ctx.Value("uid") 可获取 jwt token 中自定义的参数
		if int64Uid, err := jsonUid.Int64(); err == nil {
			userId = int64Uid
		} else {
			return nil, response.Error(401, err.Error()) // 这里并不是token错误，在进入这个接口前，token就已经被解析了
		}
	}

	userInfoResp, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
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
