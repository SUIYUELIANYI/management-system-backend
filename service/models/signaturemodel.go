package models

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SignatureModel = (*customSignatureModel)(nil)

type (
	// SignatureModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSignatureModel.
	SignatureModel interface {
		signatureModel
		RowBuilder() squirrel.SelectBuilder
		CountByRescueTargetId(ctx context.Context, rescueTargetId int64) (int64, error)
		FindAllByRescueTargetId(ctx context.Context, rowBuilder squirrel.SelectBuilder, rescueTargetId int64) ([]Signature, error)
	}

	customSignatureModel struct {
		*defaultSignatureModel
	}
)

// NewSignatureModel returns a model for the database table.
func NewSignatureModel(conn sqlx.SqlConn) SignatureModel {
	return &customSignatureModel{
		defaultSignatureModel: newSignatureModel(conn),
	}
}

func (m *defaultSignatureModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(signatureRows).From(m.table)
}

// 查询救援对象对应的签字人数
func (m *defaultSignatureModel) CountByRescueTargetId(ctx context.Context, rescueTargetId int64) (int64, error) {
	query := "select COUNT(*) AS number FROM signature WHERE rescue_target_id = ?"
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, rescueTargetId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询救援对象对应的签字人数
func (m *defaultSignatureModel) FindAllByRescueTargetId(ctx context.Context, rowBuilder squirrel.SelectBuilder, rescueTargetId int64) ([]Signature, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	query, values, err := rowBuilder.Where("rescue_target_id = ? and del_state = 0", rescueTargetId).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []Signature
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
