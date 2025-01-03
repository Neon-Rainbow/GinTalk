package logger

import (
	"GinTalk/settings"
	"bytes"
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// newLogger 创建并初始化zap日志库
func newLogger(cfg *settings.LoggerConfig) (*zap.Logger, error) {
	// 检查 OutputPaths 是否包含文件路径，并创建其目录
	for _, path := range cfg.OutputPaths {
		// 如果路径不是标准输出或标准错误，则处理为文件路径
		if path != "stdout" && path != "stderr" {
			logDir := filepath.Dir(path)
			if _, err := os.Stat(logDir); os.IsNotExist(err) {
				err = os.MkdirAll(logDir, 0755)
				if err != nil {
					return nil, fmt.Errorf("failed to create log directory: %w", err)
				}
			}
		}
	}

	// 设置日志级别,从配置文件中读取
	var level zapcore.Level

	logLevelMapper := map[int]zapcore.Level{
		-1: zap.DebugLevel,
		0:  zap.InfoLevel,
		1:  zap.WarnLevel,
		2:  zap.ErrorLevel,
		3:  zap.DPanicLevel,
		4:  zap.PanicLevel,
		5:  zap.FatalLevel,
	}

	if l, ok := logLevelMapper[cfg.Level]; ok {
		level = l
	} else {
		level = zap.InfoLevel
	}

	// 设置日志格式（JSON 或 Console）
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// 使用lumberjack进行日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.OutputPaths[1], // 假设 [1] 是日志文件路径，应该处理 OutputPaths 中文件类型的路径
		MaxSize:    cfg.MaxSize,        // 单个日志文件最大尺寸（MB）
		MaxBackups: cfg.MaxBackups,     // 保留的日志文件个数
		MaxAge:     cfg.MaxAge,         // 日志文件最大保留天数
		Compress:   cfg.Compress,       // 是否压缩旧日志文件
	}

	// 设置日志输出路径和错误输出路径
	cores := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),        // 输出到标准输出
			zapcore.AddSync(lumberJackLogger), // 输出到日志文件
		),
		level,
	)

	// 构建zap日志对象
	logger := zap.New(cores, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 由于c.Request.Body只允许读一次,读完之后内容会被删除,因此在这里我们读取Body中的内容,然后存储在变量中
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		cost := time.Since(start)
		c.Set("cost", cost) // 将耗时存储在context中, 以便后续使用
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.Any("head", c.Request.Header),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("body", string(bodyBytes)),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
			zap.String("request_id", c.Request.Header.Get("X-Request-Id")),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
