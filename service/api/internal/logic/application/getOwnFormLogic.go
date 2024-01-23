package application

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOwnFormLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOwnFormLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOwnFormLogic {
	return &GetOwnFormLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOwnFormLogic) GetOwnForm() (resp *types.FormResp, err error) {
	// 用户个人查看自己的申请表
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	form, err := l.svcCtx.ApplicationFormModel.FindOneByUserId(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "申请表不存在")
	default:
		return nil, response.Error(500, err.Error())
	}

	var typeApplicationForm types.ApplicationForm
	typeApplicationForm.Id = form.Id
	typeApplicationForm.CreateTime = form.CreateTime.Format("2006-01-02 15:04:05")
	typeApplicationForm.UpdateTime = form.UpdateTime.Format("2006-01-02 15:04:05")
	typeApplicationForm.UserId = form.UserId
	typeApplicationForm.Mobile = form.Mobile
	typeApplicationForm.Username = form.Username
	typeApplicationForm.Sex = form.Sex
	typeApplicationForm.Address = form.Address
	typeApplicationForm.Birthday = form.Birthday
	typeApplicationForm.Email = form.Email
	typeApplicationForm.Status = form.Status
	typeApplicationForm.RegionalHeadId = form.RegionalHeadId
	typeApplicationForm.OrganizingCommitteeId = form.OrganizingCommitteeId
	typeApplicationForm.SubmissionTime = form.SubmissionTime

	return &types.FormResp{
		Form: typeApplicationForm,
	}, nil
}
