package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/service/class"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	ClassServiceGroup   class.ServiceGroup
	BlogServiceGroup    blog.ServiceGroup
}
