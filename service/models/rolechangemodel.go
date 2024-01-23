package models

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleChangeModel = (*customRoleChangeModel)(nil)

type (
	// RoleChangeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleChangeModel.
	RoleChangeModel interface {
		roleChangeModel
		FindAllByUserId(ctx context.Context, userId int64) ([]RoleChange, error)
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
	}

	customRoleChangeModel struct {
		*defaultRoleChangeModel
	}
)

// NewRoleChangeModel returns a model for the database table.
func NewRoleChangeModel(conn sqlx.SqlConn) RoleChangeModel {
	return &customRoleChangeModel{
		defaultRoleChangeModel: newRoleChangeModel(conn),
	}
}

// 事务
func (m *defaultRoleChangeModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

// 查询user_id对应用户的所有身份更改记录
func (m *defaultRoleChangeModel) FindAllByUserId(ctx context.Context, userId int64) ([]RoleChange, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ?", roleChangeRows, m.table) // %s为string类型或[]byte的占位符
	var resp []RoleChange
	err := m.conn.QueryRowsCtx(context.Background(), &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, sqlx.ErrNotFound
	default:
		return nil, err
	}
}
