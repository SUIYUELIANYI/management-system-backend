package rescue

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/models"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ReleaseResuceInfoByExcelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewReleaseResuceInfoByExcelLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *ReleaseResuceInfoByExcelLogic {
	return &ReleaseResuceInfoByExcelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *ReleaseResuceInfoByExcelLogic) ReleaseResuceInfoByExcel() error {
	file, _, err := l.r.FormFile("file")
	if err != nil {
		response.Error(450, err.Error())
	}

	// 读取excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		response.Error(450, "open excel error:"+err.Error())
	}

	// 解析excel数据
	rescueInfos, lxRrr := readExcel(xlsx)
	if lxRrr != nil {
		response.Error(450, "解析excel数据失败:"+err.Error())
	}

	for _, rescueInfo := range rescueInfos {
		// 先查重
		_, err := l.svcCtx.RescueInfoModel.FindOneByReleaseTimeWeiboAddress(l.ctx, rescueInfo.ReleaseTime, rescueInfo.WeiboAddress)
		if err != nil && err != models.ErrNotFound {
			return response.Error(500, err.Error())
		} else if err == models.ErrNotFound {
			// 救援信息导入时，先根据微博地址查询是否存在该救援对象，如果不存在，先创建救援对象RescueTarget
			rescueTarget, err := l.svcCtx.RescueTargetModel.FindOneByWeiboAddress(l.ctx, rescueInfo.WeiboAddress)
			if err != nil && err != models.ErrNotFound {
				return response.Error(500, err.Error())
			} else if err == models.ErrNotFound {
				// 如果以前没有存在该救援对象，有两个事务，一个是创建救援对象，一个是创建救援信息
				if err := l.svcCtx.RescueTargetModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
					t := new(models.RescueTarget)
					t.WeiboAddress = rescueInfo.WeiboAddress
					t.Status = 0
					t.StartTime = time.Now()

					result, err := l.svcCtx.RescueTargetModel.Insert(l.ctx, t)
					if err != nil {
						return err
					}

					rescueInfo.RescueTargetId, err = result.LastInsertId()
					if err != nil {
						return err
					}
					if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, &rescueInfo); err != nil {
						return err
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
					t.WeiboAddress = rescueInfo.WeiboAddress
					t.Status = 0
					t.StartTime = time.Now()
					result, err := l.svcCtx.RescueTargetModel.Insert(l.ctx, t)
					if err != nil {
						return err
					}
					// 创建救援信息
					rescueInfo.RescueTargetId, err = result.LastInsertId()
					if err != nil {
						return err
					}
					if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, &rescueInfo); err != nil {
						return err
					}
					return nil
				}); err != nil {
					return response.Error(500, err.Error())
				}
			} else { // 已经存在该救援对象，且当前救援对象正在救援或者待救援中，直接插入救援信息即可
				rescueInfo.RescueTargetId = rescueTarget.Id
				if _, err := l.svcCtx.RescueInfoModel.Insert(l.ctx, &rescueInfo); err != nil {
					return response.Error(500, err.Error())
				}
			}
		}
	}
	return nil
}

// readExcel：读取excel转成切片
func readExcel(xlsx *excelize.File) ([]models.RescueInfo, error) {
	// 根据名字获取cells的内容，返回的是一个[][]string
	// GetRows按给定的工作表名称返回工作表中的所有行，返回为二维数组
	var lxProduct []models.RescueInfo
	for _, sheetName := range xlsx.GetSheetList() {
		rows, err := xlsx.GetRows(sheetName)
		if err != nil {
			response.Error(450, "excel get rows error:"+err.Error())
		}

		for i, row := range rows {
			// 去掉第一行是excel表头部分
			if i == 0 {
				continue
			}
			var data models.RescueInfo
			for k, v := range row {
				// 第一列是发布时间
				if k == 0 {
					data.ReleaseTime = v
				}
				// 第二列是微博账号
				if k == 1 {
					data.WeiboAccount = v
				}
				// 第三列是微博地址
				if k == 2 {
					data.WeiboAddress = v
				}
				// 第四列是昵称
				if k == 3 {
					data.Nickname = v
				}
				// 第五列是危险级别
				if k == 4 {
					data.RiskLevel = v
				}
				// 第六列是所在城市
				if k == 5 {
					data.Area = v
				}
				// 第七列是性别（0-男，1-女）
				if k == 6 {
					if v == "男" || v == "0" {
						data.Sex = 0
					} else if v == "女" || v == "1" {
						data.Sex = 1
					}
				}
				// 第八列是生日
				if k == 7 {
					data.Birthday = v
				}
				// 第九列是简介
				if k == 8 {
					data.BriefIntroduction = v
				}
				// 第十列是信息原文
				if k == 9 {
					data.Text = sql.NullString{String: v, Valid: true}
				}
			}
			// 添加到切片中
			lxProduct = append(lxProduct, data)
		}
	}

	// 声明一个数组
	return lxProduct, nil
}
