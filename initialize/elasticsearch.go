package initialize

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

// InitESClient 创建elasticsearch连接
func InitESClient() *elasticsearch.TypedClient {
	config := elasticsearch.Config{
		Addresses: []string{global.GVA_CONFIG.Elasticsearch.Address},
	}
	//创建客户端连接
	client, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		global.GVA_LOG.Error("create elasticsearch client failed", zap.Any("err", err))
		return nil
	}
	global.GVA_LOG.Info("create elasticsearch client success")
	return client
}

// InitES 初始化ES
func InitES() {
	s := system.SyncEs{
		EsClient: global.GVA_ES_CLIENT,
		Cfg: struct {
			Index      string
			BatchSize  int
			MaxRetries int
		}{
			Index:      global.GVA_CONFIG.Elasticsearch.Index,
			BatchSize:  global.GVA_CONFIG.Elasticsearch.BatchSize,
			MaxRetries: global.GVA_CONFIG.Elasticsearch.MaxRetries,
		},
	}
	//创建索引库
	err := s.CreateIndex()
	if err != nil {
		global.GVA_LOG.Error("create elasticsearch index failed", zap.Any("err", err))
	}
	global.GVA_LOG.Info("create elasticsearch index success")
}
