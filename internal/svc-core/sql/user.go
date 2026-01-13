package coreSql

import (
	"github.com/go-konsultin/sqlk"
	"github.com/go-konsultin/sqlk/pq/query"
	"github.com/jmoiron/sqlx"
)

type User struct {
	GetUserByXid     *sqlx.Stmt
	GetUserById      *sqlx.Stmt
	FindByIdentifier *sqlx.Stmt
	Insert           *sqlx.NamedStmt
}

func NewUser(db *sqlk.DatabaseContext) *User {
	return &User{
		GetUserByXid: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserSchema).
				Where(
					query.Equal(query.Column("xid")),
				).Build(),
		),
		GetUserById: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserSchema).
				Where(
					query.Equal(query.Column("id")),
				).Build(),
		),
		FindByIdentifier: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserSchema).
				Where(
					query.Or(
						query.Equal(query.Column("email")),
						query.Equal(query.Column("phone")),
						query.Equal(query.Column("username")),
					),
				).
				Limit(1).Build(),
		),
		Insert: db.MustPrepareNamed(
			query.Insert(UserSchema,
				"xid",
				"username",
				"fullName",
				"email",
				"phone",
				"age",
				"avatar",
				"statusId",
				"createdAt",
				"updatedAt",
				"modifiedBy",
				"version",
				"metadata",
			).Build(),
		),
	}
}
