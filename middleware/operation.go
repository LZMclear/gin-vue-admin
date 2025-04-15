package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/segmentio/kafka-go"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

// OperationRecord 主要用来记录http请求和响应的相关信息
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet { //其他请求
			var err error
			//读取请求体，存入body变量中
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", zap.Error(err))
			} else {
				//将请求体重新设置到c.Request.Body中，以便后续处理函数可以继续读取请求体
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			//将请求的参数以键值对方式存储
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			//JSON序列化
			body, _ = json.Marshal(&m)
		}
		claims, _ := utils.GetClaims(c)
		if claims != nil && claims.BaseClaims.ID != 0 {
			userId = int(claims.BaseClaims.ID)
		} else {
			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		record := system.SysOperationRecordsKafka{
			CreatedAt: time.Now(),
			Ip:        c.ClientIP(),          //IP地址
			Method:    c.Request.Method,      //请求方法
			Path:      c.Request.URL.Path,    //请求路由
			Agent:     c.Request.UserAgent(), //请求代理
			Body:      "",
			UserID:    userId, //用户ID
		}

		// 上传文件时 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			record.Body = "[文件]"
		} else {
			if len(body) > bufferSize {
				record.Body = "[超出记录长度]"
			} else {
				record.Body = string(body)
			}
		}
		//创建自定义的responseBodyWriter用于记录响应内容
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		//将ResponseWriter改为自定义的writer，因为writer以匿名字段封装了ResponseWriter，所以可以直接调用它的方法
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now) //运行时间
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status() //响应状态码
		record.Latency = latency
		//解析响应体
		var resp response.Response
		err := json.Unmarshal(writer.body.Bytes(), &resp)
		if err != nil {
			fmt.Println(err)
		}
		record.Resp = resp

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
			if len(writer.body.String()) > bufferSize {
				// 截断
				record.Resp.Msg = "超出记录长度" //将record.Body ——>record.Resp
			}
		}
		//序列化为字节
		marshal, err := json.Marshal(record)
		if err != nil {
			global.GVA_LOG.Error("marshal json records error:", zap.Error(err))
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()
		// 执行kafka日志异步生产
		if err = global.GVA_WRITER.WriteMessages(ctx, kafka.Message{
			Value: marshal,
		}); err != nil {
			global.GVA_LOG.Info("write messages to kafka error:", zap.Error(err))
		}

		//if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
		//	global.GVA_LOG.Error("create operation record error:", zap.Error(err))
		//}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
