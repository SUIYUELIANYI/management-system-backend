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

type GetClaimedRescueInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClaimedRescueInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClaimedRescueInfoLogic {
	return &GetClaimedRescueInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClaimedRescueInfoLogic) GetClaimedRescueInfo(req *types.RescueInfosReq) (resp *types.RescueInfosResp, err error) {
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
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}

	RescueInfos, err := l.svcCtx.RescueInfoModel.FindClaimedByUserIdWithPage(l.ctx, l.svcCtx.RescueInfoModel.RowBuilder(), req.Page, req.PageSize, userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.RescueInfo
	if len(RescueInfos) > 0 {
		for _, rescueInfoModel := range RescueInfos {
			var typeRescueInfo types.RescueInfo
			_, err := l.svcCtx.RescueProcessModel.FindOneByInfoIdTeacherId(l.ctx, rescueInfoModel.Id, userId) // 查询是否评价过该救援信息
			switch err {
			case nil:
			case models.ErrNotFound:
				_ = copier.Copy(&typeRescueInfo, rescueInfoModel)
				typeRescueInfo.CreateTime = rescueInfoModel.CreateTime.Format("2006-01-02 15:04:05")
				typeRescueInfo.UpdateTime = rescueInfoModel.UpdateTime.Format("2006-01-02 15:04:05")
				list = append(list, typeRescueInfo)
			default:
				return nil, response.Error(500, err.Error())
			}

		}
	}
	return &types.RescueInfosResp{
		List: list,
	}, nil
}
