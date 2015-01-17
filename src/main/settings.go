package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	SaveDir      string
	BingUrl      string
	IntervalTime int
}

const (
	config_file = "config.json"
)

func loadConfig() Settings {
	data, err := ioutil.ReadFile(config_file)
	if err != nil {
		fmt.Println("无法读取配置文件[" + config_file + "]!")
		os.Exit(-1)
	}
	var setting Settings
	json.Unmarshal(data, &setting)
	return setting
}
