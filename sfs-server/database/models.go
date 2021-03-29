package database

const Schema = `
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS group_memberships;
DROP TABLE IF EXISTS file_permissions;
DROP TABLE IF EXISTS check_sums;
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
    check_sum VARCHAR
);
`
