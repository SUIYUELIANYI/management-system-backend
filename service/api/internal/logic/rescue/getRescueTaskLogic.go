package rescue

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRescueTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRescueTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRescueTaskLogic {
	return &GetRescueTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRescueTaskLogic) GetRescueTask(req *types.GetRescueTaskReq) (resp *types.GetRescueTaskResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 获取救援任务就是返回当前用户认领的救援对象
	rescueTargets, err := l.svcCtx.RescueTargetModel.FindAllByUserIdWithPage(l.ctx, l.svcCtx.RescueTargetModel.RowBuilder(), req.Page, req.PageSize, userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.RescueTarget
	if len(rescueTargets) > 0 {
		for _, rescueTargetModel := range rescueTargets {
			var typeRescueTarget types.RescueTarget
			_ = copier.Copy(&typeRescueTarget, rescueTargetModel)
			typeRescueTarget.CreateTime = rescueTargetModel.CreateTime.Format("2006-01-02 15:04:05")
			typeRescueTarget.UpdateTime = rescueTargetModel.UpdateTime.Format("2006-01-02 15:04:05")
			typeRescueTarget.StartTime = rescueTargetModel.StartTime.Format("2006-01-02 15:04:05")
			if rescueTargetModel.EndTime.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
				typeRescueTarget.EndTime = "0000-00-00 00:00:00"
			} else {
				typeRescueTarget.EndTime = rescueTargetModel.EndTime.Format("2006-01-02 15:04:05")
			}
			if typeRescueTarget.RescueTeacher1Id != 0 {
				typeRescueTarget.RescueTeacher1Name, err = l.svcCtx.UserModel.FindUsernameById(l.ctx, typeRescueTarget.RescueTeacher1Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
				typeRescueTarget.RescueTeacher1Role, err = l.svcCtx.UserModel.FindRoleById(l.ctx, typeRescueTarget.RescueTeacher1Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
			}
			if typeRescueTarget.RescueTeacher2Id != 0 {
				typeRescueTarget.RescueTeacher2Name, err = l.svcCtx.UserModel.FindUsernameById(l.ctx, typeRescueTarget.RescueTeacher2Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
				typeRescueTarget.RescueTeacher2Role, err = l.svcCtx.UserModel.FindRoleById(l.ctx, typeRescueTarget.RescueTeacher2Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
			}
			if typeRescueTarget.RescueTeacher3Id != 0 {
				typeRescueTarget.RescueTeacher3Name, err = l.svcCtx.UserModel.FindUsernameById(l.ctx, typeRescueTarget.RescueTeacher3Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
				typeRescueTarget.RescueTeacher3Role, err = l.svcCtx.UserModel.FindRoleById(l.ctx, typeRescueTarget.RescueTeacher3Id)
				if err != nil {
					return nil, response.Error(500, err.Error())
				}
			}
			// 添加昵称
			rescueInfo, err := l.svcCtx.RescueInfoModel.FindOneByWeiboAddress(l.ctx, rescueTargetModel.WeiboAddress)
			if err != nil {
				return nil, response.Error(500, err.Error())
			}
			typeRescueTarget.Nickname = rescueInfo.Nickname

			list = append(list, typeRescueTarget)
		}
	}

	return &types.GetRescueTaskResp{
		List: list,
	}, nil
}
