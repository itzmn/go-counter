package zlog

import (
	"testing"
)

func TestLog(t *testing.T) {

	InitLogger(&LogConf{})

	Log(LL_INFO, "test", "test msg")
	Log(LL_ERROR, "test", "test msg")

	//logger := createJsonLogger()
	//logger.Info("haha")

}
