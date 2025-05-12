package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ComerModel = (*customComerModel)(nil)

type (
	// ComerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customComerModel.
	ComerModel interface {
		comerModel
		withSession(session sqlx.Session) ComerModel
	}

	customComerModel struct {
		*defaultComerModel
	}
)

// NewComerModel returns a model for the database table.
func NewComerModel(conn sqlx.SqlConn) ComerModel {
	return &customComerModel{
		defaultComerModel: newComerModel(conn),
	}
}

func (m *customComerModel) withSession(session sqlx.Session) ComerModel {
	return NewComerModel(sqlx.NewSqlConnFromSession(session))
}
