package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GovernanceSettingModel = (*customGovernanceSettingModel)(nil)

type (
	// GovernanceSettingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGovernanceSettingModel.
	GovernanceSettingModel interface {
		governanceSettingModel
		withSession(session sqlx.Session) GovernanceSettingModel
	}

	customGovernanceSettingModel struct {
		*defaultGovernanceSettingModel
	}
)

// NewGovernanceSettingModel returns a model for the database table.
func NewGovernanceSettingModel(conn sqlx.SqlConn) GovernanceSettingModel {
	return &customGovernanceSettingModel{
		defaultGovernanceSettingModel: newGovernanceSettingModel(conn),
	}
}

func (m *customGovernanceSettingModel) withSession(session sqlx.Session) GovernanceSettingModel {
	return NewGovernanceSettingModel(sqlx.NewSqlConnFromSession(session))
}
