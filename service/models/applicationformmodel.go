package models

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApplicationFormModel = (*customApplicationFormModel)(nil)

type (
	// ApplicationFormModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApplicationFormModel.
	ApplicationFormModel interface {
		applicationFormModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]ApplicationForm, error)
		FindAllByStatusWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, status string) ([]ApplicationForm, error)
		FindAllByAddressWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string) ([]ApplicationForm, error)
		FindAllByAddressAndStatusWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string, status int) ([]ApplicationForm, error)
	}

	customApplicationFormModel struct {
		*defaultApplicationFormModel
	}
)

// NewApplicationFormModel returns a model for the database table.
func NewApplicationFormModel(conn sqlx.SqlConn) ApplicationFormModel {
	return &customApplicationFormModel{
		defaultApplicationFormModel: newApplicationFormModel(conn),
	}
}

// 事务
func (m *defaultApplicationFormModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultApplicationFormModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(applicationFormRows).From(m.table)
}

func (m *defaultApplicationFormModel) FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]ApplicationForm, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []ApplicationForm
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultApplicationFormModel) FindAllByStatusWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, status string) ([]ApplicationForm, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where(status).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []ApplicationForm
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultApplicationFormModel) FindAllByAddressWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string) ([]ApplicationForm, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("LEFT(address,2)=?", address).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []ApplicationForm
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// 2023.12.12
// 将对申请表的筛选放入sql查询中
func (m *defaultApplicationFormModel) FindAllByAddressAndStatusWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string, status int) ([]ApplicationForm, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("status=? and LEFT(address,2)=?", status, address).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []ApplicationForm
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
