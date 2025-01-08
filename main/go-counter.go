package main

import (
	"context"
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
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var logDir = flag.String("log", "./log", "log dir")
	//test_counter_mode := flag.Bool("test_counter_mode", false, "enable test_counter mode, use default counter variable")
	start := time.Now().UnixNano()
	flag.Parse()
	if !zlog.InitLogger(&zlog.LogConf{
		LogDir: *logDir,
	}) {
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

	mux := http.NewServeMux()

	// 初始化三方库
	mux.HandleFunc("/counter", internal.Counter)
	// 注册 prometheus 监控指标
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{Addr: ":" + strconv.Itoa(config.GetConfig().ServerPort), Handler: mux}

	cost := (time.Now().UnixNano() - start) / 1e6
	// 启动一个 HTTP 服务器，默认监听 6060 端口
	go func() {
		zlog.Log(zlog.LL_INFO, "start", "Starting pprof server at :"+strconv.Itoa(config.GetConfig().ServerPort+1))
		log.Println("Starting pprof server at :" + strconv.Itoa(config.GetConfig().ServerPort+1))
		if err := http.ListenAndServe(":"+strconv.Itoa(config.GetConfig().ServerPort+1), nil); err != nil {
			log.Println("Error starting pprof server:", err)
		}
	}()
	go func() {
		zlog.Log(zlog.LL_INFO, "start", "server start at :"+strconv.Itoa(config.GetConfig().ServerPort)+", cost:"+strconv.Itoa(int(cost)))
		log.Println("server at :" + strconv.Itoa(config.GetConfig().ServerPort) + ", cost:" + strconv.Itoa(int(cost)))
		// 启动服务
		if err := server.ListenAndServe(); err != nil {
			log.Println("server start err, ", err)
			os.Exit(1)
		}

	}()

	// 创建一个信号通道，用于接收操作系统的信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // 捕获 SIGINT ctrlC 和 SIGTERM kill 信号

	sig := <-signalChan
	fmt.Println("receive signal: ", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("shut down err", err)
	} else {
		fmt.Println("shut down over", err)
	}

}
