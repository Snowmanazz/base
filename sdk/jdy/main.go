package main

import (
	"context"
	"fmt"
	"github.com/snowmanazz/base/requestx"
	"io"
	"net/http"
)

type JdyClient struct {
	//地址
	host string
	//授权token
	token string
	//密钥
	secret string
	//重试次数
	retryTimes int
	//是否调试
	debug bool

	req *requestx.Client
}

type Option func(*JdyClient)

func WithDefault(host string, token string, secret string) Option {
	return func(client *JdyClient) {
		client.host = host
		client.token = token
		client.secret = secret
	}
}

func WithRetryTimes(retryTimes int) Option {

	return func(client *JdyClient) { client.retryTimes = retryTimes }
}
func WithDebug(debug bool) Option {
	return func(client *JdyClient) { client.debug = debug }
}
func WithSecret(secret string) Option {
	return func(client *JdyClient) { client.secret = secret }
}
func WithHost(host string) Option {
	return func(client *JdyClient) { client.host = host }
}
func WithToken(token string) Option {
	return func(client *JdyClient) { client.token = token }
}

func NewJdyClient(opts ...Option) *JdyClient {
	c := new(JdyClient)

	for _, opt := range opts {
		opt(c)
	}
	if c.retryTimes <= 0 {
		c.retryTimes = 3
	}
	if c.host == "" {
		c.host = JdyUrlPrefix
	}

	afterFunc := func(c context.Context, res *http.Response) (context.Context, error) {
		if res.StatusCode != 200 {
			bytes, _ := io.ReadAll(res.Body)
			return c, fmt.Errorf(string(bytes))
		}
		return c, nil
	}

	c.req = requestx.NewClient(requestx.Config{
		RetryTimes: c.retryTimes,
		AfterFunc:  afterFunc,
	})

	return c
}
