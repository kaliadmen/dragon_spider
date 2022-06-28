package dragonSpider

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (ds *DragonSpider) routes() http.Handler {
	mux := chi.NewRouter()
	//add default middleware
	//inserts a request id to each request
	mux.Use(middleware.RequestID)
	//gets ip of where request is coming from
	mux.Use(middleware.RealIP)

	if ds.Debug {
		//logs request to console
		mux.Use(middleware.Logger)
	}
	//recovers from panics
	mux.Use(middleware.Recoverer)

	//use sessions
	mux.Use(ds.SessionLoadAndSave)

	//use csrf
	mux.Use(ds.NoSurf)

	mux.Use(ds.CheckForMaintenanceMode)

	return mux
}

// Routes are dragonSpider specific routes, which are mounted in the routes file
// in DragonSpider applications
func Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/test-ds", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("it works!"))
		if err != nil {
			return
		}
	})
	return r
}
