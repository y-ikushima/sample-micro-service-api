package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// LogConfig はログ設定を管理する構造体
type LogConfig struct {
	Level       string // debug, info, warn, error
	Environment string // development, production
	Service     string
	Version     string
}

// InitFromEnv は環境変数からログを初期化する便利関数
func InitFromEnv() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Init(LogConfig{
		Level:       logLevel,
		Environment: env,
		Service:     "app-service",
		Version:     "1.0.0",
	})
}

// Init はグローバルロガーを初期化
func Init(config LogConfig) error {
	var logger *zap.Logger
	var err error

	logger, err = newLogger(config)

	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// newLogger はCloud Run本のロガーを作成
func newLogger(config LogConfig) (*zap.Logger, error) {
	// Cloud Loggingに最適化されたEncoder設定
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity", // Cloud Loggingの標準フィールド
		NameKey:        "logger",
		CallerKey:      "sourceLocation",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		// Cloud Loggingの標準時刻形式（RFC3339）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339Nano))
		},
		// Cloud Loggingの標準重要度レベル
		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch level {
			case zapcore.DebugLevel:
				enc.AppendString("DEBUG")
			case zapcore.InfoLevel:
				enc.AppendString("INFO")
			case zapcore.WarnLevel:
				enc.AppendString("WARNING")
			case zapcore.ErrorLevel:
				enc.AppendString("ERROR")
			case zapcore.DPanicLevel, zapcore.PanicLevel:
				enc.AppendString("CRITICAL")
			case zapcore.FatalLevel:
				enc.AppendString("ALERT")
			default:
				enc.AppendString("DEFAULT")
			}
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 標準的なcaller形式を使用
	}

	// Cloud Loggingに最適化されたCore設定
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout), // Cloud Runはstdoutを自動収集
		parseLogLevel(config.Level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	
	// Cloud Loggingで認識される追加フィールド
	logger = logger.With(
		zap.String("service", config.Service),
		zap.String("version", config.Version),
	)

	return logger, nil
}


// parseLogLevel は文字列からzapcore.Levelに変換
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "critical":
		return zapcore.DPanicLevel
	case "alert":
		return zapcore.PanicLevel
	case "emergency":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// GetLogger はグローバルロガーを取得
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// デフォルトロガーで初期化
		Init(LogConfig{
			Level:       "info",
			Environment: "development",
			Service:     "unknown",
			Version:     "unknown",
		})
	}
	return globalLogger
}

// Info はInfoレベルのログを出力
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error はErrorレベルのログを出力
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug はDebugレベルのログを出力
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn はWarnレベルのログを出力
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal はFatalレベルのログを出力してプログラムを終了
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Sync はバッファされたログをフラッシュ（終了時に呼び出し推奨）
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// HttpRequest はCloud Loggingの標準httpRequestフィールド構造
type HttpRequest struct {
	RequestMethod string        `json:"requestMethod,omitempty"`
	RequestUrl    string        `json:"requestUrl,omitempty"`
	RequestSize   string        `json:"requestSize,omitempty"`
	Status        int           `json:"status,omitempty"`
	ResponseSize  string        `json:"responseSize,omitempty"`
	UserAgent     string        `json:"userAgent,omitempty"`
	RemoteIp      string        `json:"remoteIp,omitempty"`
	ServerIp      string        `json:"serverIp,omitempty"`
	Referer       string        `json:"referer,omitempty"`
	Latency       string        `json:"latency,omitempty"`
	Protocol      string        `json:"protocol,omitempty"`
}

// LogHttpRequest はHTTPリクエスト情報をCloud Logging形式でログ出力
func LogHttpRequest(msg string, req HttpRequest, fields ...zap.Field) {
	allFields := append(fields, zap.Any("httpRequest", req))
	GetLogger().Info(msg, allFields...)
} 