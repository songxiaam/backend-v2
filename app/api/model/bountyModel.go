package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyModel = (*customBountyModel)(nil)

type (
	// BountyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyModel.
	BountyModel interface {
		bountyModel
		withSession(session sqlx.Session) BountyModel
	}

	customBountyModel struct {
		*defaultBountyModel
	}
)

// NewBountyModel returns a model for the database table.
func NewBountyModel(conn sqlx.SqlConn) BountyModel {
	return &customBountyModel{
		defaultBountyModel: newBountyModel(conn),
	}
}

func (m *customBountyModel) withSession(session sqlx.Session) BountyModel {
	return NewBountyModel(sqlx.NewSqlConnFromSession(session))
}
