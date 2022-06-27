package dragonSpider

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"strconv"
	"strings"
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

func (ds *DragonSpider) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maintenanceMode {
			if !strings.Contains(r.URL.Path, "/public/maintenance.html") {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", ds.RootPath))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
