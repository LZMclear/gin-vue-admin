package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/router/class"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Class   class.RouterGroup
	Blog    blog.RouterGroup
}

// RouterGroupApp 初始化结构体，返回的是一个指向该结构体的指针
