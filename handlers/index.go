package handlers

import (
	"net/http"
)

// IndexHandler - Handles calls to the index '/'
type IndexHandler struct {
	Handler
}

// ServeHTTP - Handles the http request made to the index
func (i IndexHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	i.Handler.Log.Infof("'%s' serving request on route '%s'", i.Name, i.Route)
	w.Write([]byte("you have successfully connected to Dohyo"))
}
