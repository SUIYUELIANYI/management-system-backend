package rescue

import (
	"context"
	"database/sql"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignLogic {
	return &SignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignLogic) Sign(req *types.SignReq) error {
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
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	// 查看该救援对象是不是当前用户的
	rescueTarget, err := l.svcCtx.RescueTargetModel.FindOneByIdUserId(l.ctx, req.RescueTargetId, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "您没有认领该救援对象！")
	default:
		return response.Error(500, err.Error())
	}
	// 判断是否完成救援
	if rescueTarget.Status == 2 {
		return response.Error(100, "当前救援对象已完成救援！")
	}
	// 判断是否签过字
	_, err = l.svcCtx.SignatureModel.FindOneByRescueTeacherIdRescueTargetId(l.ctx, req.RescueTargetId, userId)
	switch err {
	case nil:
		return response.Error(100, "您已完成签字！")
	case models.ErrNotFound:
	default:
		return response.Error(500, err.Error())
	}

	signature := new(models.Signature)
	signature.RescueTeacherId = userId
	signature.RescueTargetId = rescueTarget.Id
	signature.Image = sql.NullString{String: req.Image, Valid: true}
	// 判断逻辑：如果当前是最后一位签字的，则该救援对象的状态更新为2，否则不变，具体操作是去签字表里搜索rescueTargeId有几条数据，然后和rescueTargetId有几位老师相比
	// 统计救援对象老师人数
	teacherNumber := 0
	if rescueTarget.RescueTeacher1Id != 0 {
		teacherNumber += 1
	}
	if rescueTarget.RescueTeacher2Id != 0 {
		teacherNumber += 1
	}
	if rescueTarget.RescueTeacher3Id != 0 {
		teacherNumber += 1
	}
	// 已经签过字的数量
	count, err := l.svcCtx.SignatureModel.CountByRescueTargetId(l.ctx, rescueTarget.Id)
	if err != nil {
		return response.Error(500, err.Error())
	}
	// 本次签字还没有录入数据库，所以当签字数量=老师人数-1时，说明这是最后一位签字的老师，需要进行救援对象的状态更新
	if count == int64(teacherNumber)-1 {
		// 事务包含更新救援对象状态和插入签字信息
		if err := l.svcCtx.RescueTargetModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			rescueTarget.Status = 2
			if err := l.svcCtx.RescueTargetModel.Update(l.ctx, rescueTarget); err != nil {
				return err
			}
			if _, err := l.svcCtx.SignatureModel.Insert(l.ctx, signature); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return response.Error(500, err.Error())
		}
	} else {
		if _, err := l.svcCtx.SignatureModel.Insert(l.ctx, signature); err != nil {
			return response.Error(500, err.Error())
		}
	}
	return nil
}
