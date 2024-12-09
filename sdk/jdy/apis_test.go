package main

import (
	"fmt"
	"github.com/snowmanazz/base/requestx"
	"testing"
)

func TestJdyClient_DataAdd(t *testing.T) {
	type fields struct {
		host       string
		token      string
		secret     string
		retryTimes int
		debug      bool
		req        *requestx.Client
	}
	type args struct {
		param DataAddParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{},
			args:    args{param: DataAddParam{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewJdyClient()
			if err := j.DataAdd(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("DataAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJdyClient_DataDel(t *testing.T) {
	x()
}

func x(args ...string) {
	fmt.Println(len(args))
}
