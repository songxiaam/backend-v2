package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceAdminModel = (*customGovernanceAdminModel)(nil)

type (
	// GovernanceAdminModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceAdminModel.
	GovernanceAdminModel interface {
		governanceAdminModel
		withSession(session sqlx.Session) GovernanceAdminModel
	}

	customGovernanceAdminModel struct {
		*defaultGovernanceAdminModel
	}
)

// NewGovernanceAdminModel returns a model for the database table.
func NewGovernanceAdminModel(conn sqlx.SqlConn) GovernanceAdminModel {
	return &customGovernanceAdminModel{
		defaultGovernanceAdminModel: newGovernanceAdminModel(conn),
	}
}

func (m *customGovernanceAdminModel) withSession(session sqlx.Session) GovernanceAdminModel {
	return NewGovernanceAdminModel(sqlx.NewSqlConnFromSession(session))
}
