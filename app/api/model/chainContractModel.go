package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ChainContractModel = (*customChainContractModel)(nil)

type (
	// ChainContractModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChainContractModel.
	ChainContractModel interface {
		chainContractModel
		withSession(session sqlx.Session) ChainContractModel
	}

	customChainContractModel struct {
		*defaultChainContractModel
	}
)

// NewChainContractModel returns a model for the database table.
func NewChainContractModel(conn sqlx.SqlConn) ChainContractModel {
	return &customChainContractModel{
		defaultChainContractModel: newChainContractModel(conn),
	}
}

func (m *customChainContractModel) withSession(session sqlx.Session) ChainContractModel {
	return NewChainContractModel(sqlx.NewSqlConnFromSession(session))
}
