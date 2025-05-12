package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StartupGroupMemberRelModel = (*customStartupGroupMemberRelModel)(nil)

type (
	// StartupGroupMemberRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupGroupMemberRelModel.
	StartupGroupMemberRelModel interface {
		startupGroupMemberRelModel
		withSession(session sqlx.Session) StartupGroupMemberRelModel
	}

	customStartupGroupMemberRelModel struct {
		*defaultStartupGroupMemberRelModel
	}
)

// NewStartupGroupMemberRelModel returns a model for the database table.
func NewStartupGroupMemberRelModel(conn sqlx.SqlConn) StartupGroupMemberRelModel {
	return &customStartupGroupMemberRelModel{
		defaultStartupGroupMemberRelModel: newStartupGroupMemberRelModel(conn),
	}
}

func (m *customStartupGroupMemberRelModel) withSession(session sqlx.Session) StartupGroupMemberRelModel {
	return NewStartupGroupMemberRelModel(sqlx.NewSqlConnFromSession(session))
}
