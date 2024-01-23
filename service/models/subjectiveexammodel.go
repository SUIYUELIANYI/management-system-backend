package models

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubjectiveExamModel = (*customSubjectiveExamModel)(nil)

type (
	// SubjectiveExamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubjectiveExamModel.
	SubjectiveExamModel interface {
		subjectiveExamModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindOneByUserIdResult(ctx context.Context, userId int64, result int64) (*SubjectiveExam, error)
		CountByUserIdResult(ctx context.Context, userId int64, result int64) (int64, error)
		FindAllByUserId(ctx context.Context, rowBuilder squirrel.SelectBuilder, userId int64) ([]SubjectiveExam, error)
		DeleteByUserId(ctx context.Context, userId int64) error
	}

	customSubjectiveExamModel struct {
		*defaultSubjectiveExamModel
	}
)

// NewSubjectiveExamModel returns a model for the database table.
func NewSubjectiveExamModel(conn sqlx.SqlConn) SubjectiveExamModel {
	return &customSubjectiveExamModel{
		defaultSubjectiveExamModel: newSubjectiveExamModel(conn),
	}
}

// 事务
func (m *defaultSubjectiveExamModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultSubjectiveExamModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(subjectiveExamRows).From(m.table)
}

func (m *defaultSubjectiveExamModel) FindOneByUserIdResult(ctx context.Context, userId int64, result int64) (*SubjectiveExam, error) {
	var resp SubjectiveExam
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `result` = ? and `del_state` = 0 limit 1", subjectiveExamRows, m.table)
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

func (m *defaultSubjectiveExamModel) CountByUserIdResult(ctx context.Context, userId int64, result int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `result` = ? and `del_state` = 0", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, result)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultSubjectiveExamModel) FindAllByUserId(ctx context.Context, rowBuilder squirrel.SelectBuilder, userId int64) ([]SubjectiveExam, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC") // 按更新时间排序，先查询到最先更新的
	query, values, err := rowBuilder.Where("user_id = ? and `del_state` = 0", userId).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []SubjectiveExam
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSubjectiveExamModel) DeleteByUserId(ctx context.Context, userId int64) error {
	deleteTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("update %s set `delete_time` = ?,`del_state` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, deleteTime, 1, userId)
	return err
}
