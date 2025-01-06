package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	zlog "go-counter/logs"
	"go-counter/thirdpart"
)

type CounterC struct {
	Timestamp   int
	SlotTs      int
	MinSlotTs   int
	DimArrStr   string
	DataValue   interface{}
	ResultValue interface{}
	Filter      []struct {
		Func   string      `json:"func"`
		Path   string      `json:"path"`
		Params interface{} `json:"params"`
		Type   string      `json:"type"`
	} `json:"filter"`
	Function   string `json:"function"`
	Dimensions []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"dimensions"`
	Data struct {
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"data"`
	Window struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	} `json:"window"`
	Type string `json:"type"`
	Name string `json:"name"`
}


func (counterC *CounterC) doFilter(reqData *string) bool {
	for _, filter := range counterC.Filter {
		filterV, err := GetValueFromJSON(reqData, filter.Path, filter.Type)
		if err != nil {
			return false
		}
		if filter.Params != filterV {
			return false
		}
	}
	return true
}

func (counterC *CounterC) updateWindow(window interface{}) {
	updateStr, _ := json.Marshal(window)
	setErr := thirdpart.HMSetRedisData(counterC.DimArrStr, counterC.Name, updateStr)
	if setErr != nil {
		zlog.LogReqErr(zlog.LL_ERROR, "hmset", "", "set result err", setErr)
	}
}

func (counterC *CounterC) doData(reqData *string) bool {

	v := gjson.Get(*reqData, counterC.Data.Path)
	if !v.Exists() {
		return false
	}

	if counterC.Function == "count" {
		counterC.DataValue = 1
	} else {
		counterC.DataValue = v.Value()
		if counterC.Data.Type == "int" {
			counterC.DataValue = int(v.Int())
		}
	}
	return true
}

func (counterC *CounterC) prePareSlot(reqData *string) {
	ts := gjson.Get(*reqData, "timestamp").Int()

	// 获取时间槽
	window := counterC.Window.Size
	slotSize := getSlotSize(window)

	slotTime := int(ts) / 1000 / slotSize * slotSize
	counterC.SlotTs = slotTime
	counterC.MinSlotTs = slotTime - window

}

func (counterC *CounterC) count(fieldV interface{}) {
	if counterC.Function == "count" {
		counterC.DataValue = 1
	}
	window := make([]ValueSlot, 0)
	total := 0
	dv := counterC.DataValue.(int)
	update := false
	if fieldV != nil {
		if err := json.Unmarshal([]byte(fieldV.(string)), &window); err != nil {
		}

		for i := range window {
			if window[i].Ts == counterC.SlotTs {
				update = true
				window[i].Val += dv
				total += window[i].Val
			}
			if window[i].Ts >= counterC.MinSlotTs && window[i].Ts < counterC.SlotTs {
				total += window[i].Val
			}
		}
	}

	if !update {
		v := ValueSlot{
			Ts:  counterC.SlotTs,
			Val: dv,
		}
		window = append(window, v)
		total += dv
	}
	counterC.ResultValue = total
	counterC.updateWindow(window)

}

func (counterC *CounterC) counterFloat(fieldV interface{}) {

	update := false
	window := make([]FloatSlot, 0)
	dv := counterC.DataValue.(float32)
	var total = float32(0)

	if fieldV != nil {
		if err := json.Unmarshal([]byte(fieldV.(string)), &window); err != nil {
			return
		}

		for i := range window {
			if window[i].Ts == counterC.SlotTs {
				update = true
				window[i].Val += dv
				total += window[i].Val
			} else if window[i].Ts >= counterC.MinSlotTs && window[i].Ts < counterC.SlotTs {
				total += window[i].Val
			}
		}
	}

	// 更新window且写回redis
	if !update {
		v := FloatSlot{
			Ts:  counterC.SlotTs,
			Val: dv,
		}
		total += dv
		window = append(window, v)
	}
	counterC.ResultValue = total
	go counterC.updateWindow(window)
}

func (counterC *CounterC) distinct(fieldV interface{}) {
	window := make([]DistinctSlot, 0)
	dv := counterC.DataValue.(string)
	update := false
	tmpMap := map[string]struct{}{}
	if fieldV != nil {
		if err := json.Unmarshal([]byte(fieldV.(string)), &window); err != nil {
		}

		for i := range window {
			if window[i].Val == dv {
				update = true
				window[i].Ts = counterC.SlotTs
			}
			tmpMap[window[i].Val] = struct{}{}
		}
	}

	if !update {
		v := DistinctSlot{
			Ts:  counterC.SlotTs,
			Val: dv,
		}
		window = append(window, v)
		tmpMap[dv] = struct{}{}
	}
	counterC.ResultValue = len(tmpMap)
	go counterC.updateWindow(window)

}

// GetValueFromJSON 从 JSON 字符串中提取指定键的值，并将其转换为目标类型
func GetValueFromJSON(jsonStr *string, key string, targetType string) (interface{}, error) {
	// 定义一个 map 来解析 JSON 数据
	var data map[string]interface{}

	// 解析 JSON 字符串
	err := json.Unmarshal([]byte(*jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// 检查 key 是否存在
	rawValue, exists := data[key]
	if !exists {
		return nil, errors.New("key not found in JSON")
	}

	// 根据目标类型转换
	switch targetType {
	case "string":
		value, ok := rawValue.(string)
		if !ok {
			return nil, errors.New("value type is not string")
		}
		return value, nil
	case "int":
		// JSON 默认将数字解析为 float64，需要手动转换
		floatValue, ok := rawValue.(float64)
		if !ok {
			return nil, errors.New("value type is not numeric")
		}
		return int(floatValue), nil
	case "float64":
		value, ok := rawValue.(float64)
		if !ok {
			return nil, errors.New("value type is not float64")
		}
		return value, nil
	case "bool":
		value, ok := rawValue.(bool)
		if !ok {
			return nil, errors.New("value type is not bool")
		}
		return value, nil
	default:
		return nil, errors.New("unsupported target type")
	}
}
