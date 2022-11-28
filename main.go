package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	appVersion = "v0.0.1"
)

var (
	configDir  = ""
	configPath = ""
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func deleteConfig() {
	os.RemoveAll(configDir)
}

func createConfig(data []byte) error {
	if !PathExists(configDir) {
		err := os.Mkdir(configDir, 0777)
		if err != nil {
			return fmt.Errorf("创建配置文件目录失败")
		}
	}

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("创建配置文件失败")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("写入配置文件失败")
	}

	return nil
}

func readConfig() (*Config, error) {
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

func main() {
	// 初始化目录等参数
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("获取用户目录失败")
		os.Exit(1)
	}
	configDir = filepath.Join(userDir, ".flomocli")
	configPath = filepath.Join(configDir, "config.json")

	app := &cli.App{
		Name:    "flomo",
		Version: appVersion,
	}
	app.Commands = []*cli.Command{
		{
			Name:  "token",
			Usage: "set token eg:flomo token xxxxx",
			Action: func(c *cli.Context) error {

				token := c.Args().First()
				// 设置新的config
				_, err := NewWithToken(token)
				if err != nil {
					return err
				}
				conf := Config{}
				conf.Token = token

				// 删除之前的全部config
				deleteConfig()

				// 创建config
				data, err := json.Marshal(conf)
				if err != nil {
					return err
				}
				if err := createConfig(data); err != nil {
					return err
				}
				fmt.Println("set token success")
				return nil
			},
		},
		{
			Name:  "login",
			Usage: "login with email or phone and password",
			Action: func(c *cli.Context) error {
				email := c.Args().First()
				password := c.Args().Get(1)
				if len(email) == 0 || len(password) == 0 {
					fmt.Println("email phone password cannot be empty")
				}
				fmt.Println("TODO")
				return nil
			},
		},
		{
			Name:  "new",
			Usage: "new a mono",
			Action: func(c *cli.Context) error {
				conf, err := readConfig()
				if err != nil {
					return err
				}
				client, err := NewWithToken(conf.Token)
				if err != nil {
					return err
				}
				content := c.Args().First()
				err = client.InsertMono(content)
				if err != nil {
					return err
				}
				fmt.Println("success!")
				return nil
			},
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
