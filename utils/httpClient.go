package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	URL string

	Timeout uint
}

func NewHttpClient(url string, timeout uint) *HttpClient {
	return &HttpClient{
		URL:     url,
		Timeout: timeout,
	}
}

func (c *HttpClient) req(method string, headers map[string]string, data io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, c.URL, data)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := http.Client{Timeout: time.Second * time.Duration(c.Timeout)}
	resp, err = client.Do(req)
	return
}

func (c *HttpClient) Post(headers map[string]string, data io.Reader) (resp *http.Response, err error) {
	return c.req("POST", headers, data)
}

func (c *HttpClient) Get(headers map[string]string) (resp *http.Response, err error) {
	return c.req("GET", headers, nil)
}

func (c *HttpClient) PostStruct(headers map[string]string, obj any) (resp *http.Response, err error) {
	bz, err := json.Marshal(obj)
	if err != nil {
		return
	}

	headers0 := make(map[string]string)
	for k, v := range headers {
		headers0[k] = v
	}
	headers0["Content-Type"] = "application/json"

	return c.Post(headers0, strings.NewReader(string(bz)))
}
