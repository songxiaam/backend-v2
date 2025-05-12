package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceVoteModel = (*customGovernanceVoteModel)(nil)

type (
	// GovernanceVoteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceVoteModel.
	GovernanceVoteModel interface {
		governanceVoteModel
		withSession(session sqlx.Session) GovernanceVoteModel
	}

	customGovernanceVoteModel struct {
		*defaultGovernanceVoteModel
	}
)

// NewGovernanceVoteModel returns a model for the database table.
func NewGovernanceVoteModel(conn sqlx.SqlConn) GovernanceVoteModel {
	return &customGovernanceVoteModel{
		defaultGovernanceVoteModel: newGovernanceVoteModel(conn),
	}
}

func (m *customGovernanceVoteModel) withSession(session sqlx.Session) GovernanceVoteModel {
	return NewGovernanceVoteModel(sqlx.NewSqlConnFromSession(session))
}
