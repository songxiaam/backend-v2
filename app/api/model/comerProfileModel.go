package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ComerProfileModel = (*customComerProfileModel)(nil)

type (
	// ComerProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customComerProfileModel.
	ComerProfileModel interface {
		comerProfileModel
		withSession(session sqlx.Session) ComerProfileModel
	}

	customComerProfileModel struct {
		*defaultComerProfileModel
	}
)

// NewComerProfileModel returns a model for the database table.
func NewComerProfileModel(conn sqlx.SqlConn) ComerProfileModel {
	return &customComerProfileModel{
		defaultComerProfileModel: newComerProfileModel(conn),
	}
}

func (m *customComerProfileModel) withSession(session sqlx.Session) ComerProfileModel {
	return NewComerProfileModel(sqlx.NewSqlConnFromSession(session))
}
