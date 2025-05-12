package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TransactionModel = (*customTransactionModel)(nil)

type (
	// TransactionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionModel.
	TransactionModel interface {
		transactionModel
		withSession(session sqlx.Session) TransactionModel
	}

	customTransactionModel struct {
		*defaultTransactionModel
	}
)

// NewTransactionModel returns a model for the database table.
func NewTransactionModel(conn sqlx.SqlConn) TransactionModel {
	return &customTransactionModel{
		defaultTransactionModel: newTransactionModel(conn),
	}
}

func (m *customTransactionModel) withSession(session sqlx.Session) TransactionModel {
	return NewTransactionModel(sqlx.NewSqlConnFromSession(session))
}
