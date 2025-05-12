package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyPaymentPeriodModel = (*customBountyPaymentPeriodModel)(nil)

type (
	// BountyPaymentPeriodModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyPaymentPeriodModel.
	BountyPaymentPeriodModel interface {
		bountyPaymentPeriodModel
		withSession(session sqlx.Session) BountyPaymentPeriodModel
	}

	customBountyPaymentPeriodModel struct {
		*defaultBountyPaymentPeriodModel
	}
)

// NewBountyPaymentPeriodModel returns a model for the database table.
func NewBountyPaymentPeriodModel(conn sqlx.SqlConn) BountyPaymentPeriodModel {
	return &customBountyPaymentPeriodModel{
		defaultBountyPaymentPeriodModel: newBountyPaymentPeriodModel(conn),
	}
}

func (m *customBountyPaymentPeriodModel) withSession(session sqlx.Session) BountyPaymentPeriodModel {
	return NewBountyPaymentPeriodModel(sqlx.NewSqlConnFromSession(session))
}
