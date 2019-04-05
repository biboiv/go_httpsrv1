package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"
	//	gl "wt/wt_global"
)

func GenerateStringId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type Session struct {
	sessionId string
	token     string
	dateStart time.Time
	isLogin   bool
	data      map[string]interface{}
}

func (s *Session) GetId() string {
	return s.sessionId
}

func (s *Session) LogOut() {
	s.isLogin = false
}

func (s *Session) GetNewToken() string {
	s.token = GenerateStringId()
	return s.token
}

func (s *Session) SetToken(token string) error {
	s.token = token
	return nil
}
func (s *Session) GetToken() string {
	return s.token
}

func (s *Session) Set(key string, val interface{}) error {
	s.data[key] = val
	return nil
}

func (s Session) Get(key string) interface{} {
	return s.data[key]
}

func (s Session) GetInt64(key string) int64 {
	var r int64 = 0
	v, ok := s.data[key]
	if ok {
		r = v.(int64)
	}
	return r
}

func (s Session) GetFloat64(key string) float64 {
	var r float64 = 0
	v, ok := s.data[key]
	if ok {
		r = InterfaceToFloat64(v)
	}
	return r
}

func (s Session) GetDateTime(key string) time.Time {
	var r time.Time = NilDateTime()
	v, ok := s.data[key]
	if ok {
		r = InterfaceToDateTime(v)
	}
	return r
}

func (s Session) GetString(key string) string {
	var r string = ""
	v, ok := s.data[key]
	if ok {
		r = InterfaceToString(v)
	}
	return r
}

type sessionManager struct {
	cookieName         string
	lock               sync.Mutex
	maxLifeTimeMinutes int64
	sessions           map[string]*Session
}

var sesMan sessionManager //sington

func (sm *sessionManager) getSession(w http.ResponseWriter, r *http.Request) *Session {
	//fmt.Println("GetSession...")

	cookie, err := r.Cookie(sm.cookieName)
	if err == nil {
		if cookie != nil {
			id := cookie.Value
			s, ok := sm.sessions[id]
			if ok == true {
				return s
			}
		}
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sessionId := GenerateStringId()
	cookie = &http.Cookie{
		Name:    sm.cookieName,
		Value:   sessionId,
		Expires: time.Now().Add(time.Duration(sm.maxLifeTimeMinutes) * time.Minute),
	}
	http.SetCookie(w, cookie)
	//var s *Session
	s := new(Session)
	s.sessionId = sessionId
	s.isLogin = false
	s.data = make(map[string]interface{})

	sm.sessions[sessionId] = s
	return s
}

func IfSessionManagerInit() bool {
	return sesMan.cookieName != ""
}

func InitSessionManager2(cookieName string, maxLifeTimeMinutes int64) {
	if !IfSessionManagerInit() {
		sesMan.cookieName = cookieName
		sesMan.maxLifeTimeMinutes = maxLifeTimeMinutes
		sesMan.sessions = make(map[string]*Session)
		fmt.Println("Create Session Manager - Ok")
	}
}

func InitSessionManager() {
	InitSessionManager2(GenerateStringId(), 1440) //1440=60*24
}

func GetSession(w http.ResponseWriter, r *http.Request) *Session {
	if IfSessionManagerInit() {
		ss := sesMan.getSession(w, r)
		//fmt.Println(ss.GetId())
		return ss
	}
	return nil
}
