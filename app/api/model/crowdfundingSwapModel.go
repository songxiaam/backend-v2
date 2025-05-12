package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CrowdfundingSwapModel = (*customCrowdfundingSwapModel)(nil)

type (
	// CrowdfundingSwapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCrowdfundingSwapModel.
	CrowdfundingSwapModel interface {
		crowdfundingSwapModel
		withSession(session sqlx.Session) CrowdfundingSwapModel
	}

	customCrowdfundingSwapModel struct {
		*defaultCrowdfundingSwapModel
	}
)

// NewCrowdfundingSwapModel returns a model for the database table.
func NewCrowdfundingSwapModel(conn sqlx.SqlConn) CrowdfundingSwapModel {
	return &customCrowdfundingSwapModel{
		defaultCrowdfundingSwapModel: newCrowdfundingSwapModel(conn),
	}
}

func (m *customCrowdfundingSwapModel) withSession(session sqlx.Session) CrowdfundingSwapModel {
	return NewCrowdfundingSwapModel(sqlx.NewSqlConnFromSession(session))
}
