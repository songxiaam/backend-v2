package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ImageModel = (*customImageModel)(nil)

type (
	// ImageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImageModel.
	ImageModel interface {
		imageModel
		withSession(session sqlx.Session) ImageModel
	}

	customImageModel struct {
		*defaultImageModel
	}
)

// NewImageModel returns a model for the database table.
func NewImageModel(conn sqlx.SqlConn) ImageModel {
	return &customImageModel{
		defaultImageModel: newImageModel(conn),
	}
}

func (m *customImageModel) withSession(session sqlx.Session) ImageModel {
	return NewImageModel(sqlx.NewSqlConnFromSession(session))
}
