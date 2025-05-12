package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CrowdfundingModel = (*customCrowdfundingModel)(nil)

type (
	// CrowdfundingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCrowdfundingModel.
	CrowdfundingModel interface {
		crowdfundingModel
		withSession(session sqlx.Session) CrowdfundingModel
	}

	customCrowdfundingModel struct {
		*defaultCrowdfundingModel
	}
)

// NewCrowdfundingModel returns a model for the database table.
func NewCrowdfundingModel(conn sqlx.SqlConn) CrowdfundingModel {
	return &customCrowdfundingModel{
		defaultCrowdfundingModel: newCrowdfundingModel(conn),
	}
}

func (m *customCrowdfundingModel) withSession(session sqlx.Session) CrowdfundingModel {
	return NewCrowdfundingModel(sqlx.NewSqlConnFromSession(session))
}
