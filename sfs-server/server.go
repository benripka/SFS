package main

import (
	"./fs"
	"./session"
	_ "./session/providers/memory"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	FilePathParam  = "filepath"
	NewPathParam   = "newpath"
	UserNameParam  = "username"
	GroupNameParam = "groupname"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	valid, err := fs.Authenticate(creds.Username, creds.Password)
	if err != nil || !valid {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	sess := session.SessionManager.SessionStart(w, r)
	if err := sess.Set(session.Username, creds.Username); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	homeDir, err := fs.GetHomeDir(creds.Username)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	if err := sess.Set(session.WorkingDir, homeDir); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	output, err := fs.ValidateCheckSums(creds.Username)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := fmt.Fprint(w, "Logged in :)\n"+output); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := session.SessionManager.SessionEnd(w, r)
	var message = new(bytes.Buffer)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		message.WriteString("The server failed to log you out.")
	} else {
		message.WriteString("Logged out :)")
	}
	w.Write(message.Bytes())
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	err = fs.AddUser(creds.Username, creds.Password)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		_, _ = w.Write([]byte("Login failed. Try a different username."))
		return
	}
	sess := session.SessionManager.SessionStart(w, r)
	if err := sess.Set(session.Username, creds.Username); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
	}
	homeDir, err := fs.GetHomeDir(creds.Username)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
		return
	}
	if err := sess.Set(session.WorkingDir, homeDir); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Login failed. please try again."))
	}
	if _, err := fmt.Fprint(w, fmt.Sprintf("Logged in. Welcome to SFS %s :)", creds.Username)); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		fmt.Println("Failed to start session ...")
	}
}

func lsHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	username, workingDir := getSessionInfo(w, r)
	output, err := fs.Ls(workingDir, username)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Unable to read contents of dir"
	}
	w.Write([]byte(output))
}

func pwdHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	username, workingDir := getSessionInfo(w, r)
	output, err := fs.Pwd(workingDir, username)
	if err != nil {
		output = "Failed to print working directory"
	}
	w.Write([]byte(output))
}

func mkdirHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	username, workingDir := getSessionInfo(w, r)
	path := r.URL.Query().Get(FilePathParam)
	output, err := fs.Mkdir(workingDir, username, path)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Failed to make directory"))
	}
	w.Write([]byte(output))
}

func cdHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	username, workingDir := getSessionInfo(w, r)
	path := r.URL.Query().Get(FilePathParam)
	newDir, err := fs.Cd(workingDir, username, path)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Failed to change directory"))
		return
	}
	err = changeWorkingDir(w, r, newDir)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Failed to change directory"))
	} else {
		w.Write([]byte("Done."))
	}
}

func catHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	username, workingDir := getSessionInfo(w, r)
	path := r.URL.Query().Get(FilePathParam)
	output, err := fs.Cat(workingDir, username, path)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Command failed, try again."))
	}
	w.Write([]byte(output))
}

func touchHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	fileName := r.URL.Query().Get(FilePathParam)
	username, workingDir := getSessionInfo(w, r)
	output, err := fs.Touch(workingDir, username, fileName)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Failed to create file " + fileName
	}
	w.Write([]byte(output))
}

func mvHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	path := r.URL.Query().Get(FilePathParam)
	newPath := r.URL.Query().Get(NewPathParam)
	username, workingDir := getSessionInfo(w, r)
	var output string
	output, err := fs.Mv(workingDir, username, path, newPath)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Failed to move file " + path
	}
	w.Write([]byte(output))
}

func rmHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	filepath := r.URL.Query().Get(FilePathParam)
	username, workingDir := getSessionInfo(w, r)
	var output string
	if err := fs.Rm(workingDir, username, filepath); err != nil || filepath == "" {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Failed to remove file " + filepath
	} else {
		output = "Done."
	}
	w.Write([]byte(output))
}

func addGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	groupname := r.URL.Query().Get(GroupNameParam)
	var output string
	if err := fs.AddGroup(groupname); err != nil || groupname == "" {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Failed to add group " + groupname
	} else {
		output = "Done."
	}
	w.Write([]byte(output))
}

func addUserToGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	groupname := r.URL.Query().Get(GroupNameParam)
	username := r.URL.Query().Get(UserNameParam)
	var output string
	if err := fs.AddUserToGroup(username, groupname); err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		output = "Failed to add user to group"
	} else {
		output = "Done."
	}
	w.Write([]byte(output))
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	if !session.SessionManager.SessionExists(w, r) {
		w.Write([]byte("Not logged in"))
		return
	}
	fileName := r.URL.Query().Get(FilePathParam)
	username, workingDir := getSessionInfo(w, r)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Write failed. please try again."))
		return
	}
	output, err := fs.Write(workingDir, username, fileName, data)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		w.Write([]byte("Write failed. please try again."))
		return
	}
	w.Write([]byte(output))
}

func getSessionInfo(w http.ResponseWriter, r *http.Request) (string, string) {
	sess := session.SessionManager.SessionStart(w, r)
	workingDir := sess.Get(session.WorkingDir)
	username := sess.Get(session.Username)
	return fmt.Sprint(username), fmt.Sprint(workingDir)
}

func changeWorkingDir(w http.ResponseWriter, r *http.Request, workingDir string) error {
	sess := session.SessionManager.SessionStart(w, r)
	err := sess.Set(session.WorkingDir, workingDir)
	if err != nil {
		log.Println(fmt.Errorf("error thrown: %w", err))
		return err
	}
	return nil
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/ls", lsHandler)
	http.HandleFunc("/pwd", pwdHandler)
	http.HandleFunc("/mkdir", mkdirHandler)
	http.HandleFunc("/cd", cdHandler)
	http.HandleFunc("/cat", catHandler)
	http.HandleFunc("/touch", touchHandler)
	http.HandleFunc("/mv", mvHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/rm", rmHandler)
	http.HandleFunc("/addgroup", addGroupHandler)
	http.HandleFunc("/addtogroup", addUserToGroupHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Printf("Open https://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
