package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyContactModel = (*customBountyContactModel)(nil)

type (
	// BountyContactModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyContactModel.
	BountyContactModel interface {
		bountyContactModel
		withSession(session sqlx.Session) BountyContactModel
	}

	customBountyContactModel struct {
		*defaultBountyContactModel
	}
)

// NewBountyContactModel returns a model for the database table.
func NewBountyContactModel(conn sqlx.SqlConn) BountyContactModel {
	return &customBountyContactModel{
		defaultBountyContactModel: newBountyContactModel(conn),
	}
}

func (m *customBountyContactModel) withSession(session sqlx.Session) BountyContactModel {
	return NewBountyContactModel(sqlx.NewSqlConnFromSession(session))
}
