package requestx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DEFAULT_TIMEOUT = 30 * time.Second
)

var (
	defaultBeforeFunc BeforeFunc = func(ctx context.Context, request *http.Request) (context.Context, error) {
		return ctx, nil
	}

	defaultAfterFunc AfterFunc = func(ctx context.Context, response *http.Response) (context.Context, error) {
		return ctx, nil
	}

	defaultClient = NewClient(Config{})
)

type Config struct {
	BeforeFunc BeforeFunc
	AfterFunc  AfterFunc
	ProxyUrl   *url.URL
	Logger     Logger
	RetryTimes int
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config Config) *Client {

	client := &Client{
		config: &Config{
			BeforeFunc: config.BeforeFunc,
			AfterFunc:  config.AfterFunc,
			RetryTimes: config.RetryTimes,
			Logger:     config.Logger,
			ProxyUrl:   config.ProxyUrl,
		},
	}
	if client.config.AfterFunc == nil {
		client.config.AfterFunc = defaultAfterFunc
	}

	if client.config.BeforeFunc == nil {
		client.config.BeforeFunc = defaultBeforeFunc
	}

	if client.config.RetryTimes <= 0 {
		client.config.RetryTimes = 1
	}
	if client.config.Logger == nil {
		client.config.Logger = &ZapLog{}
	}

	c := &http.Client{Timeout: DEFAULT_TIMEOUT}
	if client.config.ProxyUrl != nil {
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(client.config.ProxyUrl),
		}
	}

	client.httpClient = c
	return client
}

// Post 例子： Post("http://%s:%d//%s","localhost",8080,"xxxx" )
func (c *Client) Post(url string, args ...any) *Request {
	return c.Request(http.MethodPost, url, args...)
}

func (c *Client) Get(url string, args ...any) *Request {
	return c.Request(http.MethodGet, url, args...)
}

func (c *Client) Request(method string, url string, args ...any) *Request {
	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}

	r := &Request{
		ctx:        context.Background(),
		client:     c.httpClient,
		method:     strings.ToUpper(method),
		url:        url,
		logger:     c.config.Logger,
		retryTimes: c.config.RetryTimes,
		beforeFunc: c.config.BeforeFunc,
		afterFunc:  c.config.AfterFunc,
		headers:    http.Header{},
	}
	r.SetEncoder(JsonCodec)
	return r
}
