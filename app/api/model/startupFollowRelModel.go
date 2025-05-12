package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StartupFollowRelModel = (*customStartupFollowRelModel)(nil)

type (
	// StartupFollowRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupFollowRelModel.
	StartupFollowRelModel interface {
		startupFollowRelModel
		withSession(session sqlx.Session) StartupFollowRelModel
	}

	customStartupFollowRelModel struct {
		*defaultStartupFollowRelModel
	}
)

// NewStartupFollowRelModel returns a model for the database table.
func NewStartupFollowRelModel(conn sqlx.SqlConn) StartupFollowRelModel {
	return &customStartupFollowRelModel{
		defaultStartupFollowRelModel: newStartupFollowRelModel(conn),
	}
}

func (m *customStartupFollowRelModel) withSession(session sqlx.Session) StartupFollowRelModel {
	return NewStartupFollowRelModel(sqlx.NewSqlConnFromSession(session))
}
