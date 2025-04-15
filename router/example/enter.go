package example

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	CustomerRouter              //客户路由
	FileUploadAndDownloadRouter //文件上传和下载路由
	AttachmentCategoryRouter    //分类路由
}

var (
	exaCustomerApi              = api.ApiGroupApp.ExampleApiGroup.CustomerApi
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
	attachmentCategoryApi       = api.ApiGroupApp.ExampleApiGroup.AttachmentCategoryApi
)
