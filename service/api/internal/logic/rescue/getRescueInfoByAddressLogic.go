package rescue

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRescueInfoByAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRescueInfoByAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRescueInfoByAddressLogic {
	return &GetRescueInfoByAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRescueInfoByAddressLogic) GetRescueInfoByAddress(req *types.GetRescueInfoByAddressReq) (resp *types.RescueInfosResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	// 判断输入地址是否有效
	if len(req.Address) < 2 {
		return nil, response.Error(400, "地址无效")
	}
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 根据地址查询救援信息
	address := string([]rune(req.Address)[:2])
	rescueInfos, err := l.svcCtx.RescueInfoModel.FindAllByAddressWithPage(l.ctx, l.svcCtx.RescueInfoModel.RowBuilder(), req.Page, req.PageSize, address)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.RescueInfo
	if len(rescueInfos) > 0 {
		for _, rescueInfoModel := range rescueInfos {
			var typeRescueInfo types.RescueInfo

			_ = copier.Copy(&typeRescueInfo, rescueInfoModel)
			typeRescueInfo.CreateTime = rescueInfoModel.CreateTime.Format("2006-01-02 15:04:05")
			typeRescueInfo.UpdateTime = rescueInfoModel.UpdateTime.Format("2006-01-02 15:04:05")
			list = append(list, typeRescueInfo)
		}
	}
	return &types.RescueInfosResp{
		List: list,
	}, nil
}
