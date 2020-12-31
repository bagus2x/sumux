package sumux

import (
	"net/http"
)

// Mux -
type Mux struct {
	ServeGeMux
	middls []func(http.Handler) http.Handler
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer simpleLog(r.Method, r.URL.Path)()
	var h http.Handler = &m.ServeGeMux

	for i := len(m.middls) - 1; i >= 0; i-- {
		h = m.middls[i](h)
	}

	h.ServeHTTP(w, r)
}

// Use -
func (m *Mux) Use(h func(http.Handler) http.Handler) {
	m.middls = append(m.middls, h)
}

// NewMux -
func NewMux() Mux {
	return Mux{}
}
