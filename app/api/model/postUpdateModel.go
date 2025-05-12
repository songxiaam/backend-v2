package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PostUpdateModel = (*customPostUpdateModel)(nil)

type (
	// PostUpdateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPostUpdateModel.
	PostUpdateModel interface {
		postUpdateModel
		withSession(session sqlx.Session) PostUpdateModel
	}

	customPostUpdateModel struct {
		*defaultPostUpdateModel
	}
)

// NewPostUpdateModel returns a model for the database table.
func NewPostUpdateModel(conn sqlx.SqlConn) PostUpdateModel {
	return &customPostUpdateModel{
		defaultPostUpdateModel: newPostUpdateModel(conn),
	}
}

func (m *customPostUpdateModel) withSession(session sqlx.Session) PostUpdateModel {
	return NewPostUpdateModel(sqlx.NewSqlConnFromSession(session))
}
