package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BountyApplicantModel = (*customBountyApplicantModel)(nil)

type (
	// BountyApplicantModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBountyApplicantModel.
	BountyApplicantModel interface {
		bountyApplicantModel
		withSession(session sqlx.Session) BountyApplicantModel
	}

	customBountyApplicantModel struct {
		*defaultBountyApplicantModel
	}
)

// NewBountyApplicantModel returns a model for the database table.
func NewBountyApplicantModel(conn sqlx.SqlConn) BountyApplicantModel {
	return &customBountyApplicantModel{
		defaultBountyApplicantModel: newBountyApplicantModel(conn),
	}
}

func (m *customBountyApplicantModel) withSession(session sqlx.Session) BountyApplicantModel {
	return NewBountyApplicantModel(sqlx.NewSqlConnFromSession(session))
}
