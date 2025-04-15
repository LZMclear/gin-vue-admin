package task

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func initEsReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"43.138.33.205:9092", "43.138.33.205:9093", "43.138.33.205:9094"},
		GroupID:        "test",
		Topic:          "op-records",
		StartOffset:    kafka.FirstOffset,
		SessionTimeout: 5 * time.Minute,
		CommitInterval: 0, //关闭自动提交
		MinBytes:       10e3,
		MaxBytes:       10e6,
	})
}

// 创建elasticsearch连接
func initESClient() *elasticsearch.TypedClient {
	config := elasticsearch.Config{
		Addresses: []string{"http://43.138.33.205:9200"},
	}
	//创建客户端连接
	client, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("create elasticsearch client success")
	return client
}

// 测试Elasticsearch同步
func TestSyncES(t *testing.T) {
	s := system.SyncEs{
		EsClient: initESClient(),
		Cfg: struct {
			Index      string
			BatchSize  int
			MaxRetries int
		}{
			Index:      "op-records",
			BatchSize:  30,
			MaxRetries: 3,
		},
	}
	r := initEsReader()
	defer r.Close()
	err := s.SyncES(r)
	if err != nil {
		t.Error(err)
	}
	t.Log("sync es success")
}
