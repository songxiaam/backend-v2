package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyPaymentTermsModel = (*customBountyPaymentTermsModel)(nil)

type (
	// BountyPaymentTermsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyPaymentTermsModel.
	BountyPaymentTermsModel interface {
		bountyPaymentTermsModel
		withSession(session sqlx.Session) BountyPaymentTermsModel
	}

	customBountyPaymentTermsModel struct {
		*defaultBountyPaymentTermsModel
	}
)

// NewBountyPaymentTermsModel returns a model for the database table.
func NewBountyPaymentTermsModel(conn sqlx.SqlConn) BountyPaymentTermsModel {
	return &customBountyPaymentTermsModel{
		defaultBountyPaymentTermsModel: newBountyPaymentTermsModel(conn),
	}
}

func (m *customBountyPaymentTermsModel) withSession(session sqlx.Session) BountyPaymentTermsModel {
	return NewBountyPaymentTermsModel(sqlx.NewSqlConnFromSession(session))
}
