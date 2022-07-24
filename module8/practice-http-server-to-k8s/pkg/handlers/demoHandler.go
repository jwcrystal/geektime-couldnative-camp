package handlers

import (
	"context"
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, res *http.Request)

type Engine struct {
	router map[string]HandlerFunc
	server *http.Server
	//middlewares []HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) addRoute(method, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.router[key] = handler
}

// GET method request
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func (e *Engine) Run(addr string) (err error) {
	e.server = &http.Server{Addr: addr, Handler: e}
	return e.server.ListenAndServe()
	//return http.ListenAndServe(addr, e)
}

func (e *Engine) Shutdown(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}

//// Use is defined to add middleware to the Engine
//func (e *Engine) Use(middlewares ...HandlerFunc) { // 註冊中間件
//	//e.middlewares = append(e.middlewares, middlewares...)
//}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		// using middleware is a good choice, but skip this part here
		//e.middlewares = append(e.middlewares, handler)

		// Logger
		LogHandler(http.StatusOK, r)
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s\n", r.URL)
	}
}
