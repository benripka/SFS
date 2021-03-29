package test

import (
	"../database"
	"../encryption"
	"testing"
)

func TestEncryption(t *testing.T) {
	TestString1 := "golangfam"
	err := encryption.DecryptMany(&TestString1)
	err = encryption.DecryptMany(&TestString1)
	if err != nil || TestString1 != "golangfam" {
		t.Errorf("failed to encrypt string with error: %s", err)
	}
}

func TestEncrypMany(t *testing.T) {
	TestString1 := "golangfam"
	TestString2 := "golangfamily"
	err := encryption.EncryptMany(&TestString1, &TestString2)
	_ = database.Dao.AddUser(TestString1, TestString2)
	err = encryption.DecryptMany(&TestString1, &TestString2)
	err = encryption.EncryptMany(&TestString1, &TestString2)
	res, _ := database.Dao.Authenticate(TestString1, TestString2)
	if err != nil || !res {
		t.Errorf("failed to encrypt string with error: %s", err)
	}
}

func TestEncryptConsistent(t *testing.T) {
	TestString1 := "golangfam"
	TestString2 := "golangfam"
	err := encryption.EncryptMany(&TestString1)
	err = encryption.EncryptMany(&TestString2)
	err = encryption.DecryptMany(&TestString1)
	err = encryption.DecryptMany(&TestString2)
	if err != nil || TestString2 != TestString1 {
		t.Errorf("failed to encrypt string with error: %s", err)
	}
}

func TestChecksum(t *testing.T) {
	filename := "./testfile.txt"
	checkSum, err := encryption.CheckSum(filename)
	print(checkSum)
	if err != nil {
		t.Errorf("failed")
	}
}
