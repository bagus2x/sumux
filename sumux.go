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

func (sg ServeSumux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, pc := range sg.pcm[r.Method] {
		if equalPath(r.URL.Path, pc.pathXp) {
			ctx := context.WithValue(r.Context(), pathWithKey, pc.path)
			pc.f(w, r.WithContext(ctx))
			return
		}
	}

	http.Error(w, "Path Not Found", 404)
}

// Get -
func (sg *ServeSumux) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcm["GET"] = append(sg.pcm["GET"], PathCallback{path, compiledPath(path), f})
}

// Post -
func (sg *ServeSumux) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcm["POST"] = append(sg.pcm["POST"], PathCallback{path, compiledPath(path), f})
}

// Put -
func (sg *ServeSumux) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcm["PUT"] = append(sg.pcm["PUT"], PathCallback{path, compiledPath(path), f})
}

// Delete -
func (sg *ServeSumux) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcm["DELETE"] = append(sg.pcm["DELETE"], PathCallback{path, compiledPath(path), f})
}

// Patch -
func (sg *ServeSumux) Patch(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	sg.pcm["PATCH"] = append(sg.pcm["PATCH"], PathCallback{path, compiledPath(path), f})
}

// Group -
func (sg *ServeSumux) Group(path string, f func(r Router)) {
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

func equalPath(path, pathXp string) bool {
	if pathXp == "/" {
		return path == pathXp
	}

	rgx := regexp.MustCompile(pathXp)
	return rgx.MatchString(path)
}
