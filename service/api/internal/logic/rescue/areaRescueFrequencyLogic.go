package rescue

import (
	"context"
	"strings"

	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaRescueFrequencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaRescueFrequencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaRescueFrequencyLogic {
	return &AreaRescueFrequencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaRescueFrequencyLogic) AreaRescueFrequency() (resp *types.AreaRescueFrequencyResp, err error) {
	areas, err := l.svcCtx.RescueInfoModel.FindRescueFrequencyByArea(l.ctx)
	if err != nil {
		return nil, err
	}

	resp = &types.AreaRescueFrequencyResp{
		// List: make([]types.Area, len(areas)),
		List: []types.Area{
			{Name: "北京市", RescueFrequency: 0},
			{Name: "天津市", RescueFrequency: 0},
			{Name: "上海市", RescueFrequency: 0},
			{Name: "重庆市", RescueFrequency: 0},
			{Name: "河北省", RescueFrequency: 0},
			{Name: "河南省", RescueFrequency: 0},
			{Name: "云南省", RescueFrequency: 0},
			{Name: "辽宁省", RescueFrequency: 0},
			{Name: "黑龙江省", RescueFrequency: 0},
			{Name: "湖南省", RescueFrequency: 0},
			{Name: "安徽省", RescueFrequency: 0},
			{Name: "山东省", RescueFrequency: 0},
			{Name: "新疆维吾尔自治区", RescueFrequency: 0},
			{Name: "江苏省", RescueFrequency: 0},
			{Name: "浙江省", RescueFrequency: 0},
			{Name: "江西省", RescueFrequency: 0},
			{Name: "湖北省", RescueFrequency: 0},
			{Name: "广西壮族自治区", RescueFrequency: 0},
			{Name: "甘肃省", RescueFrequency: 0},
			{Name: "山西省", RescueFrequency: 0},
			{Name: "内蒙古自治区", RescueFrequency: 0},
			{Name: "陕西省", RescueFrequency: 0},
			{Name: "吉林省", RescueFrequency: 0},
			{Name: "福建省", RescueFrequency: 0},
			{Name: "贵州省", RescueFrequency: 0},
			{Name: "广东省", RescueFrequency: 0},
			{Name: "青海省", RescueFrequency: 0},
			{Name: "西藏自治区", RescueFrequency: 0},
			{Name: "四川省", RescueFrequency: 0},
			{Name: "宁夏回族自治区", RescueFrequency: 0},
			{Name: "海南省", RescueFrequency: 0},
			{Name: "台湾省", RescueFrequency: 0},
			{Name: "香港特别行政区", RescueFrequency: 0},
			{Name: "澳门特别行政区", RescueFrequency: 0},
		},
	}

	/* for i, area := range areas {
		resp.List[i] = types.Area{
			Name:            area.Name,
			RescueFrequency: area.RescueFrequency,
		}
	} */

	for _, area := range areas {
		if strings.Contains(area.Name, "北京") {
			resp.List[0].RescueFrequency = area.RescueFrequency
		} else if strings.Contains(area.Name, "天津") {
			resp.List[1].RescueFrequency = area.RescueFrequency
		} else if strings.Contains(area.Name, "上海") {
			resp.List[2].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"重庆") {
			resp.List[3].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"河北") {
			resp.List[4].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"河南") {
			resp.List[5].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"云南") {
			resp.List[6].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"辽宁") {
			resp.List[7].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"黑龙江") {
			resp.List[8].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"湖南") {
			resp.List[9].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"安徽") {
			resp.List[10].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"山东") {
			resp.List[11].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"新疆") {
			resp.List[12].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"江苏") {
			resp.List[13].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"浙江") {
			resp.List[14].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"江西") {
			resp.List[15].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"湖北") {
			resp.List[16].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"广西") {
			resp.List[17].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"甘肃") {
			resp.List[18].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"山西") {
			resp.List[19].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"内蒙古") {
			resp.List[20].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"陕西") {
			resp.List[21].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"吉林") {
			resp.List[22].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"福建") {
			resp.List[23].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"贵州") {
			resp.List[24].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"广东") {
			resp.List[25].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"青海") {
			resp.List[26].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"西藏") {
			resp.List[27].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"四川") {
			resp.List[28].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"宁夏") {
			resp.List[29].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"海南") {
			resp.List[30].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"台湾") {
			resp.List[31].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"香港") {
			resp.List[32].RescueFrequency = area.RescueFrequency
		}else if strings.Contains(area.Name,"澳门") {
			resp.List[33].RescueFrequency = area.RescueFrequency
		}
	}
	return resp, nil
}
