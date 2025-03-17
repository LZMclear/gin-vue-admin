package class

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type StudentRouter struct{}

func (s *StudentRouter) InitStudentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	stuRouter := Router.Group("stu").Use(middleware.OperationRecord())
	stuRouterWithoutRecord := Router.Group("stu")
	stuRouterWithoutAuth := PublicRouter.Group("stu")
	{
		stuRouter.POST("createStudent", stuApi.CreateStudent)
		stuRouter.DELETE("deleteStudent", stuApi.DeleteStudent)
		stuRouter.DELETE("deleteStudentByIds", stuApi.DeleteStudentByIds)
		stuRouter.PUT("updateStudent", stuApi.UpdateStudent)
	}
	{
		stuRouterWithoutRecord.GET("findStudent", stuApi.FindStudent)
		stuRouterWithoutRecord.GET("getStudentList", stuApi.GetStudentList)
	}
	{
		stuRouterWithoutAuth.GET("getStudentPublic", stuApi.GetStudentPublic)
		stuRouterWithoutAuth.GET("newMethod", stuApi.NewStudentMethod)
	}
}
