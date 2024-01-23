package role

import (
	"context"
	"encoding/json"
	"fmt"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ChangeRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeRoleLogic {
	return &ChangeRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeRoleLogic) ChangeRole(req *types.ChangeRoleReq) error {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	// 被更改身份的用户信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return response.Error(500, err.Error())
	}
	// 操作人信息
	operatorInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return response.Error(500, err.Error())
	}
	// 更改身份的权限仅限 43-组委会成员 44-组委会主任
	if operatorInfo.Role == 43 || operatorInfo.Role == 44 {
		// 创建role_change表数据
		roleChange := new(models.RoleChange)
		roleChange.UserId = req.UserId
		roleChange.OperatorId = userId
		roleChange.NewRole = req.NewRole
		roleChange.OldRole = userInfo.Role

		if req.NewRole == 2 && userInfo.Role >= 3 { // 特殊情况：如果此次修改身份是从见习队员（以上）改为岗前培训，需要将所有的考试成绩进行删除
			// 事务：同时插入relo_change表数据和更新user表数据以及删除用户所有成绩
			fmt.Println("测试点1")
			if err := l.svcCtx.RoleChangeModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
					return err
				}
				// 修改userInfo数据
				userInfo.Role = req.NewRole
				if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
					return err
				}
				if err := l.svcCtx.ObjectiveExamModel.DeleteByUserId(l.ctx, req.UserId); err != nil {
					return err
				}
				if err := l.svcCtx.SubjectiveExamModel.DeleteByUserId(l.ctx, req.UserId); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return response.Error(500, err.Error())
			}
		} else { // 一般情况
			// 事务：同时插入relo_change表数据和更新user表数据
			if err := l.svcCtx.RoleChangeModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
					return err
				}
				// 修改userInfo数据
				userInfo.Role = req.NewRole
				if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return response.Error(500, err.Error())
			}
		}
	} else {
		return response.Error(403, "权限不够")
	}

	return response.Success("更改身份成功")
}
