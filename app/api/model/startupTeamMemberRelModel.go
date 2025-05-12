package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StartupTeamMemberRelModel = (*customStartupTeamMemberRelModel)(nil)

type (
	// StartupTeamMemberRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupTeamMemberRelModel.
	StartupTeamMemberRelModel interface {
		startupTeamMemberRelModel
		withSession(session sqlx.Session) StartupTeamMemberRelModel
	}

	customStartupTeamMemberRelModel struct {
		*defaultStartupTeamMemberRelModel
	}
)

// NewStartupTeamMemberRelModel returns a model for the database table.
func NewStartupTeamMemberRelModel(conn sqlx.SqlConn) StartupTeamMemberRelModel {
	return &customStartupTeamMemberRelModel{
		defaultStartupTeamMemberRelModel: newStartupTeamMemberRelModel(conn),
	}
}

func (m *customStartupTeamMemberRelModel) withSession(session sqlx.Session) StartupTeamMemberRelModel {
	return NewStartupTeamMemberRelModel(sqlx.NewSqlConnFromSession(session))
}
