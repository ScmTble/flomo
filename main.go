package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
)

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
			Name:   "token",
			Usage:  "set token eg:flomo token xxxxx",
			Action: tokenCmd(),
		},
		{
			Name:   "login",
			Usage:  "login with email or phone and password",
			Action: loginCmd(),
		},
		{
			Name:   "new",
			Usage:  "new a mono",
			Action: newCmd(),
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
