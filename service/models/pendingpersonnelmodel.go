package models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PendingPersonnelModel = (*customPendingPersonnelModel)(nil)

type (
	// PendingPersonnelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPendingPersonnelModel.
	PendingPersonnelModel interface {
		pendingPersonnelModel
	}

	customPendingPersonnelModel struct {
		*defaultPendingPersonnelModel
	}
)

// NewPendingPersonnelModel returns a model for the database table.
func NewPendingPersonnelModel(conn sqlx.SqlConn) PendingPersonnelModel {
	return &customPendingPersonnelModel{
		defaultPendingPersonnelModel: newPendingPersonnelModel(conn),
	}
}
