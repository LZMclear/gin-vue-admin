package global

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB        *gorm.DB
	GVA_DBList    map[string]*gorm.DB
	GVA_REDIS     redis.UniversalClient
	GVA_REDISList map[string]redis.UniversalClient
	GVA_MONGO     *qmgo.QmgoClient
	GVA_CONFIG    config.Server
	GVA_VP        *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()  //time.Timer类型，直接进行了初始化
	GVA_Concurrency_Control             = &singleflight.Group{} //实现并发控制，避免多个goroutine 同时执行相同的操作。
	GVA_ROUTERS             gin.RoutesInfo
	GVA_ACTIVE_DBNAME       *string
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
	GVA_WRITER              *kafka.Writer
	GVA_MYSQL_READER        *kafka.Reader
	GVA_ES_READER           *kafka.Reader
	GVA_ES_CLIENT           *elasticsearch.TypedClient
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := GVA_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}
