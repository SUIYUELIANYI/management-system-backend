package rescue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimRescueTargetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimRescueTargetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimRescueTargetLogic {
	return &ClaimRescueTargetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimRescueTargetLogic) ClaimRescueTarget(req *types.ClaimRescueTargetReq) error {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 查看登录用户
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 查看救援对象是否是待救援或者救援中
	rescueTarget, err := l.svcCtx.RescueTargetModel.FindOne(l.ctx, req.RescueTargetId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前救援对象不存在")
	default:
		return response.Error(500, err.Error())
	}
	if rescueTarget.Status == 2 {
		return response.Error(100,"当前救援对象已完成救援！")
	}
	// 查看登录用户救援对象是否超过6个
	count, err := l.svcCtx.RescueTargetModel.CountClaimNumber(l.ctx, userId)
	if err != nil {
		return response.Error(500, err.Error())
	}
	fmt.Println(count)
	if count == 6 {
		return response.Error(100, "您的救援对象数量已达到最大值！")
	}
	// 查看该救援对象是否已经被当前用户认领
	ifClaim, err := l.svcCtx.RescueTargetModel.FindOneByIdUserId(l.ctx, req.RescueTargetId, userId)
	switch err {
	case nil:
		if ifClaim != nil {
			return response.Error(100, "您已经认领了该救援对象！")
		}
	case models.ErrNotFound:
	default:
		return response.Error(500, err.Error())
	}

	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	// 认领救援对象（三个救援老师不能都是见习队员，每个救援老师最多认领6个队员）
	// 具体逻辑
	// （如果当前志愿者身份是41-核心队员，42-区域负责人，43-组委会成员，44-组委会主任，则按顺序赋给teacher1/2/3）
	// （如果当前志愿者身份是3-见习队员 40-普通队员，则按顺序赋给teacher2/3）
	if userInfo.Role == 3 || userInfo.Role == 40 {
		if rescueTarget.Status == 0 { // 如果当前救援对象状态为待救援
			rescueTarget.StartTime = time.Now()
			rescueTarget.RescueTeacher2Id = userInfo.Id
			rescueTarget.Status = 1
		} else if rescueTarget.Status == 1 { // 如果当前救援对象状态为救援中
			if rescueTarget.RescueTeacher2Id == 0 {
				rescueTarget.RescueTeacher2Id = userInfo.Id
			} else if rescueTarget.RescueTeacher3Id == 0 {
				rescueTarget.RescueTeacher3Id = userInfo.Id
			} else if rescueTarget.RescueTeacher1Id == 0 {
				return response.Error(100, "当前救援对象需要核心队员身份以上才能认领！")
			} else {
				return response.Error(100, "救助老师人数已满！")
			}
		} else if rescueTarget.Status == 2 { // 如果当前救援对象状态为已救援
			return response.Error(100, "当前救援对象已完成救援！")
		}
		if err := l.svcCtx.RescueTargetModel.Update(l.ctx, rescueTarget); err != nil {
			return response.Error(500, err.Error())
		}
	} else {
		if rescueTarget.Status == 0 {
			rescueTarget.StartTime = time.Now()
			rescueTarget.RescueTeacher1Id = userInfo.Id
			rescueTarget.Status = 1
		} else if rescueTarget.Status == 1 {
			if rescueTarget.RescueTeacher1Id == 0 {
				rescueTarget.RescueTeacher1Id = userInfo.Id
			} else if rescueTarget.RescueTeacher2Id == 0 {
				rescueTarget.RescueTeacher2Id = userInfo.Id
			} else if rescueTarget.RescueTeacher3Id == 0 {
				rescueTarget.RescueTeacher3Id = userInfo.Id
			} else {
				return response.Error(100, "救助老师人数已满！")
			}
		} else if rescueTarget.Status == 2 {
			return response.Error(100, "当前救援对象已完成救援！")
		}
		if err := l.svcCtx.RescueTargetModel.Update(l.ctx, rescueTarget); err != nil {
			return response.Error(500, err.Error())
		}
	}
	return response.Success("认领成功！")
}
