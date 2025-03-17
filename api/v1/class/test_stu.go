package class

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/class"
	classReq "github.com/flipped-aurora/gin-vue-admin/server/model/class/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StudentApi struct{}

// CreateStudent 创建学生
// @Tags Student
// @Summary 创建学生
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body class.Student true "创建学生"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /stu/createStudent [post]
func (stuApi *StudentApi) CreateStudent(c *gin.Context) {
	var stu class.Student
	err := c.ShouldBindJSON(&stu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = stuService.CreateStudent(&stu)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteStudent 删除学生
// @Tags Student
// @Summary 删除学生
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body class.Student true "删除学生"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /stu/deleteStudent [delete]
func (stuApi *StudentApi) DeleteStudent(c *gin.Context) {
	ID := c.Query("ID")
	err := stuService.DeleteStudent(ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteStudentByIds 批量删除学生
// @Tags Student
// @Summary 批量删除学生
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /stu/deleteStudentByIds [delete]
func (stuApi *StudentApi) DeleteStudentByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := stuService.DeleteStudentByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateStudent 更新学生
// @Tags Student
// @Summary 更新学生
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body class.Student true "更新学生"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /stu/updateStudent [put]
func (stuApi *StudentApi) UpdateStudent(c *gin.Context) {
	var stu class.Student
	err := c.ShouldBindJSON(&stu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = stuService.UpdateStudent(stu)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindStudent 用id查询学生
// @Tags Student
// @Summary 用id查询学生
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询学生"
// @Success 200 {object} response.Response{data=class.Student,msg=string} "查询成功"
// @Router /stu/findStudent [get]
func (stuApi *StudentApi) FindStudent(c *gin.Context) {
	ID := c.Query("ID")
	restu, err := stuService.GetStudent(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(restu, c)
}

// GetStudentList 分页获取学生列表
// @Tags Student
// @Summary 分页获取学生列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query classReq.StudentSearch true "分页获取学生列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /stu/getStudentList [get]
func (stuApi *StudentApi) GetStudentList(c *gin.Context) {
	var pageInfo classReq.StudentSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := stuService.GetStudentInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetStudentPublic 不需要鉴权的学生接口
// @Tags Student
// @Summary 不需要鉴权的学生接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /stu/getStudentPublic [get]
func (stuApi *StudentApi) GetStudentPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	stuService.GetStudentPublic()
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的学生接口信息",
	}, "获取成功", c)
}

// NewStudentMethod 测试学生新方法
// @Tags Student
// @Summary 测试学生新方法
// @Accept application/json
// @Produce application/json
// @Param data query classReq.StudentSearch true "成功"
// @Success 200 {object} response.Response{data=object,msg=string} "成功"
// @Router /stu/newMethod [GET]
func (stuApi *StudentApi) NewStudentMethod(c *gin.Context) {
	// 请添加自己的业务逻辑
	err := stuService.NewStudentMethod()
	if err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("失败", c)
		return
	}
	response.OkWithData("返回数据", c)
}
