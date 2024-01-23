package application

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

type ProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessLogic {
	return &ProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessLogic) Process(req *types.ProcessReq) error {
	// 仅限区域负责人（42）和组织管理委员会（43，44）可以调用，且两种身份审批结果对应不同状态
	// 申请表状态 0-待审批 1-区域负责人通过 2-组织管理委员会通过 3-未通过
	// 申请表可以修改3次，如果第3次修改即form.submission_time=3时仍未通过，将该用户role设为-1，并添加到待处理人员表，理由为：申请表三次未通过
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 操作人信息（分为区域负责人和组织委员会成员）
	operatorInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return response.Error(500, err.Error())
	}
	if operatorInfo.Role != 42 && operatorInfo.Role != 43 && operatorInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	// 申请表信息
	form, err := l.svcCtx.ApplicationFormModel.FindOne(l.ctx, req.ApplicationFormId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "申请表不存在")
	default:
		return response.Error(500, err.Error())
	}

	if form.Status == 3 {
		return response.Error(100, "申请表已审批")
	}

	if operatorInfo.Role == 42 && form.Status == 0 {
		form.RegionalHeadId = operatorInfo.Id
		if req.Result == 0 { // 未通过
			form.Status = 3
			if form.SubmissionTime == 3 { // 如果是第三次提交还未通过，则用户身份为-1
				// 事务实现更新申请表状态，更新用户身份，插入身份变动数据，插入待处理人员数据
				if err := l.svcCtx.ApplicationFormModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
					if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
						return err
					}
					// 申请人信息
					applicantInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, form.UserId)
					if err != nil {
						return err
					}
					applicantInfo.Role = -1
					if err := l.svcCtx.UserModel.Update(l.ctx, applicantInfo); err != nil {
						return err
					}
					// 创建role_change表数据
					roleChange := new(models.RoleChange)
					roleChange.UserId = applicantInfo.Id
					roleChange.OperatorId = userId
					roleChange.NewRole = -1
					roleChange.OldRole = 1
					if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
						return err
					}
					// 创建待处理人员表数据
					pendingpersonnel := new(models.PendingPersonnel)
					pendingpersonnel.UserId = applicantInfo.Id
					pendingpersonnel.Reason = models.ReasonForApplicationForm
					pendingpersonnel.OperateId = userId
					if _, err := l.svcCtx.PendingPersonnelModel.Insert(l.ctx, pendingpersonnel); err != nil {
						return err
					}
					return nil
				}); err != nil {
					return response.Error(500, err.Error())
				}
			} else {
				if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
					return response.Error(500, err.Error())
				}
			}
		} else if req.Result == 1 { // 通过
			form.Status = 1
			if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
				return response.Error(500, err.Error())
			}
		}
	} else if (operatorInfo.Role == 43 || operatorInfo.Role == 44) && form.Status == 1 {
		form.OrganizingCommitteeId = operatorInfo.Id
		if req.Result == 0 {
			form.Status = 3
			if form.SubmissionTime == 3 { // 如果是第三次提交还未通过，则用户身份为-1
				// 事务实现更新申请表状态，更新用户身份，插入身份变动数据，插入待处理人员数据
				if err := l.svcCtx.ApplicationFormModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
					if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
						return err
					}
					// 申请人信息
					applicantInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, form.UserId)
					if err != nil {
						return err
					}
					applicantInfo.Role = -1
					if err := l.svcCtx.UserModel.Update(l.ctx, applicantInfo); err != nil {
						return err
					}
					// 创建role_change表数据
					roleChange := new(models.RoleChange)
					roleChange.UserId = applicantInfo.Id
					roleChange.OperatorId = userId
					roleChange.NewRole = -1
					roleChange.OldRole = 1
					if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
						return err
					}
					// 创建待处理人员表数据
					pendingpersonnel := new(models.PendingPersonnel)
					pendingpersonnel.UserId = applicantInfo.Id
					pendingpersonnel.Reason = models.ReasonForApplicationForm
					pendingpersonnel.OperateId = userId
					if _, err := l.svcCtx.PendingPersonnelModel.Insert(l.ctx, pendingpersonnel); err != nil {
						return err
					}
					return nil
				}); err != nil {
					return response.Error(500, err.Error())
				}
			} else {
				if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
					return response.Error(500, err.Error())
				}
			}
		} else if req.Result == 1 { // 只有组织管理委员会通过时才会更改身份，需要同时修改user，role_change，application_form表信息，其他情况只需要更新申请表信息即可
			form.Status = 2
			// 申请人信息
			applicantInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, form.UserId)
			if err != nil {
				response.Error(500, err.Error())
			}

			applicantInfo.Mobile = form.Mobile
			applicantInfo.Username = form.Username
			applicantInfo.Sex = form.Sex
			applicantInfo.Email = form.Email
			applicantInfo.Address = form.Address
			applicantInfo.Birthday = form.Birthday
			applicantInfo.Role = 2

			// 创建role_change表数据
			roleChange := new(models.RoleChange)
			roleChange.UserId = applicantInfo.Id
			roleChange.OperatorId = userId
			roleChange.NewRole = 2
			roleChange.OldRole = 1

			// 事务实现三个数据库操作
			if err := l.svcCtx.ApplicationFormModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
					return err
				}
				if err := l.svcCtx.UserModel.Update(l.ctx, applicantInfo); err != nil {
					return err
				}
				if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return response.Error(500, err.Error())
			}
		}
	} else {
		return response.Error(100, "身份不匹配：区域负责人和组委会负责审批不同的申请表。")
	}
	return response.Success("审批成功")
}
