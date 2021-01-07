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

// HandleFunc -
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// Method -
type Method map[string][]PathCallback

// PathCallback -
type PathCallback struct {
	path   string
	pathXp string
	f      func(w http.ResponseWriter, r *http.Request)
}

// ServeSumux -
type ServeSumux struct {
	HandleFuncs []HandleFunc
	groupPath   string
	pcm         Method
}

// NewServeSumux -
func NewServeSumux() *ServeSumux {
	pcm := make(Method)
	return &ServeSumux{pcm: pcm}
}

func (s ServeSumux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, pc := range s.pcm[r.Method] {
		if equalPath(r.URL.Path, pc.pathXp) {
			ctx := context.WithValue(r.Context(), pathWithKey, pc.path)
			pc.f(w, r.WithContext(ctx))
			return
		}
	}

	http.Error(w, "Path Not Found", 404)
}

// Get -
func (s *ServeSumux) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["GET"] = append(s.pcm["GET"], PathCallback{path, compiledPath(path), f})
}

// Post -
func (s *ServeSumux) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["POST"] = append(s.pcm["POST"], PathCallback{path, compiledPath(path), f})
}

// Put -
func (s *ServeSumux) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["PUT"] = append(s.pcm["PUT"], PathCallback{path, compiledPath(path), f})
}

// Delete -
func (s *ServeSumux) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["DELETE"] = append(s.pcm["DELETE"], PathCallback{path, compiledPath(path), f})
}

// Patch -
func (s *ServeSumux) Patch(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["PATCH"] = append(s.pcm["PATCH"], PathCallback{path, compiledPath(path), f})
}

// Options -
func (s *ServeSumux) Options(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(s.groupPath, path)
	s.pcm["OPTIONS"] = append(s.pcm["OPTIONS"], PathCallback{path, compiledPath(path), f})
}

// Group -
func (s *ServeSumux) Group(path string, f func(r Router)) {
	s.groupPath += path
	f(s)
	s.groupPath = ""
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

func equalPath(path, pathXp string) bool {
	if pathXp == "/" {
		return path == pathXp
	}

	rgx := regexp.MustCompile(pathXp)
	return rgx.MatchString(path)
}
