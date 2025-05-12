package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyDepositModel = (*customBountyDepositModel)(nil)

type (
	// BountyDepositModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyDepositModel.
	BountyDepositModel interface {
		bountyDepositModel
		withSession(session sqlx.Session) BountyDepositModel
	}

	customBountyDepositModel struct {
		*defaultBountyDepositModel
	}
)

// NewBountyDepositModel returns a model for the database table.
func NewBountyDepositModel(conn sqlx.SqlConn) BountyDepositModel {
	return &customBountyDepositModel{
		defaultBountyDepositModel: newBountyDepositModel(conn),
	}
}

func (m *customBountyDepositModel) withSession(session sqlx.Session) BountyDepositModel {
	return NewBountyDepositModel(sqlx.NewSqlConnFromSession(session))
}
