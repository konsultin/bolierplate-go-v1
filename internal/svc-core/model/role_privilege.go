package model

import "github.com/konsultin/project-goes-here/dto"

type RolePrivilege struct {
	BaseField
	Id          int64 `db:"id"`
	RoleId      int32 `db:"roleId"`
	PrivilegeId int64 `db:"privilegeId"`
}

type RolePrivilegeJoinRow struct {
	RolePrivilege *RolePrivilege `db:"RolePrivilege"`
	Role          *Role          `db:"Role"`
	Privilege     *Privilege     `db:"Privilege"`
}

func NewRolePrivilege(privilegeId int64, s *dto.Subject) *RolePrivilege {
	return &RolePrivilege{
		BaseField:   NewBaseField(s),
		PrivilegeId: privilegeId,
	}
}
