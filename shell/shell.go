package main

import (
	"../sfs-client"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Signup(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 3 {
		return "Error: not enough arguments.\nProper usage: signup <username> <password>"
	} else {
		output, err := client.Signup(tokens[1], tokens[2])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Login(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 3 {
		return "Error: not enough arguments.\nProper usage: login <username> <password>"
	} else {
		output, err := client.Login(tokens[1], tokens[2])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Logout(tokens []string, client *sfs_client.Client) string {
	client.SignedIn = false
	if len(tokens) != 1 {
		return "Error: to many arguments.\nProper usage: logout"
	} else {
		output, err := client.Logout()
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Ls(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 1 {
		return "Error: to many arguments.\nProper usage: ls"
	} else {
		output, err := client.Ls()
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Pwd(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 1 {
		return "\nError: to many arguments.\nProper usage: pwd"
	} else {
		output, err := client.Pwd()
		if err != nil {
			return "\nError: something went wrong."
		}
		return output
	}
}

func Cat(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: cat <filename>"
	} else {
		output, err := client.Cat(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Touch(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: touch <filename>"
	} else {
		output, err := client.Touch(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Mkdir(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: mkdir <filename>"
	} else {
		output, err := client.Mkdir(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Cd(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: touch <filename>"
	} else {
		output, err := client.Cd(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Mv(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 3 {
		return "Error: wrong number of arguments.\nProper usage: mv <filename> <new path>"
	} else {
		output, err := client.Mv(tokens[1], tokens[2])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Rm(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: mkdir <filename>"
	} else {
		output, err := client.Rm(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func Write(tokens []string, client *sfs_client.Client) string {
	if len(tokens) < 3 {
		return "Error: wrong number of arguments.\nProper usage: example.txt < 'data to write to file'"
	} else {
		data := strings.Join(tokens[2:], " ")
		output, err := client.Write(tokens[1], data)
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func AddGroup(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 2 {
		return "Error: wrong number of arguments.\nProper usage: mkdir <filename>"
	} else {
		output, err := client.AddGroup(tokens[1])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func AddUserToGroup(tokens []string, client *sfs_client.Client) string {
	if len(tokens) != 3 {
		return "Error: wrong number of arguments.\nProper usage: mkdir <filename>"
	} else {
		output, err := client.AddUserToGroup(tokens[1], tokens[2])
		if err != nil {
			return "Error: something went wrong."
		}
		return output
	}
}

func help() string {
	return "Here are all the commands you'll need:\n\n" +
		"signup <username> <password> \t\t Create a new account\n" +
		"login <username> <password> \t\t Log in to a user account\n" +
		"ls \t\t\t\t\t\t\t\t\t List the contents of the current directory\n" +
		"pwd \t\t\t\t\t\t\t\t show the current directory path\n" +
		"mkdir <directory_name> \t\t\t\t Create a new directory in current directory\n" +
		"cd \t\t\t\t\t\t\t\t\t Change the current directory\n" +
		"cat <file_name> \t\t\t\t\t Show contents of file, line by line.\n" +
		"touch <file_name> \t\t\t\t\t create a new file with provided name in current directory\n" +
		"mv <old_path> <new_path> \t\t\t move a file from one location to another\n" +
		"addgroup <groupname> \t\t\t\t Create a new group with given name\n" +
		"addtogroup <username> <groupname> \t Add a new user to group with provided name\n"
}

func handleInput(tokens []string, client *sfs_client.Client) string {
	switch tokens[0] {
	case "login":
		return Login(tokens, client)
	case "signup":
		return Signup(tokens, client)
	case "logout":
		return Logout(tokens, client)
	case "ls":
		return Ls(tokens, client)
	case "pwd":
		return Pwd(tokens, client)
	case "cat":
		return Cat(tokens, client)
	case "mkdir":
		return Mkdir(tokens, client)
	case "touch":
		return Touch(tokens, client)
	case "cd":
		return Cd(tokens, client)
	case "mv":
		return Mv(tokens, client)
	case "rm":
		return Rm(tokens, client)
	case "write":
		return Write(tokens, client)
	case "addgroup":
		return AddGroup(tokens, client)
	case "addtogroup":
		return AddUserToGroup(tokens, client)
	case "help":
		return help()
	default:
		return errorMessage()
	}
}

func loginMessage() string {
	return "Not logged in.\nPlease log in using the login command:\t login <username> <password>"
}

func introMessage() string {
	art := " ___________ _____  \n/  ___|  ___/  ___| \n\\ `--.| |_  \\ `--.  \n `--. \\  _|  `--. \\ \n/\\__/ / |   /\\__/ / \n\\____/\\_|   \\____/  \n                    \n                    "
	return art + "\nSignup or login to get started. If you're lost, you can always ask for: help"
}

func errorMessage() string {
	return "Command not recognized\n"
}

func main() {
	client := sfs_client.NewClient()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(introMessage())
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		preparedInput := strings.ReplaceAll(input, "\n", "")
		tokens := strings.Split(preparedInput, " ")
		if len(tokens) == 0 {
			continue
		}
		if client.SignedIn || tokens[0] == "login" || tokens[0] == "signup" || tokens[0] == "help" {
			result := handleInput(tokens, client)
			fmt.Println(result)
		} else {
			fmt.Println(loginMessage())
			continue
		}

	}
}
