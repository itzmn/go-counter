package internal

import (
	"fmt"
	"go-counter/thirdpart"
	"testing"
)

func Test_process(t *testing.T) {

	thirdpart.InitRedis()

	requestData := `{"activityId":"ac123", "amount":12, "name":"zhangsan", "user":"zhangsan", "timestamp":1733989840000, "requestId":"r123"}`
	response := process(requestData)
	fmt.Println(response)

}
