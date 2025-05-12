package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ChainEndpointModel = (*customChainEndpointModel)(nil)

type (
	// ChainEndpointModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChainEndpointModel.
	ChainEndpointModel interface {
		chainEndpointModel
		withSession(session sqlx.Session) ChainEndpointModel
	}

	customChainEndpointModel struct {
		*defaultChainEndpointModel
	}
)

// NewChainEndpointModel returns a model for the database table.
func NewChainEndpointModel(conn sqlx.SqlConn) ChainEndpointModel {
	return &customChainEndpointModel{
		defaultChainEndpointModel: newChainEndpointModel(conn),
	}
}

func (m *customChainEndpointModel) withSession(session sqlx.Session) ChainEndpointModel {
	return NewChainEndpointModel(sqlx.NewSqlConnFromSession(session))
}
