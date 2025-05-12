package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TagTargetRelModel = (*customTagTargetRelModel)(nil)

type (
	// TagTargetRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagTargetRelModel.
	TagTargetRelModel interface {
		tagTargetRelModel
		withSession(session sqlx.Session) TagTargetRelModel
	}

	customTagTargetRelModel struct {
		*defaultTagTargetRelModel
	}
)

// NewTagTargetRelModel returns a model for the database table.
func NewTagTargetRelModel(conn sqlx.SqlConn) TagTargetRelModel {
	return &customTagTargetRelModel{
		defaultTagTargetRelModel: newTagTargetRelModel(conn),
	}
}

func (m *customTagTargetRelModel) withSession(session sqlx.Session) TagTargetRelModel {
	return NewTagTargetRelModel(sqlx.NewSqlConnFromSession(session))
}
