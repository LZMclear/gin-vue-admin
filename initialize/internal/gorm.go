package internal

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
// Author [SliverHorn](https://github.com/SliverHorn)
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		general = global.GVA_CONFIG.Mysql.GeneralDB
	case "pgsql":
		general = global.GVA_CONFIG.Pgsql.GeneralDB
	case "oracle":
		general = global.GVA_CONFIG.Oracle.GeneralDB
	case "sqlite":
		general = global.GVA_CONFIG.Sqlite.GeneralDB
	case "mssql":
		general = global.GVA_CONFIG.Mssql.GeneralDB
	default:
		general = global.GVA_CONFIG.Mysql.GeneralDB
	}
	return &gorm.Config{
		//配置GORM的日志记录器，参数为一个Write类型的接口和一个配置结构体
		//实现了Print方法的结构体都是Write类型的接口
		Logger: logger.New(NewWriter(general), logger.Config{
			SlowThreshold: 200 * time.Millisecond, //设置慢查询的阈值为200秒，超过200秒的查询都被视为慢查询
			LogLevel:      general.LogLevel(),     //根据数据库配置设置日志记录级别
			Colorful:      true,                   //启用彩色日志输出
		}),
		NamingStrategy: schema.NamingStrategy{ //命名策略
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
