package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ComerAccountModel = (*customComerAccountModel)(nil)

type (
	// ComerAccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customComerAccountModel.
	ComerAccountModel interface {
		comerAccountModel
		withSession(session sqlx.Session) ComerAccountModel
	}

	customComerAccountModel struct {
		*defaultComerAccountModel
	}
)

// NewComerAccountModel returns a model for the database table.
func NewComerAccountModel(conn sqlx.SqlConn) ComerAccountModel {
	return &customComerAccountModel{
		defaultComerAccountModel: newComerAccountModel(conn),
	}
}

func (m *customComerAccountModel) withSession(session sqlx.Session) ComerAccountModel {
	return NewComerAccountModel(sqlx.NewSqlConnFromSession(session))
}
