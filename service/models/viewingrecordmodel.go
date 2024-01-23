package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ViewingRecordModel = (*customViewingRecordModel)(nil)

type (
	// ViewingRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customViewingRecordModel.
	ViewingRecordModel interface {
		viewingRecordModel
		RowBuilder() squirrel.SelectBuilder
		FindOneNotDel(ctx context.Context, id int64) (*ViewingRecord, error)
		FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]ViewingRecord, error)
	}

	customViewingRecordModel struct {
		*defaultViewingRecordModel
	}
)

// NewViewingRecordModel returns a model for the database table.
func NewViewingRecordModel(conn sqlx.SqlConn) ViewingRecordModel {
	return &customViewingRecordModel{
		defaultViewingRecordModel: newViewingRecordModel(conn),
	}
}

func (m *defaultViewingRecordModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(viewingRecordRows).From(m.table)
}

func (m *defaultViewingRecordModel) FindOneNotDel(ctx context.Context, id int64) (*ViewingRecord, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `del_state` = 0 limit 1 ", viewingRecordRows, m.table)
	var resp ViewingRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultViewingRecordModel) FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]ViewingRecord, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("`del_state` = 0").Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []ViewingRecord
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
