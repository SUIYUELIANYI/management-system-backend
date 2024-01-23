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

type GetRescueTargetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRescueTargetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRescueTargetInfoLogic {
	return &GetRescueTargetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRescueTargetInfoLogic) GetRescueTargetInfo(req *types.GetRescueTargetInfoReq) (resp *types.GetRescueTargetInfoResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	_, err = l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	rescueTargetInfo, err := l.svcCtx.RescueTargetModel.FindOne(l.ctx, req.RescueTargetId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "救援对象不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	var typeRescueTarget types.RescueTarget
	_ = copier.Copy(&typeRescueTarget, rescueTargetInfo)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	typeRescueTarget.CreateTime = rescueTargetInfo.CreateTime.Format("2006-01-02 15:04:05")
	typeRescueTarget.UpdateTime = rescueTargetInfo.UpdateTime.Format("2006-01-02 15:04:05")
	typeRescueTarget.StartTime = rescueTargetInfo.StartTime.Format("2006-01-02 15:04:05")
	if rescueTargetInfo.EndTime.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		typeRescueTarget.EndTime = "0000-00-00 00:00:00"
	} else {
		typeRescueTarget.EndTime = rescueTargetInfo.EndTime.Format("2006-01-02 15:04:05")
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
	rescueInfo, err := l.svcCtx.RescueInfoModel.FindOneByWeiboAddress(l.ctx, rescueTargetInfo.WeiboAddress)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	typeRescueTarget.Nickname = rescueInfo.Nickname

	return &types.GetRescueTargetInfoResp{
		RescueTargetInfo: typeRescueTarget,
	}, nil
}
