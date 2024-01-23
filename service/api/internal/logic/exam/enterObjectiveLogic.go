package exam

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type EnterObjectiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnterObjectiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnterObjectiveLogic {
	return &EnterObjectiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnterObjectiveLogic) EnterObjective(req *types.EnterGradeReq) error {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 确保存在该用户
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "该用户不存在！")
	default:
		return response.Error(500, err.Error())
	}
	if userInfo.Role == -1 {
		return response.Error(100, "该用户为待处理人员，无法操作！")
	}
	if userInfo.DelState == 1 {
		return response.Error(100, "该用户已删除！")
	}
	// 确保该用户身份是岗前培训
	if userInfo.Role != 2 {
		return response.Error(100, "该用户身份非岗前培训！")
	}
	// 确保该用户没有通过该考试
	ifExitPassRecord, err := l.svcCtx.ObjectiveExamModel.FindOneByUserIdResult(l.ctx, req.UserId, 1)
	if err != nil && err != models.ErrNotFound {
		return response.Error(500, err.Error())
	}
	if ifExitPassRecord != nil {
		return response.Error(400, "该用户客观题已合格！")
	}
	// 确保该条记录不重复
	ifExitSameRecord, err := l.svcCtx.ObjectiveExamModel.FindOneByUserIdTime(l.ctx, req.UserId, req.Time)
	if err != nil && err != models.ErrNotFound {
		return response.Error(500, err.Error())
	}
	if ifExitSameRecord != nil {
		return response.Error(400, "不能录入相同成绩！")
	}
	if req.Result != 0 && req.Result != 1 {
		return response.Error(400, "请输入正确的成绩格式：0-不合格，1-合格！")
	}
	if req.Time == "" {
		return response.Error(400, "日期不能为空！")
	}

	if req.Result == 0 { // 不合格的成绩直接录入，不会影响身份变动
		// 插入客观题成绩
		objectiveExam := new(models.ObjectiveExam)
		objectiveExam.UserId = req.UserId
		objectiveExam.Result = req.Result
		objectiveExam.Time = req.Time
		if _, err := l.svcCtx.ObjectiveExamModel.Insert(l.ctx, objectiveExam); err != nil {
			return response.Error(500, err.Error())
		}
	} else {
		ifBothPassRecord, err := l.svcCtx.SubjectiveExamModel.FindOneByUserIdResult(l.ctx, req.UserId, 1)
		if err != nil && err != models.ErrNotFound {
			return response.Error(500, err.Error())
		}
		if ifBothPassRecord == nil { // 如果主观题没有合格，直接插入客观题成绩
			// 插入客观题成绩
			objectiveExam := new(models.ObjectiveExam)
			objectiveExam.UserId = req.UserId
			objectiveExam.Result = req.Result
			objectiveExam.Time = req.Time
			if _, err := l.svcCtx.ObjectiveExamModel.Insert(l.ctx, objectiveExam); err != nil {
				return response.Error(500, err.Error())
			}
		} else { // 如果主观题已经合格，该事务包括插入客观题成绩，更新用户身份，插入身份变动数据
			if err := l.svcCtx.ObjectiveExamModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				// 插入客观题成绩
				objectiveExam := new(models.ObjectiveExam)
				objectiveExam.UserId = req.UserId
				objectiveExam.Result = req.Result
				objectiveExam.Time = req.Time
				if _, err := l.svcCtx.ObjectiveExamModel.Insert(l.ctx, objectiveExam); err != nil {
					return err
				}
				// 更新用户身份
				if err := l.svcCtx.UserModel.UpdateWithRole(l.ctx, req.UserId, 3); err != nil {
					return err
				}
				// 创建role_change表数据
				roleChange := new(models.RoleChange)
				roleChange.UserId = req.UserId
				roleChange.OperatorId = userId
				roleChange.NewRole = 3
				roleChange.OldRole = 2
				if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return response.Error(500, err.Error())
			}
		}
	}
	return response.Success("提交成功！")
}
