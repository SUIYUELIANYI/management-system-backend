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

type GetRescueProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRescueProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRescueProcessLogic {
	return &GetRescueProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRescueProcessLogic) GetRescueProcess(req *types.GetRescueProcessReq) (resp *types.GetRescueProcessResp, err error) {
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
	// 查看救援信息是否存在
	_, err = l.svcCtx.RescueInfoModel.FindOne(l.ctx, req.RescueInfoId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "救援信息不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	// 根据救援信息id查询对应的救援过程评价
	rescueProcesses, err := l.svcCtx.RescueProcessModel.FindAllByRescueInfoIdWithPage(l.ctx, l.svcCtx.RescueProcessModel.RowBuilder(), req.Page, req.PageSize, req.RescueInfoId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.RescueProcess
	if len(rescueProcesses) > 0 {
		for _, rescueProcessModel := range rescueProcesses {
			var typeRescueProcess types.RescueProcess

			_ = copier.Copy(&typeRescueProcess, rescueProcessModel)
			typeRescueProcess.CreateTime = rescueProcessModel.CreateTime.Format("2006-01-02 15:04:05")
			typeRescueProcess.UpdateTime = rescueProcessModel.UpdateTime.Format("2006-01-02 15:04:05")
			list = append(list, typeRescueProcess)
		}
	}
	return &types.GetRescueProcessResp{
		List: list,
	}, nil
}
