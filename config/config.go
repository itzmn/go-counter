package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var Config *TotalConfig

type RedisConf struct {
	Host   string `json:"Host"`
	Port   int    `json:"Port"`
	Passwd string `json:"Passwd"`
}

//var configPath string

type TotalConfig struct {
	ServerPort int
	RedisConf  RedisConf
}

var configPath = flag.String("config", "./config/config.json", "config.json path")

func InitConfig() error {

	file, err := os.Open(*configPath)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	c := new(TotalConfig)
	err = json.Unmarshal(bytes, c)
	if err != nil {
		return err
	}
	Config = c
	return nil
}

func GetConfig() *TotalConfig {
	return Config
}
