package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
)

func tokenCmd() cli.ActionFunc {
	return func(c *cli.Context) error {
		token := c.Args().First()
		// 设置新的config
		if _, err := NewWithToken(token); err != nil {
			return err
		}
		conf := Config{Token: token}

		// 删除之前的全部config
		deleteConfig()

		// 创建config
		var data []byte
		var err error
		if data, err = json.Marshal(conf); err != nil {
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
			return fmt.Errorf("email phone password cannot be empty")
		}
		// login
		client, err := NewWithAccount(email, password)
		if err != nil {
			return err
		}
		// 删除config
		deleteConfig()
		// 设置新的config
		conf := &Config{Token: client.token}
		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		if err := createConfig(data); err != nil {
			return err
		}
		fmt.Println("login success")
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

func vimCmd() cli.ActionFunc {
	return func(c *cli.Context) error {
		conf, err := readConfig()
		if err != nil {
			return err
		}
		client, err := NewWithToken(conf.Token)
		if err != nil {
			return err
		}

		path := filepath.Join(os.TempDir(), "tmp.txt")
		f, err := os.CreateTemp(os.TempDir(), "tmp.txt")
		if err != nil {
			return err
		}
		defer os.Remove(path)
		if err := f.Close(); err != nil {
			return err
		}

		cmd := exec.Command("vim", path, "-n")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err = client.InsertMono(string(data)); err != nil {
			return err
		}
		fmt.Println("success!")
		return nil
	}
}
