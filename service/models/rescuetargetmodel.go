package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RescueTargetModel = (*customRescueTargetModel)(nil)

type (
	// RescueTargetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRescueTargetModel.
	RescueTargetModel interface {
		rescueTargetModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		CountClaimNumber(ctx context.Context, userId int64) (int64, error)
		FindOneByIdUserId(ctx context.Context, id, userId int64) (*RescueTarget, error)
		FindAllByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueTarget, error)
		FindOneByWeiboAddress(ctx context.Context, weiboAddress string) (*RescueTarget, error)
		FindMobileById(ctx context.Context, id int64) ([]string, error)
	}

	customRescueTargetModel struct {
		*defaultRescueTargetModel
	}
)

// NewRescueTargetModel returns a model for the database table.
func NewRescueTargetModel(conn sqlx.SqlConn) RescueTargetModel {
	return &customRescueTargetModel{
		defaultRescueTargetModel: newRescueTargetModel(conn),
	}
}

// 事务
func (m *defaultRescueTargetModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultRescueTargetModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(rescueTargetRows).From(m.table)
}

// 认领人数
func (m *defaultRescueTargetModel) CountClaimNumber(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select COUNT(*) AS total_count FROM %s WHERE (rescue_teacher1_id = ? OR rescue_teacher2_id = ? OR rescue_teacher3_id = ?) AND status = 1", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userId, userId, userId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 根据userId和救援对象id查找数据
func (m *defaultRescueTargetModel) FindOneByIdUserId(ctx context.Context, id, userId int64) (*RescueTarget, error) {
	var resp RescueTarget
	query := fmt.Sprintf("select %s from %s where (`rescue_teacher1_id` = ? or `rescue_teacher2_id` = ? or `rescue_teacher3_id` = ?) and id = ? and status = 1 limit 1", rescueTargetRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, userId, userId, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// 根据用户id寻找其认领的所有救援对象
func (m *defaultRescueTargetModel) FindAllByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueTarget, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("`rescue_teacher1_id` = ? or `rescue_teacher2_id` = ? or `rescue_teacher3_id` = ?", userId, userId, userId).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []RescueTarget
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// 根据微博地址查询一个救援对象，微博地址相同优先查最近的
func (m *defaultRescueTargetModel) FindOneByWeiboAddress(ctx context.Context, weiboAddress string) (*RescueTarget, error) {
	var resp RescueTarget
	query := fmt.Sprintf("select %s from %s where `weibo_address` = ? order by id desc limit 1", rescueTargetRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, weiboAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// 根据救援对象id返回所有救援老师的电话
func (m *defaultRescueTargetModel) FindMobileById(ctx context.Context, id int64) ([]string, error) {
	query := fmt.Sprintf("select user.mobile from rescue_target rt join user on user.id in (rt.rescue_teacher1_id, rt.rescue_teacher2_id, rt.rescue_teacher3_id) where rt.id = ?")
	var resp []string
	err := m.conn.QueryRowsCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
