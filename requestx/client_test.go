package requestx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPost(t *testing.T) {
	s := Server()
	defer s.Close()

	do1 := defaultClient.Post(s.URL + "/test1").Debug().Do()
	do2 := defaultClient.Post(s.URL + "/test2").SetQuery(url.Values{
		"id": []string{"test2"},
	}).Debug().Do()
	do3 := defaultClient.Post(s.URL + "/test3").SetEncoder(FormCodec).Debug().SetBody(url.Values{
		"id": []string{"test3"},
	}).Do()
	do4 := defaultClient.Post(s.URL + "/test4").Debug().SetBody(map[string]string{
		"id": "test4",
	}).Do()
	logResponse(t, do1)
	logResponse(t, do2)
	logResponse(t, do3)
	logResponse(t, do4)
}

func TestPostResponseBind(t *testing.T) {
	s := Server()
	defer s.Close()
	var err error
	mp := make(map[string]string)
	err = defaultClient.Post(s.URL + "/test1").Do().BindJson(&mp)
	assert.NoError(t, err)
	fmt.Println(mp)
	err = defaultClient.Post(s.URL + "/test2").SetQuery(url.Values{
		"id": []string{"test2"},
	}).Do().BindJson(&mp)
	assert.NoError(t, err)
	fmt.Println(mp)

	err = defaultClient.Post(s.URL + "/test3").SetEncoder(FormCodec).SetBody(url.Values{
		"id": []string{"test3"},
	}).Do().BindJson(&mp)
	assert.NoError(t, err)
	fmt.Println(mp)

	err = defaultClient.Post(s.URL + "/test4").SetBody(map[string]string{
		"id": "test4",
	}).Do().BindJson(&mp)
	assert.NoError(t, err)
	fmt.Println(mp)
}

func TestPostHeader(t *testing.T) {
	s := Server()
	defer s.Close()

	do := defaultClient.Post(s.URL+"/auth").SetHeader("Token", "x").Debug().Do()
	logResponse(t, do)
}

func TestRetry(t *testing.T) {

	s := Server()
	defer s.Close()
	mp := make(map[string]string)
	err := defaultClient.Post(s.URL + "/test1").Debug().SetName("TTT").SetRetryTimes(10).Do().BindJson(&mp)
	assert.NoError(t, err)

}

func TestAfterBeforeFunc(t *testing.T) {
	s := Server()
	defer s.Close()

	defaultClient.Post(s.URL + "/test1").SetAfterFunc(func(ctx context.Context, response *http.Response) (context.Context, error) {
		t.Logf("status:%s\n", response.Status)
		return ctx, nil
	}).SetBeforeFunc(func(ctx context.Context, request *http.Request) (context.Context, error) {
		t.Logf("url:%s\n", request.URL)
		t.Logf("headers:%s\n", request.Header)
		t.Logf("method:%s\n", request.Method)

		return ctx, nil
	}).Debug().Do()
}

func TestReuseRequest(t *testing.T) {
	s := Server()
	defer s.Close()

	r1 := defaultClient.Post(s.URL + "/test1")
	r1.Debug().SetName("test1").SetHeader("xx", "xx").Do()
	r1.Debug().SetName("test2").SetHeader("xxx", "xxx").Do()
	//可以复用
}

func TestMultiThread(t *testing.T) {
	s := Server()
	defer s.Close()
	for i := 0; i < 10; i++ {
		go func() {
			defaultClient.Post(s.URL + "/test1").Debug().Do()
		}()
	}
}

func Server() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch r.URL.Path {
			case "/test1":
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				fmt.Fprint(w, `{"id":"test1","data":"test1-data"}`)
			case "/test2":
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				id := r.URL.Query().Get("id")
				w.Write([]byte(`{"id":" ` + id + `","data":"test2-data"}`))
			case "/test3":
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				id := r.FormValue("id")
				w.Write([]byte(`{"id":" ` + id + `","data":"test2-data"}`))
			case "/test4":
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				b, _ := io.ReadAll(r.Body)
				m := map[string]string{}
				json.Unmarshal(b, &m)
				w.Write([]byte(`{"id":" ` + m["id"] + `","data":"test2-data"}`))
			case "/auth":
				if r.Header.Get("Token") == "" {
					w.WriteHeader(http.StatusNonAuthoritativeInfo)
					fmt.Fprint(w, "NOT AUTH")

				} else {
					fmt.Fprint(w, "auth success")
				}
			}
		}
	})

	return httptest.NewServer(handler)
}

func logResponse(t *testing.T, r *Response) {
	t.Logf("response status :%v", r.Status)
	t.Logf("response Headers :%v", r.Header)
	t.Logf("response Cookies :%v", r.Cookies())
	bd, _ := io.ReadAll(r.Body)
	t.Logf("response Body :%v", string(bd))

}
