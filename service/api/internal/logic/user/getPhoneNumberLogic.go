package user

import (
	"context"
	"fmt"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPhoneNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPhoneNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPhoneNumberLogic {
	return &GetPhoneNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPhoneNumberLogic) GetPhoneNumber(req *types.GetPhoneNumberReq) (resp *types.GetPhoneNumberResp, err error) {
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	})
	authResult, err := miniprogram.GetAuth().GetPhoneNumber(req.Code)
	if err != nil || authResult.ErrCode != 0 {
		return nil, response.Error(100, fmt.Sprintf("发起授权请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult))
	}

	var phoneInfo types.PhoneInfo
	phoneInfo.PhoneNumber = authResult.PhoneInfo.PhoneNumber
	phoneInfo.PurePhoneNumber = authResult.PhoneInfo.PurePhoneNumber
	phoneInfo.CountryCode = authResult.PhoneInfo.CountryCode
	phoneInfo.WaterMark.Timestamp = authResult.PhoneInfo.WaterMark.Timestamp
	phoneInfo.WaterMark.AppID = authResult.PhoneInfo.WaterMark.AppID

	return &types.GetPhoneNumberResp{
		PhoneInfo: phoneInfo,
	}, nil
}
