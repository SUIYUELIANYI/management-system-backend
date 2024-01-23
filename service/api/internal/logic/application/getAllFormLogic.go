package application

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllFormLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllFormLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllFormLogic {
	return &GetAllFormLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllFormLogic) GetAllForm(req *types.FormsReq) (resp *types.FormsResp, err error) {
	// 仅限区域负责人（42）和组织管理委员会（43，44）可以调用
	// 申请表状态 0-待审批 1-区域负责人通过 2-组织管理委员会通过 3-未通过
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	forms, err := l.svcCtx.ApplicationFormModel.FindAllWithPage(l.ctx, l.svcCtx.ApplicationFormModel.RowBuilder(), req.Page, req.PageSize)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.ApplicationForm
	if userInfo.Role == 42 { // 区域负责人负责 0-待审批的申请表
		if len(forms) > 0 {
			for _, form := range forms {
				var typeApplicationForm types.ApplicationForm
				if form.Status == 0 {
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
		}
	} else if userInfo.Role == 43 || userInfo.Role == 44 {
		if len(forms) > 0 {
			for _, form := range forms {
				var typeApplicationForm types.ApplicationForm
				if form.Status == 1 {
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
		}
	} else {
		return nil, response.Error(403, "权限不够")
	}

	return &types.FormsResp{
		List: list,
	}, nil
}
