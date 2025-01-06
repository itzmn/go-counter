package internal

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	zlog "go-counter/logs"
	"go-counter/thirdpart"
	"go-counter/util"
	"sort"
	"strings"
	"sync"
	"time"
)

func process(requestData string) string {

	start := time.Now().UnixNano()
	reqId := util.GenReqId()
	response := requestData
	defer func() {
		zlog.LogReqEnd(zlog.LL_INFO, "RE", reqId, "start end", response, start)
	}()

	response, _ = sjson.Set(response, "requestId", reqId)
	zlog.LogReqStart(zlog.LL_INFO, "RB", reqId, "start", requestData)

	var waitGroup sync.WaitGroup

	VariableMutex.RLock()
	var resultChan = make(chan map[string]interface{}, len(counterVariables))
	for dim, statisticVars := range counterVariables {

		waitGroup.Add(1)
		go statisticOneDim(reqId, dim, statisticVars, &requestData, resultChan, &waitGroup)

	}
	VariableMutex.RUnlock()

	// 开启计时器
	//timeout := time.After(50 * time.Millisecond)

	// 当所有任务线程执行完成时候，关闭结果channel
	go func(wg *sync.WaitGroup, c chan map[string]interface{}) {
		wg.Wait()
		close(c)
	}(&waitGroup, resultChan)

	//time.Sleep(2*time.Second)
	over := false
	for true {
		select {
		case oneMap, ok := <-resultChan:
			if !ok {
				over = true
			}
			for k, v := range oneMap {
				response, _ = sjson.Set(response, k, v)
			}
		}
		if over {
			break
		}
	}

	return response
}

type DimProcessor struct {
	DimArrStr   string
	ProcessVars []*CounterC
}

func (p *DimProcessor) processCounter(counterC *CounterC, reqData *string) {

	filterRes := counterC.doFilter(reqData)
	if !filterRes {
		return
	}
	haveData := counterC.doData(reqData)
	if !haveData {
		return
	}

	counterC.DimArrStr = p.DimArrStr
	counterC.prePareSlot(reqData)

	p.ProcessVars = append(p.ProcessVars, counterC)
}

func (p *DimProcessor) statisticCounter() map[string]interface{} {

	var vars []string

	for _, counterC := range p.ProcessVars {
		vars = append(vars, counterC.Name)
	}
	resMap := make(map[string]interface{}, 0)

	fieldVs := thirdpart.HMGetRedisData(p.DimArrStr, vars...)
	for index, counterC := range p.ProcessVars {
		fieldV := fieldVs[index]

		if counterC.Function == "count" || counterC.Function == "sum" {
			//counterC.DataValue = 1
			//
			//if counterC.Type == "int" {
			counterC.count(fieldV)
			//} else if counterC.Type == "float32" {
			//	counterC.counterFloat(fieldV)
			//} else {
			//	continue
			//}

		} else if counterC.Function == "distinct" {
			counterC.distinct(fieldV)
		} else {
			continue
		}

		resMap[counterC.Name] = counterC.ResultValue
	}

	return resMap
}

func statisticOneDim(reqId, dim string, counterCs []CounterC, requestData *string,
	resultChan chan map[string]interface{}, waitGroup *sync.WaitGroup) {

	dimStart := time.Now().UnixNano()
	defer waitGroup.Done()

	dimProcessor := new(DimProcessor)

	// dim处理提前处理
	var dimArr []string
	dims := strings.Split(dim, VariableSplit)
	for _, dimension := range dims {
		r := gjson.Get(*requestData, dimension)
		if !r.Exists() {
			return
		}
		dimStr := r.String()
		dimArr = append(dimArr, dimStr)
	}
	sort.Strings(dimArr)
	dimArrStr := strings.Join(dimArr, VariableSplit)
	dimProcessor.DimArrStr = "counter:" + util.Md5Sum(dimArrStr)

	for _, counterC := range counterCs {
		dimProcessor.processCounter(&counterC, requestData)
	}
	if len(dimProcessor.ProcessVars) == 0 {
		return
	}
	statisticCounterRes := dimProcessor.statisticCounter()

	zlog.LogReq(zlog.LL_INFO, "process_dim", reqId, fmt.Sprintf("process dim over, dimName: %v, dimValue: %v, cost:%v",
		dimProcessor.DimArrStr, dimProcessor.DimArrStr,
		(time.Now().UnixNano()-dimStart)/1e6))
	resultChan <- statisticCounterRes

}
