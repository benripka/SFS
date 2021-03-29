package fs

import (
	"../database"
	"../encryption"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const HomeDir = "/home/ubuntu/ECE_422_Project_1/home/"

func Authenticate(username string, password string) (bool, error) {
	if err := encryption.EncryptMany(&username, &password); err != nil {
		return false, err
	}
	loggedIn, err := database.Dao.Authenticate(username, password)
	return loggedIn, err
}

func AddUser(username string, password string) error {
	if err := encryption.EncryptMany(&username, &password); err != nil {
		return err
	}
	exists, err := database.Dao.CheckUserExists(username)
	if err != nil || exists {
		return err
	}
	err = database.Dao.AddUser(username, password)
	if err != nil {
		return err
	}
	homeDir, err := GetHomeDir(username)
	if err != nil {
		return err
	}
	err = database.Dao.AddUserPermission(username, filepath.Dir(homeDir))
	if err != nil {
		return err
	}
	return nil
}

func AddGroup(name string) error {
	if err := encryption.EncryptMany(&name); err != nil {
		return err
	}
	return database.Dao.AddGroup(name)
}

func AddUserToGroup(username string, groupname string) error {
	if err := encryption.EncryptMany(&username, &groupname); err != nil {
		return err
	}
	err := database.Dao.AddUserToGroup(username, groupname)
	return err
}

func ValidateCheckSums(username string) (string, error) {
	homeDir, err := GetHomeDir(username)
	if err != nil {
		return "", err
	}
	files := make([]string, 0)
	err = filepath.Walk(homeDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return "", err
	}
	tamperedFiles := make([]string, 0)
	for _, path := range files {
		checkSum, err := encryption.CheckSum(path)
		if err != nil {
			return "Tamper check failed.", err
		}
		sum, err := database.Dao.GetCheckSum(path)
		if sum != string(checkSum) {
			tamperedFiles = append(tamperedFiles, path)
		}
	}
	if len(tamperedFiles) == 0 {
		return "All files are un-tampered-with :)", nil
	} else {
		return "TAMPER ALERT!!!\nThe following files have been tampered with:\n" + strings.Join(tamperedFiles, "\n"), nil
	}
}

func Ls(workingDir string, username string) (string, error) {
	if err := encryption.EncryptMany(&username); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}
	result := make([]string, 0)
	for _, f := range files {
		path := f.Name()
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		hasPermission, err := database.Dao.CheckUserPermission(username, absPath)
		if hasPermission {
			if err := encryption.DecryptMany(&path); err != nil {
				return "", nil
			}
		}
		result = append(result, path)
	}
	return strings.Join(result, "\n"), nil
}

// mkdir <directory_name> - Create a new directory in current directory
// function to get functionality
func Mkdir(workingDir string, username string, directoryName string) (string, error) {
	if err := encryption.EncryptMany(&username, &directoryName); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	err := os.Mkdir(directoryName, os.ModeDir)
	// looks like file exists
	if os.IsExist(err) {
		return "", err
	}
	absPath, err := filepath.Abs(directoryName)
	if err != nil {
		return "", err
	}
	err = database.Dao.AddUserPermission(username, absPath)
	if err != nil {
		return "", err
	}
	err = database.Dao.UpdatePermissionForAllUsersGroups(username, absPath)
	if err != nil {
		return "", err
	}
	return "Folder created", nil
}

// cd - Change the current directory (support ~/./..)
// function to get functionality
// this changes the 'virtual' directory - the shell directory stays the same
func Cd(workingDir string, username string, newDir string) (string, error) {
	if strings.Contains(newDir, "~") {
		rest := strings.Replace(newDir, "~", "", 1)
		if rest != "" {
			err := encryption.EncryptMany(&rest)
			if err != nil {
				return "", nil
			}
		}
		homeDir, err := GetHomeDir(username)
		if err != nil {
			return "", err
		}
		newDir = homeDir + rest
	} else {
		err := encryption.EncryptMany(&newDir)
		if err != nil {
			return "", err
		}
	}
	if err := encryption.EncryptMany(&username); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(newDir)
	if err != nil {
		return "", err
	}
	groupPermission, err := database.Dao.CheckUsersGroupPermission(username, absPath)
	if err != nil {
		return "", err
	}
	userPermission, err := database.Dao.CheckUserPermission(username, absPath)
	if err != nil {
		return "", err
	}
	if !groupPermission && !userPermission {
		return "", errors.New("Permission denied.")
	}
	if err := os.Chdir(absPath); err != nil {
		return "", err
	}
	newCurDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return newCurDir, nil
}

// cat <file_name> - Show contents of file, line by line.
// example: "cat file1"
func Cat(workingDir string, username string, filePath string) (string, error) {
	if err := encryption.EncryptMany(&username, &filePath); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	permission, err := database.Dao.CheckUserPermission(username, absPath)
	if err != nil {
		return "", err
	}
	if !permission {
		return "You are not authorized to access this file.", nil
	}
	// https://golangcode.com/read-a-files-contents/
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", nil
	}
	// Convert []byte to string and print to screen
	text := string(content)
	return text, nil
}

// touch <file_name> - create a new file with provided name in current directory
// example: "touch file1"
func Touch(workingDir string, username string, filename string) (string, error) {
	if err := encryption.EncryptMany(&username, &filename); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	permission, err := database.Dao.CheckUserPermission(username, filepath.Dir(absPath))
	if err != nil {
		return "", err
	}
	if !permission {
		return "You do not have authorization to create a file in this location.", nil
	}
	if !pathExists(absPath) {
		if err := database.Dao.AddUserPermission(username, absPath); err != nil {
			return "", err
		}
	}
	_, err = os.OpenFile(absPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	err = database.Dao.AddCheckSum(absPath, "")
	if err != nil {
		return "", err
	}
	return "Done.", nil
}

// mv <old_path> <new_path> - move a file from one location to another
// example: "mv /home/folder1/file1 /home/folder1/folder2/file1
func Mv(workingDir string, username string, oldPath string, newPath string) (string, error) {
	if err := encryption.EncryptMany(&username, &oldPath, &newPath); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", nil
	}
	oldPath, err := filepath.Abs(oldPath)
	if err != nil {
		return "", err
	}
	newPath, err = filepath.Abs(newPath)
	if err != nil {
		return "", err
	}
	oldPathPermission, err := database.Dao.CheckUserPermission(username, oldPath)
	if err != nil {
		return "", err
	}
	newPathPermission, err := database.Dao.CheckUserPermission(username, filepath.Dir(newPath))
	if err != nil {
		return "", err
	}
	if !newPathPermission {
		return "You are not authorized to write to this location", nil
	}
	if !oldPathPermission {
		return "You are not authorized to move this object", nil
	}
	err = database.Dao.ChangeFilePath(oldPath, newPath)
	if err != nil {
		return "", err
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return "", err
	}
	return "Done.", nil
}

// rm <file_name> - Delete file
func Rm(workingDir string, username string, path string) error {
	if err := encryption.EncryptMany(&username, &path); err != nil {
		return err
	}
	if err := os.Chdir(workingDir); err != nil {
		return err
	}
	// Removing file from the directory
	absPath, err := filepath.Abs(path)
	permission, err := database.Dao.CheckUserPermission(username, absPath)
	if err != nil {
		return err
	}
	if !permission {
		return errors.New("No permission")
	}
	if err != nil {
		return err
	}
	err = os.RemoveAll(absPath)
	if err != nil {
		return err
	}
	return nil
}

func Write(workingDir string, username string, filename string, data []byte) (string, error) {
	if err := encryption.EncryptMany(&username, &filename); err != nil {
		return "", err
	}
	if err := os.Chdir(workingDir); err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(filename)
	if pathExists(absPath) {
		if err != nil {
			return "", err
		}
		permission, err := database.Dao.CheckUserPermission(username, absPath)
		if err != nil {
			return "", err
		}
		if !permission {
			return "You are not authorized to write to this file", nil
		}
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return "", err
		}
		_, err = f.Write(data)
		if err != nil {
			return "", err
		}
		f.Close()
		checksum, err := encryption.CheckSum(filename)
		err = database.Dao.UpdateCheckSum(absPath, string(checksum))
		if err != nil {
			return "", nil
		}
		return "Done.", nil
	}
	return "File does not exist.", nil
}

func Pwd(workingDir string, username string) (string, error) {
	if err := encryption.EncryptMany(&username); err != nil {
		return "", err
	}
	tokens := strings.Split(workingDir, "/")
	for i, _ := range tokens {
		permission, err := database.Dao.CheckUserPermission(username, strings.Join(tokens[:i], "/"))
		if err != nil {
			return "", err
		}
		if permission && tokens[i] != "home" {
			err := encryption.DecryptMany(&tokens[i])
			if err != nil {
				return "", err
			}
		}
	}
	return strings.Join(tokens, "/"), nil
}

func GetHomeDir(username string) (string, error) {
	if err := encryption.EncryptMany(&username); err != nil {
		return "", err
	}
	path := HomeDir + username
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModeDir); err != nil {
			return "", err
		}
	}
	if err := database.Dao.AddUserPermission(username, path); err != nil {
		return "", err
	}
	err := database.Dao.UpdatePermissionForAllUsersGroups(username, path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
