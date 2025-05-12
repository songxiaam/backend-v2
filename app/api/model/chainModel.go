package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ChainModel = (*customChainModel)(nil)

type (
	// ChainModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChainModel.
	ChainModel interface {
		chainModel
		withSession(session sqlx.Session) ChainModel
	}

	customChainModel struct {
		*defaultChainModel
	}
)

// NewChainModel returns a model for the database table.
func NewChainModel(conn sqlx.SqlConn) ChainModel {
	return &customChainModel{
		defaultChainModel: newChainModel(conn),
	}
}

func (m *customChainModel) withSession(session sqlx.Session) ChainModel {
	return NewChainModel(sqlx.NewSqlConnFromSession(session))
}
