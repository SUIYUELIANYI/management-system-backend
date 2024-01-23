package rescue

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ReleaseRescueInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReleaseRescueInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseRescueInfoLogic {
	return &ReleaseRescueInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReleaseRescueInfoLogic) ReleaseRescueInfo(req *types.ReleaseRescueInfoReq) error {
	// Wechat-Mini
	/* miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.AppSecret,
		Cache:     cache.NewMemory(),
	}) */

	// 获取access_token
	/* redisKey := "access_token"
	accessToken, err := l.svcCtx.RedisClient.Get(redisKey)
	if err != nil && err != red.Nil {
		return response.Error(500, "redis error:"+err.Error())
	}
	if accessToken == "" {
		accessToken, err := miniprogram.GetAuth().GetAccessToken()
		if err != nil {
			return response.Error(100, fmt.Sprintf("获取接口调用凭证失败 err : %v , authResult: %v", err, accessToken))
		}
		// 设置缓存,access_token过期时间为7200s（2小时）
		err = l.svcCtx.RedisClient.Setex(redisKey, accessToken, 7200)
		if err != nil {
			return response.Error(500, "redis error:"+err.Error())
		}
	}
	fmt.Println(accessToken) */
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	// 因为存在客户端的rescueToken是不会被销毁的，我们需要再去缓存中判断该token是否有效
	key := "management_system_rescue_token_" + strconv.Itoa(int(userId))
	_, err := l.svcCtx.RedisClient.Get(key)
	switch err {
	case nil:
	case red.Nil:
		return response.Error(100, "认证已失效，请重新认证！")
	default:
		return response.Error(500, err.Error())
	}

	// 设置本次救援信息需要发送报警的电话号码（如果该信息的救援对象已被认领，则只给认领的救援老师发送，否则给全部人发送）
	/* var phoneNumberList []string
	requestData := map[string]interface{}{
		"env": "management-system-3ery9yf2d68d54",
		// "url_link":            "",
		"template_id":         "844110",
		"template_param_list": []string{"树洞救援"},
		"phone_number_list":   phoneNumberList,
		"resource_appid":      "wx8b456b9fba47330f",
		"sms_type":            "Notification",
		"path":                "/index.html",
	} */

	_, err = l.svcCtx.RescueInfoModel.FindOneByReleaseTimeWeiboAddress(l.ctx, req.ReleaseTime, req.WeiboAddress)
	if err != nil && err != models.ErrNotFound {
		return response.Error(500, err.Error())
	} else if err == nil {
		return response.Error(400, "数据重复！")
	}

	// 救援信息导入时，先根据微博地址查询是否存在该救援对象，如果不存在，先创建救援对象RescueTarget
	rescueTarget, err := l.svcCtx.RescueTargetModel.FindOneByWeiboAddress(l.ctx, req.WeiboAddress)
	if err != nil && err != models.ErrNotFound {
		return response.Error(500, err.Error())
	} else if err == models.ErrNotFound {
		// 如果以前没有存在该救援对象，有两个事务，一个是创建救援对象，一个是创建救援信息
		if err := l.svcCtx.RescueTargetModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			t := new(models.RescueTarget)
			t.WeiboAddress = req.WeiboAddress
			t.Status = 0
			t.StartTime = time.Now()
			result, err := l.svcCtx.RescueTargetModel.Insert(l.ctx, t)
			if err != nil {
				return err
			}
			rescueInfo := new(models.RescueInfo)
			rescueInfo.RescueTargetId, err = result.LastInsertId()
			if err != nil {
				return err
			}
			rescueInfo.ReleaseTime = req.ReleaseTime
			rescueInfo.WeiboAccount = req.WeiboAccount
			rescueInfo.WeiboAddress = req.WeiboAddress
			rescueInfo.Nickname = req.Nickname
			rescueInfo.RiskLevel = req.RiskLevel
			rescueInfo.Area = req.Area
			rescueInfo.Sex = req.Sex
			rescueInfo.Birthday = req.Birthday
			rescueInfo.BriefIntroduction = req.BriefIntroduction
			rescueInfo.Text = sql.NullString{String: req.Text, Valid: true}
			if req.Flag == 1 { // 发布订阅消息
				if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
					return err
				}
				/* // 返回所有用户电话列表
				phoneNumberList, err = l.svcCtx.UserModel.FindAllMoblie(l.ctx)
				if err != nil {
					return err
				}
				// 测试打印
				fmt.Printf("[%s]\n", strings.Join(phoneNumberList, ", "))

				jsonData, err := json.Marshal(requestData)
				if err != nil {
					fmt.Println("Error encoding JSON data: ", err)
					return err
				}

				// 构建请求 URL
				url := fmt.Sprintf("https://api.weixin.qq.com/tcb/sendsms?access_token=%s", accessToken)

				// 发送 POST 请求
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error sending request: ", err)
					return err
				}

				defer resp.Body.Close()

				// 读取响应内容
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response: ", err)
					return err
				}

				// 打印响应内容
				fmt.Println(string(body)) */
			} else { // 不发布订阅信息
				if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return response.Error(500, err.Error())
		}
	} else if rescueTarget.Status == 2 { // 如果该救援对象已经存在且已完成救助，重新设置新的救援对象，微博地址相同但id不同
		// 事务包含两个操作，一个是更新救援对象信息，另一个是插入救援信息
		if err := l.svcCtx.RescueTargetModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			// 插入新的救援对象信息
			t := new(models.RescueTarget)
			t.WeiboAddress = req.WeiboAddress
			t.Status = 0
			t.StartTime = time.Now()
			result, err := l.svcCtx.RescueTargetModel.Insert(l.ctx, t)
			if err != nil {
				return err
			}
			// 创建救援信息
			rescueInfo := new(models.RescueInfo)
			rescueInfo.RescueTargetId, err = result.LastInsertId()
			if err != nil {
				return err
			}
			rescueInfo.ReleaseTime = req.ReleaseTime
			rescueInfo.WeiboAccount = req.WeiboAccount
			rescueInfo.WeiboAddress = req.WeiboAddress
			rescueInfo.Nickname = req.Nickname
			rescueInfo.RiskLevel = req.RiskLevel
			rescueInfo.Area = req.Area
			rescueInfo.Sex = req.Sex
			rescueInfo.Birthday = req.Birthday
			rescueInfo.BriefIntroduction = req.BriefIntroduction
			rescueInfo.Text = sql.NullString{String: req.Text, Valid: true}
			if req.Flag == 1 { // Flag=1，本次发布救援信息需要同时发布订阅消息（该部分未实现）
				if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
					return err
				}

				/* // 返回所有用户电话列表
				phoneNumberList, err = l.svcCtx.UserModel.FindAllMoblie(l.ctx)
				if err != nil {
					return err
				}
				// 测试打印
				fmt.Printf("[%s]\n", strings.Join(phoneNumberList, ", "))

				jsonData, err := json.Marshal(requestData)
				if err != nil {
					fmt.Println("Error encoding JSON data: ", err)
					return err
				}

				// 构建请求 URL
				url := fmt.Sprintf("https://api.weixin.qq.com/tcb/sendsms?access_token=%s", accessToken)

				// 发送 POST 请求
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Println("Error sending request: ", err)
					return err
				}

				defer resp.Body.Close()

				// 读取响应内容
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response: ", err)
					return err
				}

				// 打印响应内容
				fmt.Println(string(body)) */

			} else { // Flag=0，本次发布救援信息不需要发布订阅消息
				if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return response.Error(500, err.Error())
		}
	} else { // 已经存在该救援对象，且当前救援对象正在救援或者待救援中，直接插入救援信息即可
		rescueInfo := new(models.RescueInfo)
		rescueInfo.RescueTargetId = rescueTarget.Id
		rescueInfo.ReleaseTime = req.ReleaseTime
		rescueInfo.WeiboAccount = req.WeiboAccount
		rescueInfo.WeiboAddress = req.WeiboAddress
		rescueInfo.Nickname = req.Nickname
		rescueInfo.RiskLevel = req.RiskLevel
		rescueInfo.Area = req.Area
		rescueInfo.Sex = req.Sex
		rescueInfo.Birthday = req.Birthday
		rescueInfo.BriefIntroduction = req.BriefIntroduction
		rescueInfo.Text = sql.NullString{String: req.Text, Valid: true}
		if req.Flag == 1 { // 发布订阅消息
			if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
				return response.Error(500, err.Error())
			}

			/* 	if rescueTarget.Status == 0 { // 如果该救援信息的救援对象仍然是待救援状态
				// 返回所有用户电话列表
				phoneNumberList, err = l.svcCtx.UserModel.FindAllMoblie(l.ctx)
				if err != nil {
					return err
				}

				// 测试打印
				fmt.Printf("[%s]\n", strings.Join(phoneNumberList, ", "))
			} else if rescueTarget.Status == 1 {
				phoneNumberList, err = l.svcCtx.RescueTargetModel.FindMobileById(l.ctx, rescueTarget.Id)
				if err != nil {
					return err
				}

				// 测试打印
				fmt.Printf("[%s]\n", strings.Join(phoneNumberList, ", "))
			}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				fmt.Println("Error encoding JSON data: ", err)
				return err
			}

			// 构建请求 URL
			url := fmt.Sprintf("https://api.weixin.qq.com/tcb/sendsms?access_token=%s", accessToken)

			// 发送 POST 请求
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error sending request: ", err)
				return err
			}

			defer resp.Body.Close()

			// 读取响应内容
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response: ", err)
				return err
			}

			// 打印响应内容
			fmt.Println(string(body)) */
		} else {
			if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, rescueInfo); err != nil {
				return response.Error(500, err.Error())
			}
		}
	}
	return nil
}
