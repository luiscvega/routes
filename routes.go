package routes

import (
	"net/http"
	"regexp"
)

type Handler interface {
	Serve(http.ResponseWriter, *http.Request, map[string]string)
}

type route struct {
	method  string
	pattern *regexp.Regexp
	keys    []string
	handler Handler
}

type Routes []route

var re = regexp.MustCompile(`:(\w+)`)

func (rs *Routes) Add(method, path string, handler Handler) {
	r := route{
		method:  method,
		pattern: regexp.MustCompile("^" + re.ReplaceAllLiteralString(path, `([^\\/]+)`) + "/?$"),
		handler: handler,
	}

	result := re.FindAllStringSubmatch(path, -1)
	for _, tuple := range result {
		r.keys = append(r.keys, tuple[1])
	}

	*rs = append(*rs, r)
}

func (rs Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rs {
		if r.Method != route.method {
			continue
		}

		matches := route.pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			continue
		}

		params := map[string]string{}
		for i, value := range matches[1:] {
			key := route.keys[i]
			params[key] = value
		}

		route.handler.Serve(w, r, params)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("routes: not found"))
}
