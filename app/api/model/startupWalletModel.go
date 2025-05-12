package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StartupWalletModel = (*customStartupWalletModel)(nil)

type (
	// StartupWalletModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStartupWalletModel.
	StartupWalletModel interface {
		startupWalletModel
		withSession(session sqlx.Session) StartupWalletModel
	}

	customStartupWalletModel struct {
		*defaultStartupWalletModel
	}
)

// NewStartupWalletModel returns a model for the database table.
func NewStartupWalletModel(conn sqlx.SqlConn) StartupWalletModel {
	return &customStartupWalletModel{
		defaultStartupWalletModel: newStartupWalletModel(conn),
	}
}

func (m *customStartupWalletModel) withSession(session sqlx.Session) StartupWalletModel {
	return NewStartupWalletModel(sqlx.NewSqlConnFromSession(session))
}
