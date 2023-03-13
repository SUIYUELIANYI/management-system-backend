package models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

var UserAuthTypeSystem string ="system" // 平台内部手机密码登录
var UserAuthTypeSmallWX string = "wxMini" // 微信小程序登录
