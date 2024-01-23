package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	JwtAuthForRescue struct {
		AccessSecret string
		AccessExpire int64
	}
	DB struct {
		DataSource string
	}
	WxMiniConf struct {
		AppId     string `json:"AppId"`     // 微信AppId
		AppSecret string `json:"AppSecret"` // 微信AppSecret
	}
	Cache cache.CacheConf
	Redis redis.RedisConf
}
