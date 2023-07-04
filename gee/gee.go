package gee

import (
	"net/http"
)

// 定义HandlerFunc为路由的请求函数
type HandlerFunc func(*Context)

// router 路由映射表，由请求方法和静态路由地址构成。如GET-/hello
type enginer struct {
	router *router
}

func New() *enginer {
	return &enginer{router: newRouter()}
}

func (engine *enginer) addRoute(method, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// 访问方法
func (engine *enginer) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *enginer) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *enginer) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute("PUT", pattern, handler)
}
func (engine *enginer) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute("DELETE", pattern, handler)
}

// 启动监听
func (engine *enginer) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *enginer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
