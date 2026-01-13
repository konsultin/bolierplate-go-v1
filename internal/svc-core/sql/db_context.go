package coreSql

import (
	"github.com/go-konsultin/sqlk/schema"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
)

// User Schemas
var (
	UserSchema           = schema.New(schema.FromModelRef(new(model.User)), schema.As("User"))
	UserCredentialSchema = schema.New(schema.FromModelRef(new(model.UserCredential)), schema.As("UserCredential"))
	ClientAuthSchema     = schema.New(schema.FromModelRef(new(model.ClientAuth)), schema.As("ClientAuth"))
	RoleSchema           = schema.New(schema.FromModelRef(new(model.Role)), schema.As("Role"))
	RolePrivilegeSchema  = schema.New(schema.FromModelRef(new(model.RolePrivilege)), schema.As("RolePrivilege"))
	PrivilegeSchema      = schema.New(schema.FromModelRef(new(model.Privilege)), schema.As("Privilege"))
)
