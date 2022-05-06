package dragonSpider

import (
	"github.com/justinas/nosurf"
	"net/http"
	"strconv"
)

func (ds *DragonSpider) SessionLoadAndSave(next http.Handler) http.Handler {
	return ds.Session.LoadAndSave(next)
}

func (ds *DragonSpider) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	secure, _ := strconv.ParseBool(ds.config.cookie.Secure)

	//exclude csrf token
	csrfHandler.ExemptGlobs("/api/*/*", "/api/*/*/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   ds.config.cookie.domain,
	})

	return csrfHandler
}
