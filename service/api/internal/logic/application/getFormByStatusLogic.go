package application

import (
	"context"
	"fmt"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFormByStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFormByStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFormByStatusLogic {
	return &GetFormByStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFormByStatusLogic) GetFormByStatus(req *types.GetByStatusReq) (resp *types.GetByStatusResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))

	status := fmt.Sprintf("status = %d", req.Status)
	forms, err := l.svcCtx.ApplicationFormModel.FindAllByStatusWithPage(l.ctx, l.svcCtx.ApplicationFormModel.RowBuilder(), req.Page, req.PageSize, status)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	var list []types.ApplicationForm
	if len(forms) > 0 {
		for _, form := range forms {
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
			list = append(list, typeApplicationForm)
		}
	}
	return &types.GetByStatusResp{
		List: list,
	}, nil
}
