package rescue

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RescueProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRescueProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RescueProcessLogic {
	return &RescueProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RescueProcessLogic) RescueProcess(req *types.RescueProcessReq) error {
	// 救援信息评价（志愿者每次进行救援后都要对救援信息进行评价，并记录时间）
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	// 查看救援信息是否存在
	_, err = l.svcCtx.RescueInfoModel.FindOne(l.ctx, req.RescueInfoId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "救援信息不存在")
	default:
		return response.Error(500, err.Error())
	}
	// 插入救援过程数据
	rescueProcess := new(models.RescueProcess)
	rescueProcess.RescueTeacherId = userId
	rescueProcess.RescueInfoId = req.RescueInfoId
	rescueProcess.StartTime = req.StartTime // 格式：0000-00-00 00:00
	rescueProcess.EndTime = req.EndTime
	rescueProcess.Evaluation = req.Evaluation
	rescueProcess.Duration = getDuration(req.StartTime, req.EndTime)
	if _, err := l.svcCtx.RescueProcessModel.Insert(l.ctx, rescueProcess); err != nil {
		return response.Error(500, err.Error())
	}

	if userInfo.Role == 3 { // 如果是见习队员，提交后需要计算救援时长判断是否转正
		rescueProcesses, err := l.svcCtx.RescueProcessModel.FindAllByRescueTeacherId(l.ctx, l.svcCtx.RescueProcessModel.RowBuilder(), userId)
		if err != nil {
			return response.Error(500, err.Error())
		}
		durations := []string{}
		if len(rescueProcesses) > 0 {
			for _, rescueProcessModel := range rescueProcesses {
				durations = append(durations, rescueProcessModel.Duration)
			}
		}
		totalDuration := time.Duration(0)
		for _, durationStr := range durations {
			parts := strings.Split(durationStr, "h") // 按字符'h'拆分字符串，例如"2h15m"，拆分后的结果为 []string{"2","15m"}
			hours, _ := strconv.Atoi(parts[0])
			minutes, _ := strconv.Atoi(strings.TrimSuffix(parts[1], "m"))
			duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
			totalDuration += duration
		}
		hours := int(totalDuration.Hours()) // 救援时长
		// 缓存中查询救援时长阈值
		var rescueDurationThreshold int
		key := "rescue_duration_threshold"
		value, err := l.svcCtx.RedisClient.Get(key)
		if err != nil && err != red.Nil {
			return response.Error(500, "redis error:"+err.Error())
		} else if err != red.Nil {
			rescueDurationThreshold, err = strconv.Atoi(value)
			if err != nil {
				return response.Error(400, err.Error())
			}
		}
		// 如果设置了救援时长阈值就按救援阈值比较，否则按30小时算
		if (err == red.Nil && hours > 30) || (err != red.Nil && hours > rescueDurationThreshold) { // 如果大于规定的救援时长，role变为40-普通队员
			roleChange := new(models.RoleChange)
			roleChange.UserId = userId
			roleChange.OperatorId = userId // 由于是系统自动判断，操作人是自己
			roleChange.NewRole = 40
			roleChange.OldRole = 3
			userInfo.Role = 40
			if err := l.svcCtx.RoleChangeModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
				if _, err := l.svcCtx.RoleChangeModel.Insert(l.ctx, roleChange); err != nil {
					return err
				}
				if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return response.Error(500, err.Error())
			}
		}
	}
	return nil
}

func getDuration(start, end string) string {
	layout := "2006-01-02 15:04"
	startTime, _ := time.Parse(layout, start)
	endTime, _ := time.Parse(layout, end)

	duration := endTime.Sub(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	durationStr := strconv.Itoa(hours) + "h" + strconv.Itoa(minutes) + "m"
	return durationStr
}
