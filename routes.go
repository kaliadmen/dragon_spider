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

	return mux
}
