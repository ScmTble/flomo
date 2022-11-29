package main

import (
	"encoding/json"
	"fmt"

	cli "github.com/urfave/cli/v2"
)

func tokenCmd() cli.ActionFunc {
	return func(c *cli.Context) error {
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
	}
}

func loginCmd() cli.ActionFunc {
	return func(c *cli.Context) error {
		email := c.Args().First()
		password := c.Args().Get(1)
		if len(email) == 0 || len(password) == 0 {
			fmt.Println("email phone password cannot be empty")
		}
		fmt.Println("TODO")
		return nil
	}
}

func newCmd() cli.ActionFunc {
	return func(c *cli.Context) error {
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
	}
}
