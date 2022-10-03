package net

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Headers http.Header
	//Cookies    http.Cookie // 功能未实现
	Body       string
	Url        string
	IP         string
	Content    []byte
	StatusCode int
}

func (r *Response) sendRequest(url string, data string, timeout int, headers map[string]string, method string) error {
	var req *http.Request

	switch method {
	case http.MethodGet:
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		break
	case http.MethodPost:
		req, _ = http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
		break
	default:
		return errors.New(fmt.Sprintln("http.go http.sendRequest Request method not allowed"))
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := net.Dial(network, addr)
				// get remote address
				if conn != nil {
					req.RemoteAddr = conn.RemoteAddr().String()
				} else {
					req.RemoteAddr = ""
				}
				return conn, err
			},
		},
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return errors.New(fmt.Sprintf("http.go http.sendRequest request error, %v", err))
	}

	r.format(url, resp, req)
	return nil
}

// format the result
func (r *Response) format(url string, resp *http.Response, req *http.Request) {
	defer resp.Body.Close()

	r.StatusCode = resp.StatusCode
	r.Headers = resp.Header
	r.IP = strings.Split(req.RemoteAddr, ":")[0]
	body, _ := ioutil.ReadAll(resp.Body)
	r.Body = string(body)
	r.Content = body
	r.Url = url
}

// Get Send get request
func (r *Response) Get(url string, timeout int, headers map[string]string) error {
	return r.sendRequest(url, "", timeout, headers, http.MethodGet)
}

// Post Send post request
func (r *Response) Post(url string, data string, timeout int, headers map[string]string) error {
	if headers == nil {
		headers = map[string]string{}
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return r.sendRequest(url, data, timeout, headers, http.MethodPost)
}

// IsOk Check whether the StatusCode is equal to 200
func (r *Response) IsOk() bool {
	return r.StatusCode == 200
}

// IsBodyNull Check if the Body content is empty
func (r *Response) IsBodyNull() bool {
	return len(r.Body) == 0
}
