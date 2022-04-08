package dragonSpider

import "net/http"

func (ds *DragonSpider) SessionLoadAndSave(next http.Handler) http.Handler {
	return ds.Session.LoadAndSave(next)
}
