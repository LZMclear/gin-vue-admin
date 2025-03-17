package class

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ StudentRouter }

var stuApi = api.ApiGroupApp.ClassApiGroup.StudentApi
