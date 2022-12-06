// Package api provides structure of the api server
package api

import "net/http"

// every version of the api should implement this interface
type ApiInterface interface {
	getVersion() string
	registerHandle(patterns string, f func(w http.ResponseWriter, r *http.Request))
	StartServer(addr string)
}
