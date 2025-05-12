package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DictDataModel = (*customDictDataModel)(nil)

type (
	// DictDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDictDataModel.
	DictDataModel interface {
		dictDataModel
		withSession(session sqlx.Session) DictDataModel
	}

	customDictDataModel struct {
		*defaultDictDataModel
	}
)

// NewDictDataModel returns a model for the database table.
func NewDictDataModel(conn sqlx.SqlConn) DictDataModel {
	return &customDictDataModel{
		defaultDictDataModel: newDictDataModel(conn),
	}
}

func (m *customDictDataModel) withSession(session sqlx.Session) DictDataModel {
	return NewDictDataModel(sqlx.NewSqlConnFromSession(session))
}
