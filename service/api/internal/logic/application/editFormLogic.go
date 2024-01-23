package application

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

type EditFormLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditFormLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditFormLogic {
	return &EditFormLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditFormLogic) EditForm(req *types.EditReq) error {
	// 申请表状态 0-待审批 1-区域负责人通过 2-组织管理委员会通过 3-未通过
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 先根据userId判断身份，必须是0-非在册队员 1-申请队员才能调用该接口
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return response.Error(500, err.Error())
	}
	if userInfo.Role != 1 && userInfo.Role != 0 {
		return response.Error(401, "您不需要提交/修改申请表！")
	}
	// 再根据userId找到申请表
	form, err := l.svcCtx.ApplicationFormModel.FindOneByUserId(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		// 如果申请表不存在，说明是非在册成员
		form = new(models.ApplicationForm)
		form.UserId = userId
		form.Mobile = req.Mobile
		if req.Mobile == "" {
			return response.Error(400, "电话不能为空！")
		} else {
			for _, char := range req.Mobile {
				if string(char) == " " {
					return response.Error(400, "电话中不能含有空格！")
				}
			}
		}
		form.Username = req.Username
		if req.Username == "" {
			return response.Error(400, "昵称不能为空！")
		}
		form.Sex = req.Sex
		form.Address = req.Address
		if req.Address == "" {
			return response.Error(400, "地址不能为空！")
		}
		form.Birthday = req.Birthday
		if req.Birthday == "" {
			return response.Error(400, "生日不能为空！")
		}
		form.Email = req.Email
		if req.Email == "" {
			return response.Error(400, "邮箱不能为空！")
		}
		form.SubmissionTime = 1
		// 修改用户信息的身份
		userInfo.Role = 1
		if err := l.svcCtx.ApplicationFormModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			if _, err := l.svcCtx.ApplicationFormModel.Insert(l.ctx, form); err != nil {
				return err
			}
			if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return response.Error(500, err.Error())
		}
		return response.Success("提交申请表成功！")
	default:
		return response.Error(500, err.Error())
	}

	form.UserId = userId
	form.Mobile = req.Mobile
	form.Username = req.Username
	form.Sex = req.Sex
	form.Address = req.Address
	form.Birthday = req.Birthday
	form.Email = req.Email

	if form.Status == 0 || form.Status == 1 {
		return response.Success("您的申请表正在审批中！")
	} else if form.Status == 2 {
		return response.Success("您的申请表已通过！")
	} else if form.Status == 3 && form.SubmissionTime == 3 {
		return response.Error(401, "很遗憾，您的3次申请均未通过，无法再次修改！")
	} else {
		form.SubmissionTime = form.SubmissionTime + 1
		form.Status = 0 // 重置申请表状态
		err = l.svcCtx.ApplicationFormModel.Update(l.ctx, form)
		if err != nil {
			return response.Error(500, err.Error())
		}
	}

	return response.Success(fmt.Sprintf("成功修改申请表！您还有%d次机会！", 3-form.SubmissionTime))
}
