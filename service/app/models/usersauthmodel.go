package models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UsersAuthModel = (*customUsersAuthModel)(nil)

type (
	// UsersAuthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersAuthModel.
	UsersAuthModel interface {
		usersAuthModel
	}

	customUsersAuthModel struct {
		*defaultUsersAuthModel
	}
)

// NewUsersAuthModel returns a model for the database table.
func NewUsersAuthModel(conn sqlx.SqlConn) UsersAuthModel {
	return &customUsersAuthModel{
		defaultUsersAuthModel: newUsersAuthModel(conn),
	}
}
