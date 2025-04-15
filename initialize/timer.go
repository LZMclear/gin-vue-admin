package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option                    // 用于存储定时器的配置选项
		option = append(option, cron.WithSeconds()) //该配置允许在corn中使用秒级精度
		// 清理DB定时任务  @daily corn表达式，表示任务每天执行一次
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 同步Kafka数据到Mysql数据库定时任务  @daily corn表达式，表示任务每天执行一次
		_, err = global.GVA_Timer.AddTaskByFunc("SyncMySQL", "@daily", func() {
			err := task.SyncMySQL(global.GVA_DB, global.GVA_MYSQL_READER)
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时同步Kafka【日志】数据到MySQL数据库", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		//同步kafka数据到ES index库定时任务
		_, err = global.GVA_Timer.AddTaskByFunc("SyncES", "0 */5 * * * *", func() {
			err := task.SyncES()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "同步Kafka【日志】数据到ES索引库中", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}
	}()
}
