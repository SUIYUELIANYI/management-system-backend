package rescue

import (
	"context"

	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type YearRescueFrequencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewYearRescueFrequencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *YearRescueFrequencyLogic {
	return &YearRescueFrequencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *YearRescueFrequencyLogic) YearRescueFrequency() (resp *types.YearRescueFrequencyResp, err error) {
	areas, err := l.svcCtx.RescueInfoModel.FindRescueFrequencyByYear(l.ctx)
	if err != nil {
		return nil, err
	}

	resp = &types.YearRescueFrequencyResp{
		List: make([]types.Year, len(areas)),
	}

	for i, area := range areas {
		resp.List[i] = types.Year{
			Name:            area.Name,
			RescueFrequency: area.RescueFrequency,
		}
	}

	return resp, nil
}
