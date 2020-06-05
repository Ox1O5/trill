package trill

import (
	"fmt"
	"net/http"
)

//定义了trill要使用的handle func
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

//实现HttpServer的引擎
type Engine struct {
	router map[string]HandlerFunc
}

//引擎构造函数
func New() *Engine{
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (engine *Engine)addRoute(method string, pattern string, handler HandlerFunc)  {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine)GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine)POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine)Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL.Path)
	}
}



