package system

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type SyncEs struct {
	EsClient *elasticsearch.TypedClient
	Cfg      struct {
		Index      string
		BatchSize  int
		MaxRetries int
	}
}

// IndexExist 检查索引是否存在
func (s *SyncEs) IndexExist() (bool, error) {
	//创建请求
	request := esapi.IndicesExistsRequest{
		Index: []string{s.Cfg.Index},
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	response, err := request.Do(ctx, s.EsClient)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.IsError() {
		//索引不存在返回404
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, errors.New(response.String())
	}
	return true, nil
}

// CreateIndex 创建索引
func (s *SyncEs) CreateIndex() error {
	//先判断索引库是否存在
	if IsExist, _ := s.IndexExist(); IsExist {
		return nil
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	indexFalse := false
	analyzer := "ik_smart"
	//使用面向对象创建索引库
	response, err := s.EsClient.Indices.Create(s.Cfg.Index).Mappings(&types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": types.NewDateProperty(),
			"ip":         types.NewKeywordProperty(),
			"method":     types.NewKeywordProperty(),
			"path": types.TextProperty{
				Analyzer: &analyzer,
				CopyTo:   []string{"fulltext"},
			},
			"status":  types.NewIntegerNumberProperty(),
			"latency": types.NewLongNumberProperty(),
			"agent": &types.TextProperty{
				Index: &indexFalse,
			},
			"errorMessage": &types.TextProperty{
				Analyzer: &analyzer,
				CopyTo:   []string{"fulltext"},
			},
			"body": &types.TextProperty{
				Analyzer: &analyzer,
				CopyTo:   []string{"fulltext"},
			},
			"resp": &types.TypeMapping{
				Properties: map[string]types.Property{
					"code": types.KeywordProperty{Index: &indexFalse},
					"data": types.NewObjectProperty(),
					"msg": &types.TextProperty{
						Analyzer: &analyzer,
						CopyTo:   []string{"fulltext"},
					},
				},
			},
			"userId": types.NewKeywordProperty(),
			"fulltext": types.TextProperty{
				Analyzer: &analyzer,
			},
		},
	}).Do(ctx)
	if err != nil {
		global.GVA_LOG.Error("create index failed", zap.Any("err", err))
		return err
	}
	global.GVA_LOG.Info("create index success", zap.String("Index:", response.Index))
	return nil
}

type ESMessage struct {
	KafkaMessage kafka.Message
	Data         []byte
}

// SyncES 同步ES  拉取kafka数据，随后直接生成ES文档
func (s *SyncEs) SyncES(r *kafka.Reader) error {
	fmt.Println(r.Config().GroupID)
	//拉取请求总共一分钟时间，超时说明拉取完毕，没有拉取的下次拉取
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	partitionBuffer := make(map[int][]ESMessage)
	sum := 0
	for {
		message, err := r.FetchMessage(ctx)
		fmt.Println(string(message.Value))
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, context.DeadlineExceeded) {
				break //表明正常读取数据完毕，结束读取
			}
			return err
		}
		sum++
		partition := message.Partition
		partitionBuffer[partition] = append(partitionBuffer[partition], ESMessage{
			KafkaMessage: message,
			Data:         message.Value,
		})
		//按分区批处理
		if len(partitionBuffer[partition]) >= s.Cfg.BatchSize {
			fmt.Println("partitionBuffer[", partition, "]")
			err := s.bulkInsertToES(partitionBuffer[partition])
			if err != nil {
				return err
			}
			//提交分区偏移量最高的消息
			if err = r.CommitMessages(ctx, partitionBuffer[partition][len(partitionBuffer[partition])-1].KafkaMessage); err != nil {
				return fmt.Errorf("commit message offset failed : %w", err)
			}
			//清空批次
			delete(partitionBuffer, partition)
		}
	}
	//处理剩余消息
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, ms := range partitionBuffer {
		err := s.bulkInsertToES(ms)
		if err != nil {
			return err
		}
		//提交分区内偏移量最高的消息
		if err = r.CommitMessages(c, ms[len(ms)-1].KafkaMessage); err != nil {
			return fmt.Errorf("commit message offset failed : %w", err)
		}
		fmt.Println("commit message offset success")
	}
	global.GVA_LOG.Info("sync ES success", zap.Int("sync data num: ", sum))
	return nil
}

// 批量创建ES文档
func (s *SyncEs) bulkInsertToES(ems []ESMessage) error {
	var bulkBuffer bytes.Buffer
	for _, em := range ems {
		//构建每个文档的元数据
		meta := []byte(fmt.Sprintf(`{ "index": { "_index": "%s" } }%s`,
			s.Cfg.Index, "\n"))
		//元数据添加到请求体中
		bulkBuffer.Write(meta)

		em.Data = append(em.Data, '\n')
		_, err := bulkBuffer.Write(em.Data)
		if err != nil {
			return err
		}
	}
	//创建批量插入请求
	request := esapi.BulkRequest{
		Index: s.Cfg.Index,
		Body:  bytes.NewReader(bulkBuffer.Bytes()),
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	//执行请求
	response, err := request.Do(ctx, s.EsClient)
	if err != nil {
		return fmt.Errorf("perform bulk insert error: %w", err)
	}
	all, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(all))
	defer response.Body.Close()
	return nil
}
