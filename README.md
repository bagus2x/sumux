# Simple Mux "sumux"

## Getting Started
```bash
$ go get https://github.com/bagus2x/sumux
```

run example

```bash
$ git clone https://github.com/bagus2x/sumux.git
$ cd sumux
$ make example

```

## Example
    
[link](https://github.com/bagus2x/sumux/tree/main/example)

### Init

```go
r := sumux.NewMux()
```

### Method

```go
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    res.JSON(w, 200, map[string]string{"greeting": "test 1 2 3"})
})

r.Get("/mw", mw(func(w http.ResponseWriter, r *http.Request) {
    res.Plain(w, 200, "hello world")
}))

r.Post("/whoami", func(w http.ResponseWriter, r *http.Request) {
    b, _ := ioutil.ReadAll(r.Body)
    res.Plain(w, 200, string(b))
})
```

### Group & Param

```go
r.Group("/api/v1/user", func(r sumux.Router) {
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        res.JSON(w, 200, map[string]int{"oneminusone": 1})
    })

    r.Get("/<name>", func(w http.ResponseWriter, r *http.Request) {
        name, _ := sumux.Param(r, "name")
        res.JSON(w, 200, map[string]string{"name": name})
    })

    r.Get("/<name>/abc/<address>", func(w http.ResponseWriter, r *http.Request) {
        p, _ := sumux.Params(r)
        res.JSON(w, 200, p)
    })
})
```

### Middleware

```go
type Handler func(w http.ResponseWriter, r *http.Request)

func mw(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hai from local mw")
		next(w, r)
	}
}

func mw2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hai from global mw")
		next.ServeHTTP(w, r)
	})
}
```