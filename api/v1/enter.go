package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/class"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	ClassApiGroup   class.ApiGroup
	BlogApiGroup    blog.ApiGroup
}
