package sumux

import "net/http"

// Router -
type Router interface {
	Get(path string, f func(w http.ResponseWriter, r *http.Request))
	Post(path string, f func(w http.ResponseWriter, r *http.Request))
	Put(path string, f func(w http.ResponseWriter, r *http.Request))
	Delete(path string, f func(w http.ResponseWriter, r *http.Request))
	Patch(path string, f func(w http.ResponseWriter, r *http.Request))
	Group(path string, f func(r Router))
}
