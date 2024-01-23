package models

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RescueInfoModel = (*customRescueInfoModel)(nil)

type Area struct {
	Name            string `json:"name"`
	RescueFrequency int64  `json:"rescue_frequency"`
}

type Year struct {
	Name            string `json:"name"`
	RescueFrequency int64  `json:"rescue_frequency"`
}

type (
	// RescueInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRescueInfoModel.
	RescueInfoModel interface {
		rescueInfoModel
		RowBuilder() squirrel.SelectBuilder
		TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]RescueInfo, error)
		FindClaimedByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueInfo, error)
		FindUnclaimedByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueInfo, error)
		FindAllByAddressWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string) ([]RescueInfo, error)
		FindAllByRescueTargetIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, Page, PageSize int64, rescueTargetId int64) ([]RescueInfo, error)
		FindAllByRescueTeacherNameWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, Page, PageSize int64, rescueTeacherName string) ([]RescueInfo, error)
		FindRescueFrequencyByArea(ctx context.Context) ([]Area, error)
		FindRescueFrequencyByYear(ctx context.Context) ([]Year, error)
		FindOneByReleaseTimeWeiboAddress(ctx context.Context, releaseTime string, weiboAddress string) (*RescueInfo, error) // 微博地址和发布时间作为唯一标志，避免重复数据
		FindOneByWeiboAddress(ctx context.Context, weiboAddress string) (*RescueInfo, error)
	}

	customRescueInfoModel struct {
		*defaultRescueInfoModel
	}
)

// NewRescueInfoModel returns a model for the database table.
func NewRescueInfoModel(conn sqlx.SqlConn) RescueInfoModel {
	return &customRescueInfoModel{
		defaultRescueInfoModel: newRescueInfoModel(conn),
	}
}

// 事务
func (m *defaultRescueInfoModel) TransCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}

func (m *defaultRescueInfoModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(rescueInfoRows).From(m.table)
}

func (m *defaultRescueInfoModel) FindAllWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64) ([]RescueInfo, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []RescueInfo
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultRescueInfoModel) FindAllByAddressWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, address string) ([]RescueInfo, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("LEFT(area,2) = ? and del_state = 0", address).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []RescueInfo
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultRescueInfoModel) FindAllByRescueTargetIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, rescueTargetId int64) ([]RescueInfo, error) {
	rowBuilder = rowBuilder.OrderBy("update_time DESC")
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query, values, err := rowBuilder.Where("rescue_target_id = ? and del_state = 0", rescueTargetId).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []RescueInfo
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// 根据救援老师姓名模糊查询救援信息
func (m *defaultRescueInfoModel) FindAllByRescueTeacherNameWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, rescueTeacherName string) ([]RescueInfo, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	// query := "SELECT ri.* FROM rescue_info ri JOIN rescue_target rt ON ri.rescue_target_id = rt.id JOIN user u1 ON rt.rescue_teacher1_id = u1.id JOIN user u2 ON rt.rescue_teacher2_id = u2.id JOIN user u3 ON rt.rescue_teacher3_id = u3.id WHERE u1.username LIKE CONCAT('%', ?, '%') OR u2.username LIKE CONCAT('%', ?, '%') OR u3.username LIKE CONCAT('%', ?, '%') LIMIT ? OFFSET ?"
	// 上面的query语句会导致只返回一条数据，具体原因未知。
	query := "SELECT ri.* FROM rescue_info ri JOIN rescue_target rt ON ri.rescue_target_id = rt.id JOIN user u1 ON rt.rescue_teacher1_id = u1.id WHERE u1.username LIKE CONCAT('%', ?, '%') union SELECT ri.* FROM rescue_info ri JOIN rescue_target rt ON ri.rescue_target_id = rt.id JOIN user u2 ON rt.rescue_teacher2_id = u2.id WHERE u2.username LIKE CONCAT('%', ?, '%') union SELECT ri.* FROM rescue_info ri JOIN rescue_target rt ON ri.rescue_target_id = rt.id JOIN user u3 ON rt.rescue_teacher3_id = u3.id WHERE u3.username LIKE CONCAT('%', ?, '%') order by update_time desc LIMIT ? OFFSET ?;"
	var resp []RescueInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query, rescueTeacherName, rescueTeacherName, rescueTeacherName, uint(pageSize), uint(offset)) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultRescueInfoModel) FindClaimedByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueInfo, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select rescue_info.`id`,rescue_info.`create_time`,rescue_info.`update_time`,rescue_info.`delete_time`,rescue_info.`del_state`,rescue_info.`rescue_target_id`,rescue_info.`release_time`,rescue_info.`weibo_account`,rescue_info.`weibo_address`,rescue_info.`nickname`,rescue_info.`risk_level`,rescue_info.`area`,rescue_info.`sex`,rescue_info.`birthday`,rescue_info.`brief_introduction`,rescue_info.`text` FROM rescue_info join rescue_target WHERE rescue_target_id = `rescue_target`.id and (rescue_teacher1_id = ? or rescue_teacher2_id = ? or rescue_teacher3_id = ? ) order by update_time desc LIMIT ? OFFSET ?;")
	var resp []RescueInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, userId, userId, uint(pageSize), uint(offset)) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// select `id`,`create_time`,`update_time`,`delete_time`,`del_state`,`rescue_target_id`,`release_time`,`weibo_account`,`weibo_address`,`nickname`,`risk_level`,`area`,`sex`,`birthday`,`brief_introduction`,`text` FROM rescue_info join rescue_target WHERE rescue_target_id = `rescue_target`.id and (resuce_teacher1_id = 1 or resuce_teacher2_id = 1 or resuce_teacher3_id = 1 )
func (m *defaultRescueInfoModel) FindUnclaimedByUserIdWithPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, userId int64) ([]RescueInfo, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select rescue_info.`id`,rescue_info.`create_time`,rescue_info.`update_time`,rescue_info.`delete_time`,rescue_info.`del_state`,rescue_info.`rescue_target_id`,rescue_info.`release_time`,rescue_info.`weibo_account`,rescue_info.`weibo_address`,rescue_info.`nickname`,rescue_info.`risk_level`,rescue_info.`area`,rescue_info.`sex`,rescue_info.`birthday`,rescue_info.`brief_introduction`,rescue_info.`text` FROM rescue_info join rescue_target WHERE rescue_target_id = `rescue_target`.id and (rescue_teacher1_id != ? and rescue_teacher2_id != ? and rescue_teacher3_id != ? ) order by update_time desc LIMIT ? OFFSET ? ;")
	var resp []RescueInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, userId, userId, uint(pageSize), uint(offset)) // 这里报错"unsupported unmarshal type"，注意是QueryRowsCtx，不是QueryRowCtx
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultRescueInfoModel) FindRescueFrequencyByArea(ctx context.Context) ([]Area, error) {
	query := "SELECT `area`, COUNT(*) AS rescue_frequency FROM rescue_info GROUP BY `area`"
	var rows []struct {
		Name            string `db:"area"` // 这里结构体标签db对应的是area，不能是name
		RescueFrequency int64  `db:"rescue_frequency"`
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	var areas []Area
	for _, row := range rows {
		area := Area{
			Name:            row.Name,
			RescueFrequency: row.RescueFrequency,
		}
		areas = append(areas, area)
	}
	return areas, nil
}

func (m *defaultRescueInfoModel) FindRescueFrequencyByYear(ctx context.Context) ([]Year, error) {
	query := "SELECT SUBSTRING(release_time, 1, 4) AS year, COUNT(*) AS rescue_frequency FROM rescue_info GROUP BY SUBSTRING(release_time, 1, 4) order by year asc;"
	var rows []struct {
		Name            string `db:"year"`
		RescueFrequency int64  `db:"rescue_frequency"`
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	var years []Year
	for _, row := range rows {

		frequency := Year{
			Name:            row.Name,
			RescueFrequency: row.RescueFrequency,
		}
		years = append(years, frequency)
	}
	return years, nil
}

func (m *defaultRescueInfoModel) FindOneByReleaseTimeWeiboAddress(ctx context.Context, releaseTime string, weiboAddress string) (*RescueInfo, error) {
	query := fmt.Sprintf("select %s from %s where `release_time` = ? and `weibo_address` = ? limit 1", rescueInfoRows, m.table)
	var resp RescueInfo
	err := m.conn.QueryRowCtx(ctx, &resp, query, releaseTime, weiboAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRescueInfoModel) FindOneByWeiboAddress(ctx context.Context, weiboAddress string) (*RescueInfo, error) {
	query := fmt.Sprintf("select %s from %s where `weibo_address` = ? order by `update_time` desc limit 1;", rescueInfoRows, m.table)
	var resp RescueInfo
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
