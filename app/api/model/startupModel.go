package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"metaLand/app/api/internal/types"
)

var _ StartupModel = (*customStartupModel)(nil)

type (
	// StartupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupModel.
	StartupModel interface {
		startupModel
		withSession(session sqlx.Session) StartupModel
		List(ctx context.Context, comerID uint64, request *types.ListStartupsRequest) (resp []*types.Startup, total int, err error)
	}

	customStartupModel struct {
		*defaultStartupModel
	}
)

// NewStartupModel returns a model for the database table.
func NewStartupModel(conn sqlx.SqlConn) StartupModel {
	return &customStartupModel{
		defaultStartupModel: newStartupModel(conn),
	}
}

func (m *customStartupModel) withSession(session sqlx.Session) StartupModel {
	return NewStartupModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customStartupModel) List(ctx context.Context, comerID uint64, request *types.ListStartupsRequest) ([]*types.Startup, int, error) {
	where := ""
	if comerID > 0 {
		where = fmt.Sprintf("comer_id = %d ", comerID)
	}

	if request.Keyword != "" {
		if len(where) > 0 {
			where = fmt.Sprintf("%s AND name = like '%%%s%%' ", where, request.Keyword)
		} else {
			where = fmt.Sprintf("name like '%%%s%%' ", request.Keyword)
		}
	}

	if request.Mode != 0 {
		if len(where) > 0 {
			where = fmt.Sprintf("%s AND mode = %d ", where, request.Mode)
		} else {
			where = fmt.Sprintf("mode = %d ", request.Mode)
		}
	}

	var query string
	var total int
	if len(where) > 0 {
		query = fmt.Sprintf("select %s from %s where = %s limit ?,?", startupRows, m.table, where)
		err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf("select count(*) from %s where = %s", m.table, where))
		if err != nil {
			return nil, 0, err
		}
	} else {
		query = fmt.Sprintf("select %s from %s limit ?,?", startupRows, m.table)
		err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf("select count(*) from %s", m.table))
		if err != nil {
			return nil, 0, err
		}
	}

	logx.Info(1111)

	var resp []*types.Startup
	err := m.conn.QueryRowsCtx(ctx, &resp, query, request.Offset, request.Limit)
	return resp, total, err
}
