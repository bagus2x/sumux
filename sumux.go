package sumux

import (
	"context"
	"net/http"
	"regexp"
)

type key int

const (
	pathWithKey key = iota
)

// PathCallback -
type PathCallback struct {
	path   string
	pathXp string
	f      func(w http.ResponseWriter, r *http.Request)
}

// ServeGeMux -
type ServeGeMux struct {
	groupPath string
	pcbPut    []PathCallback
	pcbGet    []PathCallback
	pcbPost   []PathCallback
	pcbPatch  []PathCallback
	pcbDelete []PathCallback
}

func (sg ServeGeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		for _, hp := range sg.pcbGet {
			if equalPath(r.URL.Path, hp.pathXp) {
				ctx := context.WithValue(r.Context(), pathWithKey, hp.path)
				hp.f(w, r.WithContext(ctx))
				return
			}
		}
		break

	case http.MethodPost:
		for _, hp := range sg.pcbPost {
			if equalPath(r.URL.Path, hp.pathXp) {
				ctx := context.WithValue(r.Context(), pathWithKey, hp.path)
				hp.f(w, r.WithContext(ctx))
				return
			}
		}
		break

	case http.MethodPut:
		for _, hp := range sg.pcbPut {
			if equalPath(r.URL.Path, hp.pathXp) {
				ctx := context.WithValue(r.Context(), pathWithKey, hp.path)
				hp.f(w, r.WithContext(ctx))
				return
			}
		}
		break

	case http.MethodDelete:
		for _, hp := range sg.pcbGet {
			if equalPath(r.URL.Path, hp.pathXp) {
				ctx := context.WithValue(r.Context(), pathWithKey, hp.path)
				hp.f(w, r.WithContext(ctx))
				return
			}
		}
		break

	case http.MethodPatch:
		for _, hp := range sg.pcbPatch {
			if equalPath(r.URL.Path, hp.pathXp) {
				ctx := context.WithValue(r.Context(), pathWithKey, hp.path)
				hp.f(w, r.WithContext(ctx))
				return
			}
		}
		break
	}

	http.Error(w, "Path Not Found", 404)
}

// Get -
func (sg *ServeGeMux) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcbGet = append(sg.pcbGet, PathCallback{path, compiledPath(path), f})
}

// Post -
func (sg *ServeGeMux) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcbPost = append(sg.pcbPost, PathCallback{path, compiledPath(path), f})
}

// Put -
func (sg *ServeGeMux) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcbPut = append(sg.pcbPut, PathCallback{path, compiledPath(path), f})
}

// Delete -
func (sg *ServeGeMux) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcbDelete = append(sg.pcbDelete, PathCallback{path, compiledPath(path), f})
}

// Patch -
func (sg *ServeGeMux) Patch(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcbPatch = append(sg.pcbPatch, PathCallback{path, compiledPath(path), f})
}

// Group -
func (sg *ServeGeMux) Group(path string, f func(r Router)) {
	sg.groupPath += path
	f(sg)
	sg.groupPath = ""
}

func compiledPath(path string) string {
	if path == "/" || path == "" {
		return "/"
	}

	rgx := regexp.MustCompile(`<[\w]+>`)
	return rgx.ReplaceAllString(path, `[\w]+`) + "$"
}

func concatPath(groupPath string, path string) string {
	if path != "/" {
		return groupPath + path
	}

	return groupPath
}

func equalPath(p1, p2 string) bool {
	if p2 == "/" {
		return p1 == p2
	}

	rgx := regexp.MustCompile(p2)
	return rgx.MatchString(p1)
}
