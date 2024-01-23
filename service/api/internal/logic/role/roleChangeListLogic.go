package role

import (
	"context"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleChangeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleChangeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleChangeListLogic {
	return &RoleChangeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleChangeListLogic) RoleChangeList(req *types.RoleChangeInfoListReq) (resp *types.RoleChangeInfoListResp, err error) {
	// 管理员通过此接口+要查询成员的uid来获得成员的身份变动情况
	logx.Infof("userId: %v", l.ctx.Value("userId"))

	RoleChangeInfoListResp, err := l.svcCtx.RoleChangeModel.FindAllByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var result []types.RoleChangeInfo // 将models.RoleChangeInfo转换成types.RoleChangeInfo
	if len(RoleChangeInfoListResp) > 0 {
		for _, roleschangemodel := range RoleChangeInfoListResp {
			var typeRoleChangeInfo types.RoleChangeInfo
			typeRoleChangeInfo.CreateTime = roleschangemodel.CreateTime.Format("2006-01-02 15:04:05")
			typeRoleChangeInfo.Id = roleschangemodel.Id
			typeRoleChangeInfo.NewRole = roleschangemodel.NewRole
			typeRoleChangeInfo.OldRole = roleschangemodel.OldRole
			typeRoleChangeInfo.OperatorId = roleschangemodel.OperatorId
			typeRoleChangeInfo.UserId = roleschangemodel.UserId
			result = append(result, typeRoleChangeInfo)
		}
	}
	return &types.RoleChangeInfoListResp{
		List: result,
	}, nil
}
