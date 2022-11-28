package main

import (
	"fmt"
	"log"

	rest "github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const (
	defaultAg = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) flomo/1.22.112 Chrome/100.0.4896.160 Electron/18.3.15 Safari/537.36"
	authName  = "Authorization"
)

type Client struct {
	client *rest.Client
	token  string
}

func NewWithAccount(email, password string) *Client {
	c := rest.New().EnableTrace().SetHeader("user-agent", defaultAg)
	cc := &Client{
		client: c,
		token:  "",
	}
	err := cc.login(email, password)
	if err != nil {
		// 登录出错
		log.Fatalln(err)
		return nil
	}

	return cc
}

func NewWithToken(token string) (*Client, error) {
	c := rest.New().EnableTrace().SetHeader("user-agent", defaultAg)
	cc := &Client{
		client: c,
		token:  token,
	}
	err := cc.checkToken()
	if err != nil {
		// token过期错误等
		log.Fatalln(err)
		return nil, fmt.Errorf("NewWithToken() err :%v", err)
	}

	return cc, nil
}

func (c *Client) checkResp(resp string) error {
	code := gjson.Get(resp, "code")
	message := gjson.Get(resp, "message").String()
	if !code.Exists() {
		return fmt.Errorf("checkResp() err code field is nil")
	}
	if code.Int() != 0 {
		return fmt.Errorf("checkResp() err :%s", message)
	}
	return nil
}

func (c *Client) login(email, password string) error {
	url := "https://flomoapp.com/api/v1/user/login_by_email"
	bs := NewBaParm()
	bs.Add("email", email)
	bs.Add("password", password)
	bs.GenSign()
	r, err := c.client.R().SetBody(bs).Post(url)
	if err != nil || r.IsError() {
		return fmt.Errorf("login() err:%v,err:%v", err, r.Error())
	}
	if err = c.checkResp(r.String()); err != nil {
		return fmt.Errorf("login() err :%v", err)
	}
	token := gjson.Get(r.String(), "data.access_token")
	if !token.Exists() {
		// token 不存在
		return fmt.Errorf("login() err :%s", "access_token is null")
	}
	c.token = "Bearer " + token.String()

	return nil
}

func (c *Client) checkToken() error {
	url := "https://flomoapp.com/api/v1/user/me"
	bs := NewBaParm().GenSign().Encode()
	r, err := c.client.R().SetHeader(authName, c.token).SetQueryString(bs).Get(url)
	if err != nil {
		return err
	}
	if err = c.checkResp(r.String()); err != nil {
		return fmt.Errorf("checkToken() err :%v", err)
	}

	return nil
}

func (c *Client) GetAllMono() (any, error) {
	url := "https://flomoapp.com/api/v1/memo/mine/"
	bs := NewBaParm().GenSign().Encode()
	r, err := c.client.R().SetHeader(authName, c.token).SetQueryString(bs).Get(url)
	if err != nil {
		return nil, err
	}
	if err := c.checkResp(r.String()); err != nil {
		return nil, fmt.Errorf("GetAllMono() err :%v", err)
	}
	if res := gjson.Get(r.String(), "data"); res.Exists() {
		return res.Value(), nil
	}
	return nil, fmt.Errorf("GetAllMono() is nil")
}

func (c *Client) InsertMono(content string) error {
	url := "https://flomoapp.com/api/v1/memo"
	bs := NewBaParm()
	bs.Add("created_at", bs.GetTime())
	bs.Add("file_ids", [...]string{})
	bs.Add("source", "web")
	bs.Add("tz", "8:0")
	bs.Add("source", "web")
	bs.Add("content", content)
	bs.GenSign()
	r, err := c.client.R().SetHeader(authName, c.token).SetBody(bs).Put(url)
	if err != nil {
		return err
	}
	if err = c.checkResp(r.String()); err != nil {
		return fmt.Errorf("InsertMono() err:%v", err)
	}

	return nil
}
