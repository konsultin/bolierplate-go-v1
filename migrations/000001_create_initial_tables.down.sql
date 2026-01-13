-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS "RolePrivilege";
DROP TABLE IF EXISTS "ClientAuth";
DROP TABLE IF EXISTS "Role";
DROP TABLE IF EXISTS "Privilege";
