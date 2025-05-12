package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceStrategyModel = (*customGovernanceStrategyModel)(nil)

type (
	// GovernanceStrategyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceStrategyModel.
	GovernanceStrategyModel interface {
		governanceStrategyModel
		withSession(session sqlx.Session) GovernanceStrategyModel
	}

	customGovernanceStrategyModel struct {
		*defaultGovernanceStrategyModel
	}
)

// NewGovernanceStrategyModel returns a model for the database table.
func NewGovernanceStrategyModel(conn sqlx.SqlConn) GovernanceStrategyModel {
	return &customGovernanceStrategyModel{
		defaultGovernanceStrategyModel: newGovernanceStrategyModel(conn),
	}
}

func (m *customGovernanceStrategyModel) withSession(session sqlx.Session) GovernanceStrategyModel {
	return NewGovernanceStrategyModel(sqlx.NewSqlConnFromSession(session))
}
