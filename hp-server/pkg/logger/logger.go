package logger

import (
	"bufio"
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ExpandLogger struct {
	inner *zap.SugaredLogger
}

func (e *ExpandLogger) Printf(template string, args ...interface{}) {
	e.inner.With("module", "gorm").Infof(template, args...)
}

var (
	Logger  *zap.SugaredLogger
	ELogger *ExpandLogger
	Debug   func(args ...interface{})
	Info    func(args ...interface{})
	Warn    func(args ...interface{})
	Error   func(args ...interface{})
	Fatal   func(args ...interface{})
	Panic   func(args ...interface{})
	Debugf  func(template string, args ...interface{})
	Infof   func(template string, args ...interface{})
	Warnf   func(template string, args ...interface{})
	Errorf  func(template string, args ...interface{})
	Fatalf  func(template string, args ...interface{})
	Panicf  func(template string, args ...interface{})
)

func init() {
	Logger = NewBufferRotateLogger(context.Background(), "data/logs/hp-server.log", 64, 8)
	ELogger = &ExpandLogger{
		inner: Logger,
	}
	Debug = Logger.Debug
	Info = Logger.Info
	Warn = Logger.Warn
	Error = Logger.Error
	Fatal = Logger.Fatal
	Panic = Logger.Panic
	Debugf = Logger.Debugf
	Infof = Logger.Infof
	Warnf = Logger.Warnf
	Errorf = Logger.Errorf
	Fatalf = Logger.Fatalf
	Panicf = Logger.Panicf
}

// 带缓冲、切分的日志记录器
func NewBufferRotateLogger(ctx context.Context, path string, rotateSize int, rotateCount int) *zap.SugaredLogger {
	lumLog := &lumberjack.Logger{
		Filename:   path,        // 日志文件路径
		MaxSize:    rotateSize,  // 文件最大尺寸（MB）
		MaxBackups: rotateCount, // 备份文件最大数量
		Compress:   true,        // 是否压缩/归档旧文件
	}

	buffer := bufio.NewWriterSize(lumLog, 128*1024)
	rwmutex := sync.RWMutex{}
	// 定时刷新缓存，写入到文件中
	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				rwmutex.Lock()
				buffer.Flush()
				rwmutex.Unlock()
			}
		}
	}()

	// 创建编码器
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// 创建zapcore
	core := zapcore.NewCore(encoder, zapcore.AddSync(&syncWriter{writer: buffer, rwMu: &rwmutex}), zapcore.DebugLevel)

	return zap.New(core).Sugar()
}

// 自定义的同步写入器
type syncWriter struct {
	writer *bufio.Writer
	rwMu   *sync.RWMutex
}

func (s *syncWriter) Write(p []byte) (n int, err error) {
	s.rwMu.RLock()         // 读锁
	defer s.rwMu.RUnlock() // 解锁
	return s.writer.Write(p)
}
