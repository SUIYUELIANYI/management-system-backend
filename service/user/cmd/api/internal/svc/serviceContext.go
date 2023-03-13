package svc

import (
	"management-system/service/user/cmd/api/internal/config"
	"management-system/service/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UsersModel model.UsersModel
	UsersAuthModel model.UsersAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UsersModel: model.NewUsersModel(sqlx.NewMysql(c.DB.DataSource)),
		UsersAuthModel: model.NewUsersAuthModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
