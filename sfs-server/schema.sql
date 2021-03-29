DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS group_memberships;
DROP TABLE IF EXISTS file_permissions;
CREATE TABLE users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    password VARCHAR NOT NULL,
    username VARCHAR NOT NULL
);

CREATE TABLE groups
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    group_name VARCHAR
);

CREATE TABLE group_memberships
(
    user_id  INT NOT NULL,
    group_id INT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (group_id) REFERENCES groups (id)
);

CREATE TABLE file_permissions
(
    file_path VARCHAR NOT NULL,
    user_id   INT,
    group_id  INT,
    read      BOOLEAN,
    write     BOOLEAN,
    PRIMARY KEY (file_path, user_id, group_id),
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE check_sums
(
    file_path VARCHAR PRIMARY KEY NOT NULL,
    check_sum BLOB
);

INSERT INTO users (password, username)
VALUES ('123454', 'Jake'),
       ('4321', 'Jo'),
       ('4321', 'Ben');
INSERT INTO groups (group_name)
values ('group1');

-- Add user to group (username, group_name)
INSERT OR
REPLACE
INTO group_memberships (user_id, group_id)
select u.id, g.id
from users u,
     groups g
where u.username = 'Jo'
  and g.group_name = 'group1';

-- Check if user in group (username, group_name)
SELECT CASE
           WHEN EXISTS(
                   SELECT *
                   FROM group_memberships
                   WHERE group_id = (select id from groups where group_name = 'group1')
                     and user_id = (select id from users where username = 'Jo')
               )
               THEN 'TRUE'
           ELSE 'FALSE'
           END;

-- Add user permissions
INSERT OR
REPLACE
INTO file_permissions (file_path, user_id, read, write)
select '/c/home/ben/file.txt', u.id, TRUE, TRUE
from users u
where u.username = 'Jake';

-- Add group permissions
INSERT OR
REPLACE
INTO file_permissions (file_path, group_id, read, write)
select '/c/home/ben/file.txt', g.id, TRUE, TRUE
from groups g
where g.group_name = 'group1';

-- Check user has access (path, username)
SELECT CASE
           WHEN
                   EXISTS(
                           select u.id
                           from users u
                                    join file_permissions fp on u.id = fp.user_id
                           where u.username = 'Ben'
                             and fp.file_path = '/c/home/ben/file.txt'
                       )
                   OR
                   EXISTS(
                           select u.id
                           from users u
                                    join group_memberships gm on u.id = gm.user_id
                                    join groups g on g.id = gm.group_id
                                    join file_permissions fp on g.id = fp.group_id or u.id = fp.user_id
                           where u.username = 'Ben'
                             and fp.file_path = '/c/home/ben/file.txt'
                       )
               THEN 'TRUE'
           ELSE 'FALSE'
           END;

-- Athenticate user (username, password)
SELECT CASE
           WHEN
               EXISTS(
                       select id
                       from users
                       where username = ''
                         and password = ''
                   )
               THEN 'TRUE'
           ELSE 'FALSE'
           END;

-- Add permission to path for all groups the users is part of
INSERT OR
REPLACE
INTO file_permissions (file_path, group_id, read, write)
select '<path>', g.id, TRUE, TRUE
from groups g
where g.id in (
    SELECT group_id
    from group_memberships gm
    WHERE gm.user_id = '<user_id>'
);

-- Change file path
UPDATE file_permissions
SET file_path = ''
WHERE file_path = '';

select *
from file_permissions;
SELECT *
FROM groups;

