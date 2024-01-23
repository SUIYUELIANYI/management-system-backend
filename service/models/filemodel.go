package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FileModel = (*customFileModel)(nil)

type (
	// FileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileModel.
	FileModel interface {
		fileModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindOneNotDel(ctx context.Context, id int64) (*File, error)
		FindOneByFolderIdNotDel(ctx context.Context, folderId int64) (*File, error)
		FindAllByUserRole(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]File, error)
		FindAllByFolderId(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, folderId int64) ([]File, error)
		FindFileNameById(ctx context.Context, fileId int64) (string, error)
	}

	customFileModel struct {
		*defaultFileModel
	}
)

// NewFileModel returns a model for the database table.
func NewFileModel(conn sqlx.SqlConn) FileModel {
	return &customFileModel{
		defaultFileModel: newFileModel(conn),
	}
}

// 事务
func (m *defaultFileModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultFileModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(fileRows).From(m.table)
}

func (m *defaultFileModel) FindOneNotDel(ctx context.Context, id int64) (*File, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `del_state` = 0 limit 1", fileRows, m.table)
	var resp File
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

func (m *defaultFileModel) FindOneByFolderIdNotDel(ctx context.Context, folderId int64) (*File, error) {
	query := fmt.Sprintf("select %s from %s where `folder_id` = ? and `del_state` = 0 limit 1", fileRows, m.table)
	var resp File
	err := m.conn.QueryRowCtx(ctx, &resp, query, folderId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFileModel) FindAllByUserRole(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]File, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("role <= ? and `del_state` = 0", role).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []File
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultFileModel) FindAllByFolderId(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, folderId int64) ([]File, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("folder_id = ? and `del_state` = 0", folderId).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []File
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultFileModel) FindFileNameById(ctx context.Context, fileId int64) (string, error) {
	query := fmt.Sprintf("select file_name FROM %s WHERE id = ?", m.table)
	var fileName string
	err := m.conn.QueryRowCtx(ctx, &fileName, query, fileId)
	if err != nil {
		return "", err
	}
	return fileName, err
}