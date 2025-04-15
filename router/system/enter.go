package system

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	//全部采用匿名字段，这样这些匿名字段会将其类型的方法上升到结构体RouterGroup
	ApiRouter               //接口路由
	JwtRouter               //Jwt路由
	SysRouter               //系统路由
	BaseRouter              //基础路由
	InitRouter              //初始化路由
	MenuRouter              //菜单路由
	UserRouter              //用户路由
	CasbinRouter            //casbin路由
	AutoCodeRouter          //自动化代码路由
	AuthorityRouter         //授权路由
	DictionaryRouter        //字典路由
	OperationRecordRouter   //操作记录路由
	DictionaryDetailRouter  //字典细节路由
	AuthorityBtnRouter      //鉴权按钮路由
	SysExportTemplateRouter //系统导出模板路由
	SysParamsRouter         //系统参数路由
}

var (
	dbApi               = api.ApiGroupApp.SystemApiGroup.DBApi
	jwtApi              = api.ApiGroupApp.SystemApiGroup.JwtApi
	baseApi             = api.ApiGroupApp.SystemApiGroup.BaseApi
	casbinApi           = api.ApiGroupApp.SystemApiGroup.CasbinApi
	systemApi           = api.ApiGroupApp.SystemApiGroup.SystemApi
	sysParamsApi        = api.ApiGroupApp.SystemApiGroup.SysParamsApi
	autoCodeApi         = api.ApiGroupApp.SystemApiGroup.AutoCodeApi
	authorityApi        = api.ApiGroupApp.SystemApiGroup.AuthorityApi
	apiRouterApi        = api.ApiGroupApp.SystemApiGroup.SystemApiApi
	dictionaryApi       = api.ApiGroupApp.SystemApiGroup.DictionaryApi
	authorityBtnApi     = api.ApiGroupApp.SystemApiGroup.AuthorityBtnApi
	authorityMenuApi    = api.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	autoCodePluginApi   = api.ApiGroupApp.SystemApiGroup.AutoCodePluginApi
	autocodeHistoryApi  = api.ApiGroupApp.SystemApiGroup.AutoCodeHistoryApi
	operationRecordApi  = api.ApiGroupApp.SystemApiGroup.OperationRecordApi
	autoCodePackageApi  = api.ApiGroupApp.SystemApiGroup.AutoCodePackageApi
	dictionaryDetailApi = api.ApiGroupApp.SystemApiGroup.DictionaryDetailApi
	autoCodeTemplateApi = api.ApiGroupApp.SystemApiGroup.AutoCodeTemplateApi
	exportTemplateApi   = api.ApiGroupApp.SystemApiGroup.SysExportTemplateApi
)
