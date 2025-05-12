package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CrowdfundingInvestorModel = (*customCrowdfundingInvestorModel)(nil)

type (
	// CrowdfundingInvestorModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCrowdfundingInvestorModel.
	CrowdfundingInvestorModel interface {
		crowdfundingInvestorModel
		withSession(session sqlx.Session) CrowdfundingInvestorModel
	}

	customCrowdfundingInvestorModel struct {
		*defaultCrowdfundingInvestorModel
	}
)

// NewCrowdfundingInvestorModel returns a model for the database table.
func NewCrowdfundingInvestorModel(conn sqlx.SqlConn) CrowdfundingInvestorModel {
	return &customCrowdfundingInvestorModel{
		defaultCrowdfundingInvestorModel: newCrowdfundingInvestorModel(conn),
	}
}

func (m *customCrowdfundingInvestorModel) withSession(session sqlx.Session) CrowdfundingInvestorModel {
	return NewCrowdfundingInvestorModel(sqlx.NewSqlConnFromSession(session))
}
