package gee

import (
	"log"
	"time"
)

// 中间件模板，返回的是HandleFunc处理函数
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		//处理过程
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
