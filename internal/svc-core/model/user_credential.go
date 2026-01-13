package model

import (
	"database/sql"

	"github.com/go-konsultin/timek"
	"github.com/konsultin/project-goes-here/dto"
)

type UserCredential struct {
	Id               int64                 `db:"id"`
	UserId           int64                 `db:"userId"`
	AuthProviderId   dto.AuthProvider_Enum `db:"authProviderId"`
	CredentialKey    string                `db:"credentialKey"`    // email/phone/username for PASSWORD, provider_user_id for OAuth
	CredentialSecret sql.NullString        `db:"credentialSecret"` // password_hash for PASSWORD, null for OAuth
	IsVerified       bool                  `db:"isVerified"`
	VerifiedAt       sql.NullTime          `db:"verifiedAt"`
	CreatedAt        timek.Time            `db:"createdAt"`
	UpdatedAt        timek.Time            `db:"updatedAt"`

	User *User `db:"-"`
}

func NewUserCredential() *UserCredential {
	return &UserCredential{}
}
