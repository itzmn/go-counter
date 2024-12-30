package internal

import "github.com/prometheus/client_golang/prometheus"

var (
	// 创建一个全局计数器
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)
)

func init() {
	// 注册计数器到 Prometheus
	prometheus.MustRegister(requestCount)
}
