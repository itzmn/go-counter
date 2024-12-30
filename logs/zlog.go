package zlog

import (
	"os"
	"path/filepath"
)

var gZLoger zloger

type LogConf struct {
	LogDir     string
	MaxMB      int
	MaxBackups int
	FileName   string
}

func (c *LogConf) reset(logConf *LogConf) {
	c.LogDir = "./log"
	if logConf.LogDir != "" {
		c.LogDir = logConf.LogDir
	}
	c.MaxMB = 100
	if logConf.MaxMB != 0 {
		c.MaxMB = logConf.MaxMB
	}
	c.MaxBackups = 5
	if logConf.MaxBackups != 0 {
		c.MaxBackups = logConf.MaxBackups
	}
	c.FileName = filepath.Base(os.Args[0]) + ".log"
	if logConf.FileName != "" {
		c.FileName = logConf.FileName
	}

}

func InitLogger(logConf *LogConf) bool {
	conf := new(LogConf)
	conf.reset(logConf)
	logger := createKVLogger(conf)
	gZLoger = logger
	return gZLoger != nil
}

func Log(level, stage, info string) {
	gZLoger.log(level, stage, info)
}
func LogReq(level, stage, reqId, info string) {
	gZLoger.logReq(level, stage, reqId, info)
}

func LogReqStart(level, stage, reqId, info string, reqData interface{}) {
	gZLoger.logReqStart(level, stage, reqId, info, reqData)
}

func LogReqEnd(level, stage, reqId, info string, retData interface{}, startNs int64) {
	gZLoger.logReqEnd(level, stage, reqId, info, retData, startNs)
}

func LogReqErr(level, stage, reqId, info string, err interface{}) {
	gZLoger.logReqErr(level, stage, reqId, info, err)
}

func Sync() {
	gZLoger.sync()
}
