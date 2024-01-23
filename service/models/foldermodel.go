package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FolderModel = (*customFolderModel)(nil)

type (
	// FolderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFolderModel.
	FolderModel interface {
		folderModel
		RowBuilder() squirrel.SelectBuilder
		FindOneNotDel(ctx context.Context, id int64) (*Folder, error)
		FindAllByUserRole(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]Folder, error)
	}

	customFolderModel struct {
		*defaultFolderModel
	}
)

// NewFolderModel returns a model for the database table.
func NewFolderModel(conn sqlx.SqlConn) FolderModel {
	return &customFolderModel{
		defaultFolderModel: newFolderModel(conn),
	}
}

func (m *defaultFolderModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(folderRows).From(m.table)
}

func (m *defaultFolderModel) FindOneNotDel(ctx context.Context, id int64) (*Folder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `del_state` = 0 limit 1", folderRows, m.table)
	var resp Folder
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

func (m *defaultFolderModel) FindAllByUserRole(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]Folder, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("role <= ? and `del_state` = 0", role).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []Folder
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
