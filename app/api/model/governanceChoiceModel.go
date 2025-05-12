package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceChoiceModel = (*customGovernanceChoiceModel)(nil)

type (
	// GovernanceChoiceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceChoiceModel.
	GovernanceChoiceModel interface {
		governanceChoiceModel
		withSession(session sqlx.Session) GovernanceChoiceModel
	}

	customGovernanceChoiceModel struct {
		*defaultGovernanceChoiceModel
	}
)

// NewGovernanceChoiceModel returns a model for the database table.
func NewGovernanceChoiceModel(conn sqlx.SqlConn) GovernanceChoiceModel {
	return &customGovernanceChoiceModel{
		defaultGovernanceChoiceModel: newGovernanceChoiceModel(conn),
	}
}

func (m *customGovernanceChoiceModel) withSession(session sqlx.Session) GovernanceChoiceModel {
	return NewGovernanceChoiceModel(sqlx.NewSqlConnFromSession(session))
}
