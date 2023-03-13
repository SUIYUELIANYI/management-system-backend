package svc

import (
	"management-system/service/app/api/internal/config"
	"management-system/service/app/models"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	UsersModel     models.UsersModel
	UsersAuthModel models.UsersAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		UsersModel:     models.NewUsersModel(sqlx.NewMysql(c.DB.DataSource)),
		UsersAuthModel: models.NewUsersAuthModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
