package sumux

import (
	"net/http"
)

// Mux -
type Mux struct {
	*ServeSumux
	adapter []func(http.Handler) http.Handler
	Log     bool
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Log {
		defer simpleLog(r.Method, r.URL.Path)()
	}

	var h http.Handler = m.ServeSumux

	for i := len(m.adapter) - 1; i >= 0; i-- {
		h = m.adapter[i](h)
	}

	h.ServeHTTP(w, r)
}

// Use -
func (m *Mux) Use(h func(http.Handler) http.Handler) {
	m.adapter = append(m.adapter, h)
}

// NewMux -
func NewMux() Mux {
	return Mux{ServeSumux: NewServeSumux()}
}
