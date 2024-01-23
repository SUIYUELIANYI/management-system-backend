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
	signatureFieldNames          = builder.RawFieldNames(&Signature{})
	signatureRows                = strings.Join(signatureFieldNames, ",")
	signatureRowsExpectAutoSet   = strings.Join(stringx.Remove(signatureFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	signatureRowsWithPlaceHolder = strings.Join(stringx.Remove(signatureFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	signatureModel interface {
		Insert(ctx context.Context, data *Signature) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Signature, error)
		FindOneByRescueTeacherIdRescueTargetId(ctx context.Context, rescueTeacherId int64, rescueTargetId int64) (*Signature, error)
		Update(ctx context.Context, data *Signature) error
		Delete(ctx context.Context, id int64) error
	}

	defaultSignatureModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Signature struct {
		Id              int64          `db:"id"` // 签字编号
		CreateTime      time.Time      `db:"create_time"`
		UpdateTime      time.Time      `db:"update_time"`
		DeleteTime      time.Time      `db:"delete_time"`
		DelState        int64          `db:"del_state"`
		RescueTeacherId int64          `db:"rescue_teacher_id"` // 救援老师编号
		RescueTargetId  int64          `db:"rescue_target_id"`  // 救援对象编号
		Image           sql.NullString `db:"image"`             // 签字图片编码
	}
)

func newSignatureModel(conn sqlx.SqlConn) *defaultSignatureModel {
	return &defaultSignatureModel{
		conn:  conn,
		table: "`signature`",
	}
}

func (m *defaultSignatureModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultSignatureModel) FindOne(ctx context.Context, id int64) (*Signature, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", signatureRows, m.table)
	var resp Signature
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

func (m *defaultSignatureModel) FindOneByRescueTeacherIdRescueTargetId(ctx context.Context, rescueTeacherId int64, rescueTargetId int64) (*Signature, error) {
	var resp Signature
	query := fmt.Sprintf("select %s from %s where `rescue_teacher_id` = ? and `rescue_target_id` = ? limit 1", signatureRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, rescueTeacherId, rescueTargetId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSignatureModel) Insert(ctx context.Context, data *Signature) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, signatureRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.RescueTeacherId, data.RescueTargetId, data.Image)
	return ret, err
}

func (m *defaultSignatureModel) Update(ctx context.Context, newData *Signature) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, signatureRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.DeleteTime, newData.DelState, newData.RescueTeacherId, newData.RescueTargetId, newData.Image, newData.Id)
	return err
}

func (m *defaultSignatureModel) tableName() string {
	return m.table
}
