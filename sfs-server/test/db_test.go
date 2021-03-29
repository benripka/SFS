package test

import (
	"../database"
	"testing"
)

const (
	TestFileA     = "/a/test.txt"
	TestFileB     = "/b/test.txt"
	TestDirA      = "/a/folder"
	TestUserA     = "Ben"
	TestUserB     = "Jake"
	TestPasswordA = "1234a"
	TestPasswordB = "1234b"
	TestGroupA    = "GroupA"
	TestGroupB    = "GroupB"
)

func TestAddUser(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddUser(TestUserA, TestPasswordA)
	if err != nil {
		t.Errorf("Failed to insert data with error: %s", err)
		return
	}
}

func TestAuthenticateUser(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddUser(TestUserA, TestPasswordA)
	result, err := dao.Authenticate(TestUserA, TestPasswordA)
	if err != nil || !result {
		t.Errorf("Failed to authenticate user: %s", err)
		return
	}
	result, err = dao.Authenticate("NotBen", TestPasswordA)
	if err != nil || result {
		t.Errorf("Failed to authenticate user: %s", err)
		return
	}
	result, err = dao.Authenticate(TestUserA, "badpass1234")
	if err != nil || result {
		t.Errorf("Failed to authenticate user: %s", err)
		return
	}
	result, err = dao.Authenticate("", "")
	if err != nil || result {
		t.Errorf("Failed to authenticate user: %s", err)
		return
	}
}

func TestAddGroup(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddGroup(TestGroupA)
	if err != nil {
		t.Errorf("Failed to insert data with error: %s", err)
		return
	}
}

func TestAddUserToGroup(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddGroup(TestGroupA)
	err = dao.AddUser(TestUserA, TestPasswordA)
	err = dao.AddUserToGroup(TestUserA, TestGroupA)
	if err != nil {
		t.Errorf("Failed to insert data with error: %s", err)
		return
	}
}

func TestGiveUserPermission(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddUser(TestUserA, TestPasswordA)
	err = dao.AddUserPermission(TestUserA, TestFileA)
	if err != nil {
		t.Errorf("Failed to insert data with error: %s", err)
		return
	}
}

func TestCheckUserPermissions(t *testing.T) {
	dao, err := database.NewPermissionDao()
	if err != nil {
		t.Errorf("Failed to create permissions dao: %s", err)
		return
	}
	err = dao.AddUser(TestUserA, TestPasswordA)
	err = dao.AddUser(TestUserB, TestPasswordB)
	err = dao.AddUserPermission(TestUserA, TestFileA)
	result1, err := dao.CheckUserPermission(TestUserB, TestFileA)
	if err != nil || result1 {
		t.Errorf("Failed to insert data with error: %s", err)
		return
	}
}
