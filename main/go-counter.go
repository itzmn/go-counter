package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-counter/config"
	"go-counter/internal"
	zlog "go-counter/logs"
	"go-counter/thirdpart"
	"log"
	"net/http"
	_ "net/http/pprof" // 导入 pprof 包
	"strconv"
	"time"
)

func main() {

	//test_counter_mode := flag.Bool("test_counter_mode", false, "enable test_counter mode, use default counter variable")
	start := time.Now().UnixNano()
	flag.Parse()
	if !zlog.InitLogger(&zlog.LogConf{}) {
		fmt.Println("init log err")
		return
	}
	defer zlog.Sync()

	if err := config.InitConfig(); err != nil {
		fmt.Println("init config err", err)
		return
	}
	// 加载变量
	if err := internal.LoadVariables(); err != nil {
		fmt.Println("init variable err", err)
		return
	}

	if err := thirdpart.InitRedis(); err != nil {
		fmt.Println("init redis err", err)
		return
	}

	// 初始化三方库
	http.HandleFunc("/counter", internal.Counter)
	// 注册 prometheus 监控指标
	http.Handle("/metrics", promhttp.Handler())
	cost := (time.Now().UnixNano() - start) / 1e6
	// 启动一个 HTTP 服务器，默认监听 6060 端口
	go func() {
		zlog.Log(zlog.LL_INFO, "start", "Starting pprof server at :"+strconv.Itoa(config.GetConfig().ServerPort+1))
		log.Println("Starting pprof server at :" + strconv.Itoa(config.GetConfig().ServerPort+1))
		if err := http.ListenAndServe(":"+strconv.Itoa(config.GetConfig().ServerPort+1), nil); err != nil {
			log.Println("Error starting pprof server:", err)
		}
	}()
	zlog.Log(zlog.LL_INFO, "start", "server start at :"+strconv.Itoa(config.GetConfig().ServerPort)+", cost:"+strconv.Itoa(int(cost)))
	log.Println("server at :" + strconv.Itoa(config.GetConfig().ServerPort) + ", cost:" + strconv.Itoa(int(cost)))
	// 启动服务
	if err := http.ListenAndServe(":"+strconv.Itoa(config.GetConfig().ServerPort), nil); err != nil {
		log.Println("server start err, ", err)
	}

}
