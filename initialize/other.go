package initialize

import (
	"bufio"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"os"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

func OtherInit() {
	//将过期时间转换为时间戳格式
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	//将缓冲时间转换为时间格式
	_, err = utils.ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	//初始化一个本地缓存
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr), //设置缓存默认过期时间
	)
	file, err := os.Open("go.mod")
	if err == nil && global.GVA_CONFIG.AutoCode.Module == "" {
		scanner := bufio.NewScanner(file)                                                 //创建扫描器对象scanner，用于逐行扫描文件
		scanner.Scan()                                                                    //读取文件第一行
		global.GVA_CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ") //去除module
	}
}
