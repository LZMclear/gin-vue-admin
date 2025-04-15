package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

func SyncES() error {
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
	err := s.SyncES(global.GVA_ES_READER)
	if err != nil {
		return err
	}
	return nil
}
