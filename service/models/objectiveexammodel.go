package models

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ObjectiveExamModel = (*customObjectiveExamModel)(nil)

type (
	// ObjectiveExamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customObjectiveExamModel.
	ObjectiveExamModel interface {
		objectiveExamModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindOneByUserIdResult(ctx context.Context, userId int64, result int64) (*ObjectiveExam, error)
		FindAllByUserId(ctx context.Context, rowBuilder squirrel.SelectBuilder, userId int64) ([]ObjectiveExam, error)
		DeleteByUserId(ctx context.Context, userId int64) error
	}

	customObjectiveExamModel struct {
		*defaultObjectiveExamModel
	}
)

// NewObjectiveExamModel returns a model for the database table.
func NewObjectiveExamModel(conn sqlx.SqlConn) ObjectiveExamModel {
	return &customObjectiveExamModel{
		defaultObjectiveExamModel: newObjectiveExamModel(conn),
	}
}

// 事务
func (m *defaultObjectiveExamModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultObjectiveExamModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(objectiveExamRows).From(m.table)
}

func (m *defaultObjectiveExamModel) FindOneByUserIdResult(ctx context.Context, userId int64, result int64) (*ObjectiveExam, error) {
	var resp ObjectiveExam
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `result` = ? and del_state = 0 limit 1", objectiveExamRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, result)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultObjectiveExamModel) FindAllByUserId(ctx context.Context, rowBuilder squirrel.SelectBuilder, userId int64) ([]ObjectiveExam, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的
	query, values, err := rowBuilder.Where("user_id = ? and del_state = 0", userId).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []ObjectiveExam
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultObjectiveExamModel) DeleteByUserId(ctx context.Context, userId int64) error {
	deleteTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("update %s set `delete_time` = ?,`del_state` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, deleteTime, 1, userId)
	return err
}
