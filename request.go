package qzone

import (
	"io"
	"net/http"
	"strings"
)

// RequestOptions http请求配置结构
type RequestOptions struct {
	Client   *http.Client
	Method   string
	Url      string
	Body     io.Reader
	Header   map[string]string
	RespFunc func(*http.Response) // 对http请求后返回的resp结构执行自定义函数
}

func NewRequest(options ...RequestOption) *RequestOptions {
	request := &RequestOptions{
		Client: &http.Client{},
		Method: http.MethodGet,
		Url:    "",
		Body:   nil,
		Header: map[string]string{
			"cookie":       "",
			"user-agent":   ua,
			"content-type": contentType,
		},
	}
	for _, opts := range options {
		opts.apply(request)
	}
	return request
}

type RequestOption interface {
	apply(*RequestOptions)
}

type funcRequestOption struct {
	f func(*RequestOptions)
}

func (fro *funcRequestOption) apply(do *RequestOptions) {
	fro.f(do)
}
func newFuncRequestOption(f func(*RequestOptions)) *funcRequestOption {
	return &funcRequestOption{f: f}
}

func WithClient(client *http.Client) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		if client != nil {
			o.Client = client
		}
	})
}

func WithMethod(method string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Method = strings.ToUpper(method)
	})
}

func WithUrl(url string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Url = url
	})
}

func WithBody(body io.Reader) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Body = body
	})
}

func WithHeader(header map[string]string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		for k, v := range header {
			o.Header[k] = v
		}
	})
}

func WithRespFunc(f func(*http.Response)) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.RespFunc = f
	})
}

// DialRequest 发起http请求，并返回读取的response.body数据
func DialRequest(options *RequestOptions) ([]byte, error) {
	req, err := http.NewRequest(options.Method, options.Url, options.Body)
	for k, v := range options.Header {
		req.Header.Add(k, v)
	}
	if err != nil {
		return nil, err
	}
	resp, err := options.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if options.RespFunc != nil {
		options.RespFunc(resp)
	}
	return io.ReadAll(resp.Body)
}
