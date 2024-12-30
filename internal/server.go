package internal

import (
	"io"
	"net/http"
)

func Counter(writer http.ResponseWriter, request *http.Request) {

	all, _ := io.ReadAll(request.Body)
	requestCount.WithLabelValues(request.Method, "200").Inc()
	response := process(string(all))

	writer.Write([]byte(response))

}
