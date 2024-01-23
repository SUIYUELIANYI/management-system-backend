package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RescueProcessModel = (*customRescueProcessModel)(nil)

type (
	// RescueProcessModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRescueProcessModel.
	RescueProcessModel interface {
		rescueProcessModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindAllByRescueInfoIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, rescueInfoId int64) ([]RescueProcess, error)
		FindAllByRescueTeacherId(ctx context.Context, rowBuilder squirrel.SelectBuilder, rescueTeacherId int64) ([]RescueProcess, error)
		FindOneByInfoIdTeacherId(ctx context.Context, rescueInfoId int64, rescueTeacherId int64) (*RescueProcess, error)
		CountRescueFrequency(ctx context.Context, userId int64) (int64, error)
	}

	customRescueProcessModel struct {
		*defaultRescueProcessModel
	}
)

// NewRescueProcessModel returns a model for the database table.
func NewRescueProcessModel(conn sqlx.SqlConn) RescueProcessModel {
	return &customRescueProcessModel{
		defaultRescueProcessModel: newRescueProcessModel(conn),
	}
}

// 事务
func (m *defaultRescueProcessModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultRescueProcessModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(rescueProcessRows).From(m.table)
}

// 查询救援信息的所有救援评价
func (m *defaultRescueProcessModel) FindAllByRescueInfoIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, rescueInfoId int64) ([]RescueProcess, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("rescue_info_id = ? and del_state = 0", rescueInfoId).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []RescueProcess
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// 查询救援老师的所有救援评价
func (m *defaultRescueProcessModel) FindAllByRescueTeacherId(ctx context.Context, rowBuilder squirrel.SelectBuilder, rescueTeacherId int64) ([]RescueProcess, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	query, values, err := rowBuilder.Where("rescue_teacher_id = ? and del_state = 0", rescueTeacherId).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []RescueProcess
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultRescueProcessModel) FindOneByInfoIdTeacherId(ctx context.Context, rescueInfoId int64, rescueTeacherId int64) (*RescueProcess, error) {
	query := fmt.Sprintf("select %s from %s where `rescue_info_id` = ? and `rescue_teacher_id` = ? and `del_state` = 0 limit 1", rescueProcessRows, m.table)
	var resp RescueProcess
	err := m.conn.QueryRowCtx(ctx, &resp, query, rescueInfoId, rescueTeacherId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// 查询用户救援次数
func (m *defaultRescueProcessModel) CountRescueFrequency(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select COUNT(*) AS total_count FROM %s WHERE `rescue_teacher_id` = ?  and `del_state` = 0", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
