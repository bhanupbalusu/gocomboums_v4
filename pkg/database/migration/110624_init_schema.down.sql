-- drop indexes
DROP INDEX IF EXISTS idx_user_role_user_id;
DROP INDEX IF EXISTS idx_user_role_role_id;
DROP INDEX IF EXISTS idx_role_permission_role_id;
DROP INDEX IF EXISTS idx_role_permission_permission_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_roles_role_name;
DROP INDEX IF EXISTS idx_permissions_permission_name;

-- assuming constraints are named as follows
ALTER TABLE UserRole DROP CONSTRAINT IF EXISTS UserRole_user_id_fkey;
ALTER TABLE UserRole DROP CONSTRAINT IF EXISTS UserRole_role_id_fkey;
ALTER TABLE RolePermission DROP CONSTRAINT IF EXISTS RolePermission_role_id_fkey;
ALTER TABLE RolePermission DROP CONSTRAINT IF EXISTS RolePermission_permission_id_fkey;

-- drop tables
DROP TABLE IF EXISTS UserRole;
DROP TABLE IF EXISTS RolePermission;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Roles;
DROP TABLE IF EXISTS "Permissions";
