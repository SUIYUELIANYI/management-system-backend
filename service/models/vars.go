package models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

var UserAuthTypeSystem string = "system"  // 平台内部手机密码登录
var UserAuthTypeSmallWX string = "wxMini" // 微信小程序登录

var ReasonForApplicationForm string = "申请表三次未通过 " // 申请表三次未通过
var ReasonForSubjectiveExam string = "主观题三次不合格 "  // 主观题三次不合格
