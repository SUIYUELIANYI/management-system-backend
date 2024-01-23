package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]User, error)
		FindAllByRoleWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]User, error)
		UpdateWithRole(ctx context.Context, userId, role int64) error
		FindUsernameById(ctx context.Context, userId int64) (string, error)
		FindRoleById(ctx context.Context, userId int64) (int64, error)
		FindAllMoblie(ctx context.Context) ([]string, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

// 事务
func (m *defaultUserModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultUserModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(userRows).From(m.table)
}

func (m *defaultUserModel) FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]User, error) {
	rowBuilder = rowBuilder.OrderBy("role DESC, id ASC") // 按权限排序，先查询到最高权限的
	// 2023.11.3 问题：使用过程中数据会出现重复和丢失的问题
	// 经过查阅资料，要避免在分页查询中出现重复数据，可以尝试添加一个额外的排序条件，以确保结果的唯一性，这里添加了id
	if page < 1 { // 判断页码是否小于1，如果小于1，设置为1
		page = 1
	}

	offset := (page - 1) * pageSize // 当页的偏移量
	query, values, err := rowBuilder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []User
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) UpdateWithRole(ctx context.Context, userId, role int64) error {
	query := fmt.Sprintf("update %s set `role` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, role, userId)
	return err
}

func (m *defaultUserModel) FindAllByRoleWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, role int64) ([]User, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("role = ?", role).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []User
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindUsernameById(ctx context.Context, userId int64) (string, error) {
	query := fmt.Sprintf("select username FROM %s WHERE id = ?", m.table)
	var username string
	err := m.conn.QueryRowCtx(ctx, &username, query, userId)
	if err != nil {
		return "", err
	}
	return username, err
}

func (m *defaultUserModel) FindRoleById(ctx context.Context, userId int64) (int64, error) {
	query := "select `role` FROM user WHERE `id` = ?"
	var role int64
	err := m.conn.QueryRowCtx(ctx, &role, query, userId)
	if err != nil {
		return 0, err
	}
	return role, err
}

// 返回所有用户的电话列表
func (m *defaultUserModel) FindAllMoblie(ctx context.Context) ([]string, error) {
	query := "select `mobile` from `user` where del_state = 0"
	var resp []string
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
