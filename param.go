package sumux

import (
	"net/http"
	"regexp"
	"strings"
)

// ParamResult -
type ParamResult map[string]string

// Params -
func Params(r *http.Request) (ParamResult, bool) {
	pKey, ok := r.Context().Value(pathWithKey).(string)
	if !ok {
		return nil, false
	}

	v := strings.Split(r.URL.Path, "/")
	pr := make(map[string]string)

	for i, k := range strings.Split(pKey, "/") {
		rgx, err := regexp.Compile(`<[\w]+>`)
		if err != nil {
			return nil, false
		}

		if rgx.MatchString(k) {
			pr[k[1:len(k)-1]] = v[i]
		}
	}

	return pr, true
}

// Param -
func Param(r *http.Request, key string) (string, bool) {
	pKey, ok := r.Context().Value(pathWithKey).(string)
	if !ok {
		return "", false
	}

	v := strings.Split(r.URL.Path, "/")
	for i, k := range strings.Split(pKey, "/") {
		if k != "" {
			if k[1:len(k)-1] == key {
				return v[i], true
			}
		}
	}

	return "", false
}
