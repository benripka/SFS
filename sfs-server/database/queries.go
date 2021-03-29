package database

const AddUserQuery = `INSERT INTO users (username, password) values ($1, $2)`

const AddGroupQuery = `INSERT INTO groups (group_name) values ($1)`

const AddUserToGroupQuery = `
	INSERT OR REPLACE INTO group_memberships (user_id, group_id)
	select u.id, g.id
	from users u,
		 groups g
	where u.username = $1 and g.group_name = $2;
`

const AuthenticateUserQuery = `
SELECT CASE
           WHEN
               EXISTS(
                       select id
                       from users
                       where username = $1 and password = $2
                   )
               THEN 'TRUE'
           ELSE 'FALSE'
        END;
`

const AddUserPermissionsQuery = `
INSERT
INTO file_permissions
select ?, id, null, TRUE, TRUE
from users
where username = ?;
`

const AddPermissionForAllUsersGroups = `
INSERT OR
REPLACE
INTO file_permissions (file_path, group_id, read, write)
select ?, g.id, TRUE, TRUE
from groups g
where g.id in (
    SELECT group_id
    from group_memberships gm
    join users u on gm.user_id = u.id
    WHERE u.username = ?
);
`

const UpdateGroupPermissions = `
INSERT OR REPLACE
INTO file_permissions (file_path, user_id, group_id, read, write)
SELECT fp.file_path, fp.user_id, gm.group_id, fp.read, fp.write
FROM file_permissions fp
JOIN users u on fp.user_id = u.id
JOIN group_memberships gm on u.id = gm.user_id
WHERE u.username = ?;
`

const AddGroupPermissionsQuery = `
	INSERT OR
	REPLACE
	INTO file_permissions (file_path, group_id, read, write)
	select ?, g.id, TRUE, TRUE
	from groups g
	where g.group_name = ?;
`

const AddCheckSum = `
	INSERT INTO check_sums (file_path, check_sum) VALUES (?, ?)
`

const GetCheckSum = `
	select check_sum
	from check_sums
	where file_path = ?
`

const UpdateCheckSum = `
	UPDATE check_sums
	SET check_sum = ?
	WHERE file_path = ?
`

const CheckUserHasPermissionQuery = `
SELECT CASE
           WHEN
                   EXISTS(
                           select u.id
                           from users u
                                    join file_permissions fp on u.id = fp.user_id
                           where u.username = ?
                             and fp.file_path = ?
                       )
               THEN 'TRUE'
           ELSE 'FALSE'
           END;
`

const CheckUserGroupsPermissionQuery = `
SELECT CASE
           WHEN
                   EXISTS(
                           select u.id
                           from users u
                                    join group_memberships gm on u.id = gm.user_id
                                    join groups g on g.id = gm.group_id
                                    join file_permissions fp on g.id = fp.group_id
                           where u.username = ?
                             and fp.file_path = ?
                       )
			THEN 'TRUE'
           ELSE 'FALSE'
           END;
`

const CheckUserExistsQuery = `
SELECT CASE
           WHEN
                   EXISTS(
                           select u.id
                           from users u
                           where u.username = ?
                       )
			THEN 'TRUE'
           ELSE 'FALSE'
           END;
`

const ChangeFilePathPermission = `
UPDATE file_permissions
SET file_path = ?
WHERE file_path = ?;
`

const ChangeFilePathCheckSums = `
UPDATE check_sums
SET file_path = ?
WHERE file_path = ?;
`
