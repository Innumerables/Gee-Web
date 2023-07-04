package gee

import (
	"net/http"
	"strings"
)

//在开始使用的是map[string]HandlerFunc来存储路由表，但实现不了动态路由
// type router struct {
// 	handlers map[string]HandlerFunc
// }

//	func newRouter() *router {
//		return &router{
//			handlers: make(map[string]HandlerFunc),
//		}
//	}

// 前缀树实现动态路由
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//初始的增加路由
// func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
// 	log.Printf("Router %4s - %s", method, pattern)
// 	key := method + "-" + pattern
// 	r.handlers[key] = handler
// }

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	// key := c.Method + "-" + c.Path
	// if handler, ok := r.handlers[key]; ok {
	// 	handler(c)
	// } else {
	// 	c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	// }
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
