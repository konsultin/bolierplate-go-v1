package model

import (
	"database/sql"

	"github.com/konsultin/project-goes-here/dto"
)

type User struct {
	BaseField
	Id       int64                  `db:"id"`
	Xid      string                 `db:"xid"`
	Username sql.NullString         `db:"username"`
	FullName string                 `db:"fullName"`
	Phone    sql.NullString         `db:"phone"`
	Email    sql.NullString         `db:"email"`
	Age      sql.NullString         `db:"age"`
	Avatar   sql.NullString         `db:"avatar"`
	StatusId dto.ControlStatus_Enum `db:"statusId"`
}

func NewUser() *User {
	return &User{}
}
