package main

import (
	"github.com/snowmanazz/base/requestx"
	"net/url"
)

func (j *JdyClient) generateReq(suffix string) *requestx.Request {
	uri, _ := url.JoinPath(JdyUrlPrefix, suffix)
	return j.req.Post(uri).SetHeader("Authorization", "Bearer "+j.token)
}

// 新增数据
func (j *JdyClient) DataAdd(param DataAddParam) error {
	return j.generateReq(AddUrlSuffix).SetBody(param).Do().Err()
}

// 批量新增数据
func (j *JdyClient) DataBatchAdd(param DataBatchAddParam) (res DataBatchAddResponse, err error) {
	if len(param.DataList) <= 0 {
		return
	}

	err = j.generateReq(BatchAddUrlSuffix).SetBody(param).Do().BindJson(&res)
	return
}
