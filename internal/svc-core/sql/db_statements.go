package coreSql

import "github.com/go-konsultin/sqlk"

type Statements struct {
	User           *User
	UserCredential *UserCredentialSql
	ClientAuth     *ClientAuth
	Role           *Role
}

func New(db *sqlk.DatabaseContext) *Statements {
	return &Statements{
		User:           NewUser(db),
		UserCredential: NewUserCredential(db),
		ClientAuth:     NewClientAuth(db),
		Role:           NewRole(db),
	}
}
