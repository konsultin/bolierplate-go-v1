package coreSql

import (
	"github.com/go-konsultin/sqlk/schema"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
)

// User Schemas
var (
	UserSchema           = schema.New(schema.FromModelRef(new(model.User)), schema.As("user"))
	UserCredentialSchema = schema.New(schema.FromModelRef(new(model.UserCredential)), schema.As("userCredential"))
	ClientAuthSchema     = schema.New(schema.FromModelRef(new(model.ClientAuth)), schema.As("clientAuth"))
	RoleSchema           = schema.New(schema.FromModelRef(new(model.Role)), schema.As("role"))
	RolePrivilegeSchema  = schema.New(schema.FromModelRef(new(model.RolePrivilege)), schema.As("rolePrivilege"))
	PrivilegeSchema      = schema.New(schema.FromModelRef(new(model.Privilege)), schema.As("privilege"))
)
