package session

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	CookieName     string
	CookieLifetime string
	CookiePersist  string
	CookieSecure   string
	CookieDomain   string
	SessionType    string
}

func (s *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	//how long does session last
	minutes, err := strconv.Atoi(s.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	//does cookie persist
	if strings.ToLower(s.CookiePersist) == "true" {
		persist = true
	}
	//do cookie have to be secure
	if strings.ToLower(s.CookieSecure) == "true" {
		secure = true
	}

	//create session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Secure = secure
	session.Cookie.Name = s.CookieName
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	//session store
	switch strings.ToLower(s.SessionType) {
	case "redis":

	case "mysql", "mariadb":

	case "postgres", "postgresql":

	default:
		//use cookies
	}

	return session
}
