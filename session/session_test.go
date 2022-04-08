package session

import (
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSession_InitSession(t *testing.T) {
	var sesMan *scs.SessionManager
	var sesKind reflect.Kind
	var sesType reflect.Type

	s := &Session{
		CookieName:     "Dragon Spider",
		CookieDomain:   "localhost",
		CookieLifetime: "100",
		CookiePersist:  "true",
		//CookieSecure:   "false",
		SessionType: "cookie",
	}

	ses := s.InitSession()
	reflectVal := reflect.ValueOf(ses)

	for reflectVal.Kind() == reflect.Ptr || reflectVal.Kind() == reflect.Interface {
		sesKind = reflectVal.Kind()
		sesType = reflectVal.Type()

		reflectVal = reflectVal.Elem()
	}

	if !reflectVal.IsValid() {
		t.Error("invalid kind or type:", reflectVal.Kind(), "type:", reflectVal.Type())
	}

	if sesKind != reflect.ValueOf(sesMan).Kind() {
		t.Error("Testing cookie session: Wrong kind returned. Expected", reflect.ValueOf(sesMan).Kind(), "got", sesKind)
	}

	if sesType != reflect.ValueOf(sesMan).Type() {
		t.Error("Testing cookie session: Wrong type returned. Expected", reflect.ValueOf(sesMan).Type(), "got", sesType)
	}

}
