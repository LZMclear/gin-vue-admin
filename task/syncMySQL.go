package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

// ProcessedMessage 保存原始Kafka消息和业务数据
type ProcessedMessage struct {
	KafkaMessage kafka.Message
	Data         system.SysOperationRecord
}

const batchSize = 10 //每批次插入10条数据

/*
	存在的问题：
		事务与偏移量提交的原子性问题：插入成功但提交偏移量失败时会导致数据重复。
*/

// SyncMySQL 使用kafka作为日志中间件，采用至少一次的消费机制
func SyncMySQL(db *gorm.DB, r *kafka.Reader) error {
	s := time.Now().String()
	global.GVA_LOG.Info("每五分钟执行一次，开始执行同步操作", zap.String("开始时间：", s))

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	partitionBuffer := make(map[int][]ProcessedMessage) //按照分区缓存消息
	// TODO 从kafka集群中获取数据列表
	sum := 0
	for {
		message, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("正常结束读取任务")
				break //超时，获取不到，正常结束
			}
			return fmt.Errorf("读取消息失败: %w", err)
		}
		sum += 1
		// TODO 转换为SysOperationRecord格式
		var marshal system.SysOperationRecordsKafka
		if err = json.Unmarshal(message.Value, &marshal); err != nil {
			global.GVA_LOG.Error("解析kafka操作日志记录失败", zap.String("偏移量为：", string(message.Value)), zap.Error(err))
			// 解析失败直接跳过，否则会造成数据丢失
			continue
		}
		bytes, _ := json.Marshal(marshal.Resp)
		record := system.SysOperationRecord{
			GVA_MODEL:    global.GVA_MODEL{CreatedAt: marshal.CreatedAt},
			Ip:           marshal.Ip,
			Method:       marshal.Method,
			Path:         marshal.Path,
			Status:       marshal.Status,
			Latency:      marshal.Latency,
			Agent:        marshal.Agent,
			ErrorMessage: marshal.ErrorMessage,
			Body:         marshal.Body,
			Resp:         string(bytes),
			UserID:       marshal.UserID,
			User:         system.SysUser{},
		}
		partition := message.Partition
		partitionBuffer[partition] = append(partitionBuffer[partition], ProcessedMessage{
			KafkaMessage: message,
			Data:         record,
		})

		//按分区处理批次
		if len(partitionBuffer[partition]) >= batchSize {
			err := insertBatch(partitionBuffer[partition], db)
			if err != nil {
				return err
			}
			//提交分区最高偏移量
			if err = r.CommitMessages(ctx, partitionBuffer[partition][len(partitionBuffer[partition])-1].KafkaMessage); err != nil {
				return fmt.Errorf("commit message offset failed : %w", err)
			}
			//清空批次
			delete(partitionBuffer, partition)
		}
	}
	//处理剩余消息
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, ms := range partitionBuffer {
		if err := insertBatch(ms, db); err != nil {
			return err
		}
		//提交最后一条消息的偏移量
		err := r.CommitMessages(c, ms[len(ms)-1].KafkaMessage)
		if err != nil {
			return fmt.Errorf("commit message offset failed : %w", err)
		}
	}
	global.GVA_LOG.Info("执行完毕", zap.String("结束时间", time.Now().String()), zap.String("同步数据数量", strconv.Itoa(sum)))
	return nil
}

// TODO 存入数据库
func insertBatch(messages []ProcessedMessage, db *gorm.DB) error {
	if len(messages) == 0 {
		return nil
	}
	//转换数据结构
	records := make([]system.SysOperationRecord, 0, len(messages))
	for _, message := range messages {
		records = append(records, message.Data)
	}

	//开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true, //主键冲突则更新所有字段
		}).CreateInBatches(records, batchSize).Error; err != nil {
			return fmt.Errorf("批量插入失败: %w", err)
		}
		return nil
	})
}
