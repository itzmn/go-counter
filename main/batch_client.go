package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sendRequest(url string, userCnt int, activityCnt int, amountRand int) {

	// 将请求体编码为 JSON

	start := time.Now().UnixNano()
	str := `{"activityId":"%v", "amount":%v, "organization":"itzmn", "user":"%v", "timestamp":%v}`
	// 使用当前时间戳作为随机数种子
	rand.Seed(time.Now().UnixNano())

	user := "zhangsan" + strconv.Itoa(rand.Intn(userCnt))
	activity := "ac" + strconv.Itoa(rand.Intn(activityCnt))

	body := fmt.Sprintf(str, activity, rand.Intn(amountRand), user, time.Now().UnixNano()/1e6)

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	requestId := gjson.Get(string(bytes), "requestId").String()
	defer resp.Body.Close()

	cost := time.Now().UnixNano() - start
	// 打印响应状态
	fmt.Printf("requestId: %v, Response status: %d, cost: %v\n", requestId, resp.StatusCode, cost/1e6)
}

func main() {
	// 要请求的 URL
	url := "http://localhost:19999/counter" // 替换为你要请求的目标 URL

	// 持续请求的总时长
	duration := 5 * time.Minute

	// 每秒发起的请求数
	qps := 10

	userCnt := 10000
	activityCnt := 100
	amountRand := 10000

	// 每秒启动 10 个 Goroutines 来发起请求
	ticker := time.NewTicker(time.Second) // 每秒钟触发一次

	// 持续 5 分钟
	endTime := time.Now().Add(duration)

	// 使用 Goroutines 来并发请求
	for time.Now().Before(endTime) {
		<-ticker.C
		// 每秒发起 10 次请求
		for i := 0; i < qps; i++ {
			go sendRequest(url, userCnt, activityCnt, amountRand)
		}
	}

	// 停止 ticker
	ticker.Stop()
}
