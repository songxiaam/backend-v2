package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceProposalModel = (*customGovernanceProposalModel)(nil)

type (
	// GovernanceProposalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceProposalModel.
	GovernanceProposalModel interface {
		governanceProposalModel
		withSession(session sqlx.Session) GovernanceProposalModel
	}

	customGovernanceProposalModel struct {
		*defaultGovernanceProposalModel
	}
)

// NewGovernanceProposalModel returns a model for the database table.
func NewGovernanceProposalModel(conn sqlx.SqlConn) GovernanceProposalModel {
	return &customGovernanceProposalModel{
		defaultGovernanceProposalModel: newGovernanceProposalModel(conn),
	}
}

func (m *customGovernanceProposalModel) withSession(session sqlx.Session) GovernanceProposalModel {
	return NewGovernanceProposalModel(sqlx.NewSqlConnFromSession(session))
}
