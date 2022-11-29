package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const appVersion = "v0.0.1"

var (
	configDir  = ""
	configPath = "" // configPath = configDir + "config.json"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func deleteConfig() {
	os.RemoveAll(configDir)
}

func createConfig(data []byte) error {
	if !PathExists(configDir) {
		// 配置目录不存在,则创建
		err := os.Mkdir(configDir, 0777)
		if err != nil {
			return fmt.Errorf("创建配置文件目录失败")
		}
	}

	// 创建并打开配置文件
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("创建配置文件失败")
	}
	defer file.Close()

	// data写入配置文件
	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("写入配置文件失败")
	}

	return nil
}

func readConfig() (*Config, error) {
	// 读取配置文件，反序列化到Configs
	var config Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败")
	}

	json.Unmarshal(data, &config)
	return &config, nil
}

type Config struct {
	Token string `json:"token"`
}
