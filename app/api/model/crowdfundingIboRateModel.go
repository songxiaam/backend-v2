package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CrowdfundingIboRateModel = (*customCrowdfundingIboRateModel)(nil)

type (
	// CrowdfundingIboRateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCrowdfundingIboRateModel.
	CrowdfundingIboRateModel interface {
		crowdfundingIboRateModel
		withSession(session sqlx.Session) CrowdfundingIboRateModel
	}

	customCrowdfundingIboRateModel struct {
		*defaultCrowdfundingIboRateModel
	}
)

// NewCrowdfundingIboRateModel returns a model for the database table.
func NewCrowdfundingIboRateModel(conn sqlx.SqlConn) CrowdfundingIboRateModel {
	return &customCrowdfundingIboRateModel{
		defaultCrowdfundingIboRateModel: newCrowdfundingIboRateModel(conn),
	}
}

func (m *customCrowdfundingIboRateModel) withSession(session sqlx.Session) CrowdfundingIboRateModel {
	return NewCrowdfundingIboRateModel(sqlx.NewSqlConnFromSession(session))
}
