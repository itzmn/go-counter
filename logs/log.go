package zlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	// log level
	LL_DEBUG = "[DEBUG]"
	LL_INFO  = "[INFO]"
	LL_WARN  = "[WARN]"
	LL_ERROR = "[ERROR]"
	LL_FATAL = "[FATAL]"
)

type zloger interface {
	log(level, stage, info string)
	logReq(level, stage, reqId, info string)
	logReqData(level, stage, reqId, info string, data interface{})
	logReqStart(level, stage, reqId, info string, reqData interface{})
	logReqEnd(level, stage, reqId, info string, retData interface{}, startTs int64)
	logReqErr(level, stage, reqId, info string, err interface{})
	sync()
}

type _TYPE_ZAP_LOG_fUNC func(msg string, fields ...zap.Field)

type kvLogger struct {
	*zap.Logger
	logFuncMap map[string]_TYPE_ZAP_LOG_fUNC
}

func (k *kvLogger) logReqData(level, stage, reqId, info string, data interface{}) {
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("reqId", reqId),
		zap.String("info", info),
		zap.Any("data", data),
	)
}

func (k *kvLogger) logReqStart(level, stage, reqId, info string, reqData interface{}) {
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("reqId", reqId),
		zap.String("info", info),
		zap.Any("reqData", reqData),
	)
}
func (k *kvLogger) logReqEnd(level, stage, reqId, info string, retData interface{}, startTs int64) {
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("reqId", reqId),
		zap.String("info", info),
		zap.Any("retData", retData),
		zap.Int64("cost", (time.Now().UnixNano()-startTs)/1e6),
	)
}

func (k *kvLogger) log(level, stage, info string) {
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("info", info),
	)
}
func (k *kvLogger) logReq(level, stage, reqId, info string) {
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("reqId", reqId),
		zap.String("info", info),
	)
}
func (k *kvLogger) logReqErr(level, stage, reqId, info string, err interface{}) {
	//panic("implement me")
	k.getLogFunc(level)(level,
		zap.String("stage", stage),
		zap.String("reqId", reqId),
		zap.String("info", info),
		zap.Any("err", err),
	)
}

func (k *kvLogger) sync() {
	k.Logger.Sync()
}

func (k kvLogger) getLogFunc(level string) _TYPE_ZAP_LOG_fUNC {
	unc, ok := k.logFuncMap[level]
	if !ok {
		k.Logger.Error("ERROR", zap.String("stage", "zlog"),
			zap.String("errInfo", "get level logFunc err ["+level+"], use info"))
		return k.Logger.Info
	}
	return unc
}

func getLogWriter(conf *LogConf) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.LogDir + "/" + conf.FileName,
		MaxSize:    conf.MaxMB * 1024 * 1024,
		MaxBackups: conf.MaxBackups,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func createKVLogger(conf *LogConf) zloger {

	// 1、创建config
	encoderConfig := &zapcore.EncoderConfig{
		TimeKey: "time",
		//LevelKey:    "level",
		MessageKey: "logLev", // 标识阶段
		CallerKey:  "file",
		//StacktraceKey: "stack",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// 2、使用自定义的 encoder 创建core
	core := zapcore.NewCore(
		newKeyValueEncoder(encoderConfig),
		zapcore.AddSync(getLogWriter(conf)),
		zapcore.DebugLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)) //zap.AddStacktrace(zap.WarnLevel)
	kvlogger := &kvLogger{
		Logger: logger,
	}
	kvlogger.logFuncMap = make(map[string]_TYPE_ZAP_LOG_fUNC, 5)
	kvlogger.logFuncMap[LL_DEBUG] = logger.Debug
	kvlogger.logFuncMap[LL_INFO] = logger.Info
	kvlogger.logFuncMap[LL_WARN] = logger.Warn
	kvlogger.logFuncMap[LL_ERROR] = logger.Error
	kvlogger.logFuncMap[LL_FATAL] = logger.Error
	return kvlogger
}

func createJsonLogger() *zap.Logger {
	// 1、 创建logger配置
	config := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 日志级别大写
		EncodeCaller:  zapcore.ShortCallerEncoder,  // 短路径，文件名+行号
	}
	// 2、创建日志编码器
	jsonEncoder := zapcore.NewJSONEncoder(config)
	// 2、 创建日志 core
	core := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	return logger

}

func main() {

	//jsonLogger := createJsonLogger()
	//defer jsonLogger.Sync()
	//jsonLogger.Info("is msg", zap.String("name", "zhangsan"))

	//logger := createKVLogger()
	//defer logger.Sync()

	// 示例日志
	//logger.Info("INFO", zap.String("version", "1.0.0"))
	//logger.Warn("WARN", zap.String("disk", "/dev/sda1"), zap.Int64("free_space", 100), zap.String("port", "8080"))
	//logger.Error("ERROR", zap.String("file", "/var/log/syslog"), zap.Int("retry", 3))

	//logger.Log(LL_INFO, "mysql", "read mysql success")
}
