package sfs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const HOST = "http://2b502.yeg.rac.sh:8080"
const (
	SessionIdCookieName = "sessionid"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	Client    *http.Client
	Username  string
	SessionId string
	SignedIn  bool
}

func NewClient() *Client {
	client := http.Client{}
	return &Client{
		Client:   &client,
		SignedIn: false,
	}
}

func (client *Client) Signup(username string, password string) (string, error) {
	credentials, _ := json.Marshal(Credentials{
		Username: username,
		Password: password,
	})
	output, err := client.runPostCommand("/signup", map[string]string{}, credentials)
	if err != nil {
		return "", err
	}
	client.SignedIn = true
	return output, nil
}

func (client *Client) Login(username string, password string) (string, error) {
	credentials, _ := json.Marshal(Credentials{
		Username: username,
		Password: password,
	})
	output, err := client.runPostCommand("/login", map[string]string{}, credentials)
	if err != nil {
		return "", err
	}
	client.SignedIn = true
	return output, nil
}

func (client *Client) Logout() (string, error) {
	output, _ := client.runGetCommand("/logout", map[string]string{})
	return output, nil
}

func (client *Client) Ls() (string, error) {
	if output, err := client.runGetCommand("/ls", map[string]string{}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Pwd() (string, error) {
	if output, err := client.runGetCommand("/pwd", map[string]string{}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Cat(filename string) (string, error) {
	if output, err := client.runGetCommand("/cat", map[string]string{"filepath": filename}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Touch(filename string) (string, error) {
	if output, err := client.runGetCommand("/touch", map[string]string{"filepath": filename}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Mkdir(dirname string) (string, error) {
	if output, err := client.runGetCommand("/mkdir", map[string]string{"filepath": dirname}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Cd(filename string) (string, error) {
	if output, err := client.runGetCommand("/cd", map[string]string{"filepath": filename}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Mv(filename string, newFileName string) (string, error) {
	if output, err := client.runGetCommand("/mv", map[string]string{"filepath": filename, "newpath": newFileName}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Rm(path string) (string, error) {
	if output, err := client.runGetCommand("/rm", map[string]string{"filepath": path}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) AddGroup(groupname string) (string, error) {
	if output, err := client.runGetCommand("/addgroup", map[string]string{"groupname": groupname}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) AddUserToGroup(username string, groupname string) (string, error) {
	if output, err := client.runGetCommand("/addtogroup", map[string]string{"username": username, "groupname": groupname}); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) Write(path string, data string) (string, error) {
	if output, err := client.runPostCommand("/write", map[string]string{"filepath": path}, []byte(data)); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (client *Client) runPostCommand(path string, args map[string]string, body []byte) (string, error) {
	res, err := client.post(path, args, body)
	if err != nil {
		return "", errors.New("Failed to run command")
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", errors.New("Failed to run command")
	}
	return buf.String(), nil
}

func (client *Client) runGetCommand(path string, args map[string]string) (string, error) {
	res, err := client.get(path, args, nil)
	if err != nil {
		return "", errors.New("Failed to run command ls")
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", errors.New("Failed to run command ls")
	}
	return buf.String(), nil
}

func (client *Client) get(path string, args map[string]string, body []byte) (*http.Response, error) {
	req := client.prepareRequest("GET", path, body)
	query := req.URL.Query()
	for name, value := range args {
		query.Add(name, value)
	}
	req.URL.RawQuery = query.Encode()
	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	client.updateSessionId(res)
	return res, err
}

func (client *Client) post(path string, args map[string]string, body []byte) (*http.Response, error) {
	req := client.prepareRequest("POST", path, body)
	query := req.URL.Query()
	for name, value := range args {
		query.Add(name, value)
	}
	req.URL.RawQuery = query.Encode()
	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	client.updateSessionId(res)
	return res, err
}

func (client *Client) updateSessionId(res *http.Response) {
	for _, cookie := range res.Cookies() {
		if cookie.Name == SessionIdCookieName {
			client.SessionId = cookie.Value
			client.SignedIn = true
		}
	}
}

func (client *Client) prepareRequest(method string, path string, body []byte) *http.Request {
	// Prepare request as needed, incude any headers, final encryption, etc.
	req, _ := http.NewRequest(method, HOST+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if client.SessionId != "" {
		req.AddCookie(&http.Cookie{Name: SessionIdCookieName, Value: client.SessionId})
	}
	return req
}
