package application

import (
	"context"
	"encoding/json"
	"fmt"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ProcessWithSubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessWithSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessWithSubscribeLogic {
	return &ProcessWithSubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessWithSubscribeLogic) ProcessWithSubscribe(req *types.ProcessWithSubscribeReq) error {
	// Wechat-Mini
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	})
	// 仅限区域负责人（42）和组织管理委员会（43，44）可以调用，且两种身份审批结果对应不同状态
	// 申请表状态 0-待审批 1-区域负责人通过 2-组织管理委员会通过 3-未通过
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

	// 申请人信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, form.UserId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "申请人不存在")
	default:
		return response.Error(500, err.Error())
	}

	// 申请人openId
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByUserIdAuthType(l.ctx, userInfo.Id, models.UserAuthTypeSmallWX)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "申请人授权信息不存在")
	default:
		return response.Error(500, err.Error())
	}

	if operatorInfo.Role == 42 {
		form.RegionalHeadId = operatorInfo.Id
		if req.Result == 0 { // 未通过
			// 未通过，发送订阅消息
			err := miniprogram.GetSubscribe().Send(&subscribe.Message{
				ToUser:     userAuth.AuthKey,
				TemplateID: "_xbtX5gtyvi8uhpHsAmt-XhcdfTpXd7HLWh-ahuUptU",
				Data: map[string]*subscribe.DataItem{
					"thing1":  {Value: form.Username, Color: ""}, // 申请人（字符串类型）
					"thing2":  {Value: "提交申请表", Color: ""},       // 申请内容（字符串类型）
					"phrase3": {Value: "未通过审批", Color: ""},  // 当前进度（字符串类型）有长度限制
				},
			})
			if err != nil {
				return response.Error(100, fmt.Sprintf("发起订阅消息请求失败 err : %v ", err))
			}

			form.Status = 3
			if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
				return response.Error(500, err.Error())
			}
		} else if req.Result == 1 { // 通过
			form.Status = 1
			if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
				return response.Error(500, err.Error())
			}
		}
	} else if operatorInfo.Role == 43 || operatorInfo.Role == 44 {
		form.OrganizingCommitteeId = operatorInfo.Id
		if req.Result == 0 { // 未通过
			// 未通过，发送订阅消息
			err := miniprogram.GetSubscribe().Send(&subscribe.Message{
				ToUser:     userAuth.AuthKey,
				TemplateID: "_xbtX5gtyvi8uhpHsAmt-XhcdfTpXd7HLWh-ahuUptU",
				Data: map[string]*subscribe.DataItem{
					"thing1":  {Value: form.Username, Color: ""}, // 申请人
					"thing2":  {Value: "提交申请表", Color: ""},       // 申请内容
					"phrase3": {Value: "未通过审批", Color: ""},  // 当前进度
				},
			})
			if err != nil {
				return response.Error(100, fmt.Sprintf("发起订阅消息请求失败 err : %v ", err))
			}
			form.Status = 3
			if err := l.svcCtx.ApplicationFormModel.Update(l.ctx, form); err != nil {
				return response.Error(500, err.Error())
			}
		} else if req.Result == 1 { // 只有组织管理委员会通过时才会更改身份，需要同时修改user，role_change，application_form表信息，其他情况只需要更新申请表信息即可
			// 组委会通过，发送订阅消息
			err := miniprogram.GetSubscribe().Send(&subscribe.Message{
				ToUser:     userAuth.AuthKey,
				TemplateID: "_xbtX5gtyvi8uhpHsAmt-XhcdfTpXd7HLWh-ahuUptU",
				Data: map[string]*subscribe.DataItem{
					"thing1":  {Value: form.Username, Color: ""}, // 申请人
					"thing2":  {Value: "提交申请表", Color: ""},       // 申请内容
					"phrase3": {Value: "已通过审批", Color: ""},  // 当前进度
				},
			})
			if err != nil {
				return response.Error(100, fmt.Sprintf("发起订阅消息请求失败 err : %v ", err))
			}

			form.Status = 2
			// 申请人信息
			applicantInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, form.UserId)
			if err != nil {
				response.Error(500, err.Error())
			}
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
	}
	return response.Success("审批成功")
}

// Message 订阅消息请求参数
/* type Message struct {
	ToUser           string               `json:"touser"`            // 必选，接收者（用户）的 openid
	TemplateID       string               `json:"template_id"`       // 必选，所需下发的订阅模板id
	Page             string               `json:"page"`              // 可选，点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Data             map[string]*DataItem `json:"data"`              // 必选, 模板内容
	MiniprogramState string               `json:"miniprogram_state"` // 可选，跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string               `json:"lang"`              // 入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
} */

/* data := &map[string]subscribe.DataItem{
	"thing1":  {Value: "", Color: ""}, // 申请人
	"thing2":  {Value: "", Color: ""}, // 申请内容
	"phrase3": {Value: "", Color: ""}, // 当前进度
}
err := miniprogram.GetSubscribe().Send(&subscribe.Message{
	ToUser:     "",
	TemplateID: "_xbtX5gtyvi8uhpHsAmt-XhcdfTpXd7HLWh-ahuUptU",
	Data: data,
}) */
