package storage

//
//import (
//	"encoding/json"
//	"flag"
//	"fmt"
//	"go-counter/internal"
//	zlog "go-counter/logs"
//	"io/ioutil"
//	"os"
//	"sort"
//	"strings"
//	"sync"
//	"time"
//)
//
//// 先从 json中获取 需要统计的变量，后续迁移至数据库管理
//var variablesPath = flag.String("variables", "./config/statisticVars.json", "statisticVars.json path")
//
//var counterVariables map[string][]internal.CounterC
//var VariableMutex sync.RWMutex
//
//func LoadVariables() error {
//
//	go func() {
//		for {
//			time.Sleep(60 * time.Second)
//			if err := loadVariables(); err != nil {
//				zlog.Log(zlog.LL_ERROR, "LoadVariables", fmt.Sprintf("LoadVariables err, err: %v", err))
//			}
//		}
//	}()
//
//	return loadVariables()
//}
//
//func loadVariables() error {
//	start := time.Now().UnixNano()
//	file, err := os.Open(*variablesPath)
//	if err != nil {
//		return err
//	}
//	bytes, err := ioutil.ReadAll(file)
//	if err != nil {
//		return err
//	}
//	tmp := make([]internal.CounterC, 0)
//	if err = json.Unmarshal(bytes, &tmp); err != nil {
//		return err
//	}
//
//	tmpMap := make(map[string][]internal.CounterC)
//
//	for _, c := range tmp {
//		dimArr := []string{}
//		for _, dimension := range c.Dimensions {
//			if dimension.Path != "" {
//				dimArr = append(dimArr, dimension.Path)
//			}
//		}
//		sort.Strings(dimArr)
//		join := strings.Join(dimArr, ":")
//		v, ok := tmpMap[join]
//		if !ok {
//			v = make([]internal.CounterC, 0, 1)
//		}
//		v = append(v, c)
//		tmpMap[join] = v
//	}
//
//	VariableMutex.Lock()
//	counterVariables = tmpMap
//	VariableMutex.Unlock()
//
//	zlog.Log(zlog.LL_INFO, "loadVariables", fmt.Sprintf("load variables success, load variable dim cnt: %v, cost:%v", len(tmp), (time.Now().UnixNano()-start)/1e6))
//	return nil
//}
//
//func GetVariables() map[string][]internal.CounterC {
//	return counterVariables
//}
