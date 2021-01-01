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

// PathCallbackMap -
type PathCallbackMap map[string]PathCallback

// PathCallback -
type PathCallback struct {
	method string
	path   string
	f      func(w http.ResponseWriter, r *http.Request)
}

// ServeSumux -
type ServeSumux struct {
	HandleFuncs []HandleFunc
	groupPath   string
	pcm         PathCallbackMap
}

// NewServeSumux -
func NewServeSumux() *ServeSumux {
	pcm := make(PathCallbackMap)
	return &ServeSumux{pcm: pcm}
}

func (sg ServeSumux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range sg.pcm {
		if equalPath(r.URL.Path, k) && v.method == r.Method {
			ctx := context.WithValue(r.Context(), pathWithKey, v.path)
			v.f(w, r.WithContext(ctx))
			return
		}
	}

	http.Error(w, "Path Not Found", 404)
}

// Get -
func (sg *ServeSumux) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	cp := compiledPath(path)
	sg.pcm[cp] = PathCallback{
		method: "GET",
		path:   path,
		f:      f,
	}
}

// Post -
func (sg *ServeSumux) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	cp := compiledPath(path)
	sg.pcm[cp] = PathCallback{
		method: "POST",
		path:   path,
		f:      f,
	}
}

// Put -
func (sg *ServeSumux) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	cp := compiledPath(path)
	sg.pcm[cp] = PathCallback{
		method: "PUT",
		path:   path,
		f:      f,
	}
}

// Delete -
func (sg *ServeSumux) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	cp := compiledPath(path)
	sg.pcm[cp] = PathCallback{
		method: "DELETE",
		path:   path,
		f:      f,
	}
}

// Patch -
func (sg *ServeSumux) Patch(path string, f func(w http.ResponseWriter, r *http.Request)) {
	path = concatPath(sg.groupPath, path)
	cp := compiledPath(path)
	sg.pcm[cp] = PathCallback{
		method: "PATCH",
		path:   path,
		f:      f,
	}
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
