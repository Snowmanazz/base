package requestx

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"strings"
)

const (
	CONTENT_TYPE_JSON = "application/json;charset=utf-8"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
)

type (
	Codec interface {
		Encoder
		Decoder
	}

	Encoder interface {
		Encode(any) (io.Reader, error)
		ContentType() string
	}

	Decoder interface {
		Decode(io.Reader, any) error
	}
)

var JsonCodec = new(jsonCodec)
var FormCodec = new(formCodec)

type jsonCodec struct{}

func (c *jsonCodec) Encode(v any) (r io.Reader, err error) {
	bts, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
	return bytes.NewReader(bts), err

}

func (c *jsonCodec) Decode(r io.Reader, v any) error {
	return jsoniter.ConfigFastest.NewDecoder(r).Decode(v)
}

func (c *jsonCodec) ContentType() string {
	return CONTENT_TYPE_JSON
}

type formCodec struct{}

func (c *formCodec) Encode(v any) (r io.Reader, err error) {
	if v == nil {
		return nil, nil
	}

	switch t := v.(type) {
	case url.Values:
		return strings.NewReader(t.Encode()), nil
	case string:
		return strings.NewReader(t), nil
	default:
		return nil, errors.New(fmt.Sprintf("暂不支持该参数类型[%v]", v))
	}
}

func (c *formCodec) ContentType() string {
	return CONTENT_TYPE_FORM
}
