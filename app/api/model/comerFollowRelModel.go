package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ComerFollowRelModel = (*customComerFollowRelModel)(nil)

type (
	// ComerFollowRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customComerFollowRelModel.
	ComerFollowRelModel interface {
		comerFollowRelModel
		withSession(session sqlx.Session) ComerFollowRelModel
	}

	customComerFollowRelModel struct {
		*defaultComerFollowRelModel
	}
)

// NewComerFollowRelModel returns a model for the database table.
func NewComerFollowRelModel(conn sqlx.SqlConn) ComerFollowRelModel {
	return &customComerFollowRelModel{
		defaultComerFollowRelModel: newComerFollowRelModel(conn),
	}
}

func (m *customComerFollowRelModel) withSession(session sqlx.Session) ComerFollowRelModel {
	return NewComerFollowRelModel(sqlx.NewSqlConnFromSession(session))
}
