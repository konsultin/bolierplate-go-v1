package coreSql

import (
	"github.com/go-konsultin/sqlk"
	"github.com/go-konsultin/sqlk/pq/query"
	"github.com/jmoiron/sqlx"
)

type UserCredentialSql struct {
	FindByProviderAndKey *sqlx.Stmt
	FindByUserId         *sqlx.Stmt
	Insert               *sqlx.NamedStmt
	UpdateSecret         *sqlx.Stmt
}

func NewUserCredential(db *sqlk.DatabaseContext) *UserCredentialSql {
	return &UserCredentialSql{
		FindByProviderAndKey: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserCredentialSchema).
				Where(
					query.And(
						query.Equal(query.Column("auth_provider_id")),
						query.Equal(query.Column("credential_key")),
					),
				).Build(),
		),
		FindByUserId: db.MustPrepareRebind(
			query.Select(
				query.Column("*"),
			).
				From(UserCredentialSchema).
				Where(
					query.Equal(query.Column("user_id")),
				).Build(),
		),
		Insert: db.MustPrepareNamed(
			query.Insert(UserCredentialSchema,
				"userId",
				"authProviderId",
				"credentialKey",
				"credentialSecret",
				"isVerified",
				"verifiedAt",
				"createdAt",
				"updatedAt",
			).Build(),
		),
		UpdateSecret: db.MustPrepareRebind(`
			UPDATE "UserCredential"
			SET "credentialSecret" = ?, "updatedAt" = NOW()
			WHERE "id" = ?
		`),
	}
}
