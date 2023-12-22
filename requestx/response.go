package requestx

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
	ctx context.Context
	err error
}

func (r *Response) ReadBody() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.Response == nil || r.Body == nil {
		return nil, errors.WithStack(errors.New("response响应为空"))
	}
	bodyBytes, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	return bodyBytes, errors.WithStack(err)
}

func (r *Response) Err() error {
	return r.err
}

func (r *Response) Context() context.Context {
	return r.ctx
}

func (r *Response) BindJson(v any) error {
	return r.Bind(v, JsonCodec)
}

// Bind 映射body到结构体
func (r *Response) Bind(v any, decoder Decoder) error {
	if r.err != nil {
		return r.err
	}
	if r.Response == nil || r.Body == nil {
		return errors.WithStack(errors.New("返回response或者body为空"))
	}

	err := decoder.Decode(r.Body, v)
	_ = r.Body.Close()

	return errors.WithStack(err)

}
