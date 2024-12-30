package thirdpart

import (
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	InitRedis()

	data := HMGetRedisData("counter:234", "test4", "test2")
	fmt.Println(data)
	HMSetRedisData("counter:235", "test1", "aaa", "test2", "bb", "test3", "[{\"ts\":1733726600,\"value\":5}]")

	data = HMGetRedisData("counter:235", "test3", "test2")
	fmt.Println(data)
}

//func TestHMGetRedisData(t *testing.T) {
//	InitRedis()
//	type args struct {
//		key    string
//		fields []string
//	}
//	tests := []struct {
//		name string
//		args args
//		want string
//	}{
//		// TODO: Add test cases.
//		{name: "successCase", args: args{key: "counter:234", fields: []string{"test3", "test2"}}, want: `{"test2":"[{\"ts\":1733726600,\"value\":5}]"}`},
//		{name: "fatalCase", args: args{key: "counter:234", fields: []string{"test3", "test2"}}, want: `{"test2":"[{\"ts\":1733726600,\"value\":6}]"}`},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := HMGetRedisData(tt.args.key, tt.args.fields...); got != tt.want {
//				t.Errorf("HMGetRedisData() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
