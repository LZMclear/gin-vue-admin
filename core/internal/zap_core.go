package internal

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapCore struct { //创建的自定义类型，用于扩展和定制日志库的核心功能
	level zapcore.Level
	zapcore.Core
}

func NewZapCore(level zapcore.Level) *ZapCore { //创建一个新的核心配置
	entity := &ZapCore{level: level}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	//第三个参数，任何实现了Enabled方法的类型都是LevelEnabler类型的接口
	entity.Core = zapcore.NewCore(global.GVA_CONFIG.Zap.Encoder(), syncer, levelEnabler)
	return entity
}

// WriteSyncer 根据配置信息创建一个zapcore.WriteSyncer对象
// 用于将日志数据写入指定的目标位置。可以将日志同时输出到控制台和文件，也可以只输出到文件。
func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	cutter := NewCutter( //创建一个日志切割器对象
		global.GVA_CONFIG.Zap.Director,
		z.level.String(),
		global.GVA_CONFIG.Zap.RetentionDay,
		CutterWithLayout(time.DateOnly),
		CutterWithFormats(formats...),
	)
	if global.GVA_CONFIG.Zap.LogInConsole { //将日志同时输出到控制台和文件，使用 zapcore.NewMultiWriteSyncer 函数创建一个多写入同步器 multiSyncer，将 os.Stdout（标准输出，即控制台）和日志切割器 cutter 作为参数传入
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		return zapcore.AddSync(multiSyncer)
	}
	return zapcore.AddSync(cutter)
}

// Enabled 判断某个日志级别是否启用  z.level == level使得这个zapcore.Core实例仅处理单一级别的日志
func (z *ZapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

// With 创建一个带有额外字段的新核心。复用原生核心的With方法
func (z *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

// Check 检查是否应该记录某个日志条目，级别匹配则可添加核心
func (z *ZapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

// 将日志条目和字段写入日志。
// Entry表示完整的日志消息。条目的结构化上下文已经序列化，但日志级别、时间、消息和调用站点信息可供检查和修改。
func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			syncer := z.WriteSyncer(fields[i].String)
			z.Core = zapcore.NewCore(global.GVA_CONFIG.Zap.Encoder(), syncer, z.level)
		}
	}
	return z.Core.Write(entry, fields)
}

// Sync 同步日志数据。
func (z *ZapCore) Sync() error {
	return z.Core.Sync()
}
