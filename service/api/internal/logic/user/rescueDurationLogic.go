package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RescueDurationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRescueDurationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RescueDurationLogic {
	return &RescueDurationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RescueDurationLogic) RescueDuration() (resp *types.RescueDurationResp, err error) {
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
	// 判断权限（只有见习队员以上的人才能有救援时长）
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	//
	rescueProcesses, err := l.svcCtx.RescueProcessModel.FindAllByRescueTeacherId(l.ctx, l.svcCtx.RescueProcessModel.RowBuilder(), userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	durations := []string{}
	if len(rescueProcesses) > 0 {
		for _, rescueProcessModel := range rescueProcesses {
			durations = append(durations, rescueProcessModel.Duration)
		}
	}

	return &types.RescueDurationResp{
		RescueDuration: sumDurations(durations),
	}, nil
}

func sumDurations(durations []string) string {
	totalDuration := time.Duration(0)

	for _, durationStr := range durations {
		parts := strings.Split(durationStr, "h") // 按字符'h'拆分字符串，例如"2h15m"，拆分后的结果为 []string{"2","15m"}
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(strings.TrimSuffix(parts[1], "m"))
		duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
		totalDuration += duration
	}
	hours := int(totalDuration.Hours())
	minutes := int(totalDuration.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}
