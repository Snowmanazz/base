package requestx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	BeforeFunc func(context.Context, *http.Request) (context.Context, error)
	AfterFunc  func(context.Context, *http.Response) (context.Context, error)
)

type Request struct {
	Name       string
	ctx        context.Context
	client     *http.Client
	method     string
	url        string
	headers    http.Header
	body       io.Reader
	encoder    Encoder
	retryTimes int
	beforeFunc BeforeFunc
	afterFunc  AfterFunc
	logger     Logger
	debug      bool
	err        error
}

// Debug 开启dubug输出
func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) SetName(name string) *Request {
	r.Name = name
	return r
}

// SetAfterFunc 设置回调方法
func (r *Request) SetAfterFunc(a AfterFunc) *Request {
	r.afterFunc = a
	return r
}

// SetBeforeFunc 设置请求前预处理方法
func (r *Request) SetBeforeFunc(b BeforeFunc) *Request {
	r.beforeFunc = b
	return r
}

// SetContentType 设置ContentType
func (r *Request) SetContentType(s string) *Request {
	r.headers.Set("Content-Type", s)
	return r
}

// SetEncoder 设置编码器
func (r *Request) SetEncoder(e Encoder) *Request {
	r.encoder = e
	r.SetContentType(e.ContentType())
	return r
}

// SetHeader 设置请求头
func (r *Request) SetHeader(k, v string) *Request {
	r.headers.Set(k, v)
	return r
}

// SetHeaders 批量设置请求头
func (r *Request) SetHeaders(headers http.Header) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// SetContext 设置上下文
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// SetRetryTimes 设置请求重试
func (r *Request) SetRetryTimes(i int) *Request {
	if i > 0 {
		r.retryTimes = i
	}
	return r
}

// SetQuery 设置查询参数
func (r *Request) SetQuery(query any) *Request {
	uri, err := url.Parse(r.url)
	if err != nil {
		r.err = errors.WithStack(err)
		return r
	}

	switch v := query.(type) {
	case string:
		if len(v) > 0 {
			uri.RawQuery = v
		}
	case url.Values:
		uri.RawQuery = v.Encode()
	default:
		r.err = errors.WithStack(errors.New(fmt.Sprintf("暂不支持该参数类型[%v]", v)))
		return r
	}
	r.url = uri.String()
	return r

}

// Do 执行请求
func (r *Request) Do() *Response {
	resp := &Response{ctx: r.ctx}
	if r.err != nil {
		resp.err = r.err
		return resp
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, r.url, r.body)
	if err != nil {
		resp.err = errors.WithStack(err)
		return resp
	}
	req.Header = r.headers

	//请求前
	if resp.ctx, resp.err = r.beforeFunc(r.ctx, req); resp.err != nil {
		return resp
	}

	id := uuid.New().String()
	if r.debug {
		r.printCurlLog(id, req)
	}

	//发送请求
	var i int
	for i = 0; i < r.retryTimes; i++ {

		resp.Response, resp.err = r.client.Do(req)

		if resp.err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}

	success := 0
	if r.retryTimes > i {
		success = 1
	}

	r.logger.Info("%s[%s] 请求完成，请求失败%d次，成功%d次\n", r.Name, id, i, success)

	if resp.err != nil {
		resp.err = errors.WithStack(resp.err)
		return resp
	}

	//请求后
	if resp.ctx, resp.err = r.afterFunc(r.ctx, resp.Response); resp.err != nil {
		return resp
	}

	return resp
}

// SetBody 设置请求body
func (r *Request) SetBody(body any) *Request {
	if body != nil {
		r.body, r.err = r.encoder.Encode(body)
	}
	return r
}

// printCurlLog 输出请求日志
func (r *Request) printCurlLog(id string, req *http.Request) {
	var body = bytes.NewBufferString("")
	if req.Body != nil {
		_, _ = io.Copy(body, req.Body)
	}

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s[%s] request-detail：\n", r.Name, id))
	builder.WriteString(fmt.Sprintf("curl -X %s    '%s' \\\n", req.Method, req.URL.String()))

	for k, v := range req.Header {
		for _, h := range v {
			builder.WriteString(fmt.Sprintf("    --header    '%s:%s'", k, h))
			if req.Body != nil {
				builder.WriteString("\\\n")
			} else {
				builder.WriteString("\n")
			}
		}
	}

	builder.WriteString("    --data-raw  ")
	builder.WriteString("'")
	if body.Len() < 128*1024 {
		s := strings.TrimSuffix(body.String(), "\n")
		s = strings.Replace(s, "'", "\\'", -1)
		builder.WriteString(s)
	}
	builder.WriteString("'\n")

	r.logger.Debug(builder.String())

	req.Body = io.NopCloser(body)
}
