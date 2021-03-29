package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	Username   = "username"
	WorkingDir = "wd"
)

var SessionManager *Manager

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifetime int)
}

var providers = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("session: Requier called twice for provider " + name)
	}
	providers[name] = provider
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifetime int
}

func NewManager(providerName string, cookieName string, maxLifetime int) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", providerName)
	}
	return &Manager{
		cookieName:  cookieName,
		provider:    provider,
		maxLifetime: maxLifetime,
	}, nil
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		sess, _ := manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: manager.maxLifetime}
		http.SetCookie(w, &cookie)
		return sess
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		sess, _ := manager.provider.SessionRead(sid)
		return sess
	}
}

func (manager *Manager) SessionEnd(w http.ResponseWriter, r *http.Request) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		return err
	}
	sessionId, _ := url.QueryUnescape(cookie.Value)
	err = manager.provider.SessionDestroy(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifetime)
	time.AfterFunc(time.Duration(manager.maxLifetime), func() { manager.GC() })
}

func (manager *Manager) SessionExists(w http.ResponseWriter, r *http.Request) bool {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return false
	}
	_, err = manager.provider.SessionRead(sid)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
