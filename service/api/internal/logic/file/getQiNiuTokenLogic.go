package file

import (
	"context"
	"encoding/json"

	"management_system/common/qiniu"
	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQiNiuTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQiNiuTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQiNiuTokenLogic {
	return &GetQiNiuTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQiNiuTokenLogic) GetQiNiuToken() (resp *types.GetQiNiuTokenResp, err error) {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以上传视频）
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(401, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	if userInfo.Role != 5 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 获取七牛云token
	return &types.GetQiNiuTokenResp{
		Token: qiniu.GetToken(),
	}, nil
}
