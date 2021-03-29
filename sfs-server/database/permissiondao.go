package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"sync"
)

const (
	DriverName = "sqlite3"
	DbName     = "/home/ubuntu/ECE_422_Project_1/db/sfs.db"
)

var Dao *PermissionDao

type PermissionDao struct {
	lock sync.Mutex
	db   *sqlx.DB
}

func NewPermissionDao() (*PermissionDao, error) {
	db, err := sqlx.Connect(DriverName, DbName)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(Schema)
	if err != nil {
		return nil, err
	}
	db.MustBegin()
	return &PermissionDao{
		db: db,
	}, nil
}

func (dao *PermissionDao) CheckUserExists(username string) (bool, error) {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	rows, err := dao.db.Query(CheckUserExistsQuery, username)
	if err != nil || rows == nil {
		return false, err
	}
	var exists bool
	if !rows.Next() {
		return false, errors.New("Failed to check if user exists.")
	}
	err = rows.Scan(&exists)
	if err != nil {
		return false, err
	}
	if err = rows.Close(); err != nil {
		return false, err
	}
	return exists, nil
}

func (dao *PermissionDao) AddUser(username string, password string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(AddUserQuery, username, password)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) Authenticate(username string, password string) (bool, error) {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	hasPermission := new(bool)
	err := dao.db.Get(hasPermission, AuthenticateUserQuery, username, password)
	if err != nil {
		return false, err
	}
	return *hasPermission, nil
}

func (dao *PermissionDao) AddGroup(groupName string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(AddGroupQuery, groupName)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) AddUserToGroup(username string, groupName string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(AddUserToGroupQuery, username, groupName)
	if err != nil {
		return err
	}
	_, err = tx.Exec(UpdateGroupPermissions, username)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) AddUserPermission(username string, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err = tx.Exec(AddUserPermissionsQuery, absPath, username)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) UpdatePermissionForAllUsersGroups(username string, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err = tx.Exec(AddPermissionForAllUsersGroups, absPath, username)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) AddGroupPermission(groupName string, path string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(AddGroupPermissionsQuery, path, groupName)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) AddCheckSum(path string, checkSum string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(AddCheckSum, path, checkSum)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) UpdateCheckSum(path string, checkSum string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	tx := dao.db.MustBegin()
	_, err := tx.Exec(UpdateCheckSum, checkSum, path)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dao *PermissionDao) CheckUserPermission(username string, path string) (bool, error) {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	rows, err := dao.db.Query(CheckUserHasPermissionQuery, username, path)
	if err != nil || rows == nil {
		return false, err
	}
	var permission bool
	if !rows.Next() {
		return false, errors.New("Failed to get permissions from the db")
	}
	err = rows.Scan(&permission)
	if err != nil {
		return false, err
	}
	if err = rows.Close(); err != nil {
		return false, err
	}
	return permission, nil
}

func (dao *PermissionDao) CheckUsersGroupPermission(username string, path string) (bool, error) {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	rows, err := dao.db.Query(CheckUserGroupsPermissionQuery, username, path)
	if err != nil || rows == nil {
		return false, err
	}
	var permission bool
	if !rows.Next() {
		return false, errors.New("Failed to get permissions from the db")
	}
	err = rows.Scan(&permission)
	if err != nil {
		return false, err
	}
	if err = rows.Close(); err != nil {
		return false, err
	}
	return permission, nil
}

func (dao *PermissionDao) GetCheckSum(path string) (string, error) {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	rows, err := dao.db.Query(GetCheckSum, path)
	if err != nil || rows == nil {
		return "", err
	}
	var result string
	if !rows.Next() {
		return "", errors.New("Failed to get permissions from the db")
	}
	err = rows.Scan(&result)
	if err != nil {
		return "", err
	}
	if err = rows.Close(); err != nil {
		return "", err
	}
	return result, nil
}

func (dao *PermissionDao) ChangeFilePath(oldPath string, newPath string) error {
	dao.lock.Lock()
	defer dao.lock.Unlock()
	_, err := dao.db.Exec(ChangeFilePathPermission, newPath, oldPath)
	if err != nil {
		return err
	}
	_, err = dao.db.Exec(ChangeFilePathCheckSums, newPath, oldPath)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	var err error
	Dao, err = NewPermissionDao()
	if err != nil {
		panic("Failed to create permission dao")
	}
}
