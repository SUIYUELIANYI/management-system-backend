// Code generated by goctl. DO NOT EDIT.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	folderFieldNames          = builder.RawFieldNames(&Folder{})
	folderRows                = strings.Join(folderFieldNames, ",")
	folderRowsExpectAutoSet   = strings.Join(stringx.Remove(folderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	folderRowsWithPlaceHolder = strings.Join(stringx.Remove(folderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	folderModel interface {
		Insert(ctx context.Context, data *Folder) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Folder, error)
		Update(ctx context.Context, data *Folder) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFolderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Folder struct {
		Id         int64     `db:"id"` // 文件夹编号
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		FolderName string    `db:"folder_name"` // 文件名
		Role       int64     `db:"role"`        // 文件查看权限
	}
)

func newFolderModel(conn sqlx.SqlConn) *defaultFolderModel {
	return &defaultFolderModel{
		conn:  conn,
		table: "`folder`",
	}
}

func (m *defaultFolderModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultFolderModel) FindOne(ctx context.Context, id int64) (*Folder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", folderRows, m.table)
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

func (m *defaultFolderModel) Insert(ctx context.Context, data *Folder) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, folderRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.FolderName, data.Role)
	return ret, err
}

func (m *defaultFolderModel) Update(ctx context.Context, data *Folder) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, folderRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.FolderName, data.Role, data.Id)
	return err
}

func (m *defaultFolderModel) tableName() string {
	return m.table
}