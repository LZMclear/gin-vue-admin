package class

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ StudentApi }

var stuService = service.ServiceGroupApp.ClassServiceGroup.StudentService
