package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StartupGroupModel = (*customStartupGroupModel)(nil)

type (
	// StartupGroupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupGroupModel.
	StartupGroupModel interface {
		startupGroupModel
		withSession(session sqlx.Session) StartupGroupModel
	}

	customStartupGroupModel struct {
		*defaultStartupGroupModel
	}
)

// NewStartupGroupModel returns a model for the database table.
func NewStartupGroupModel(conn sqlx.SqlConn) StartupGroupModel {
	return &customStartupGroupModel{
		defaultStartupGroupModel: newStartupGroupModel(conn),
	}
}

func (m *customStartupGroupModel) withSession(session sqlx.Session) StartupGroupModel {
	return NewStartupGroupModel(sqlx.NewSqlConnFromSession(session))
}
