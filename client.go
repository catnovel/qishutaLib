package qishutaLib

import (
	"crypto/tls"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
	"sync"
	"time"
)

type Client struct {
	Lock       sync.Mutex
	debug      bool
	retryCount int
	baseURL    string
	userAgent  string
	proxy      string
	cookie     string
	HTTPClient *resty.Client
}

func NewClient() *Client {
	client := &Client{
		Lock:       sync.Mutex{},
		retryCount: 2,
		HTTPClient: resty.New(),
		baseURL:    "https://www.qishuta.org",
	}
	client.HTTPClient.
		SetTimeout(10*time.Second).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetHeader("Host", "www.qishuta.org").
		SetHeader("Connection", "keep-alive").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	return client
}
func (client *Client) R() *APP {
	client.Lock.Lock()
	defer client.Lock.Unlock()
	client.HTTPClient.
		SetDebug(client.debug).
		SetBaseURL(client.baseURL).
		SetRetryCount(client.retryCount)
	if client.proxy != "" {
		client.HTTPClient.SetProxy(client.proxy)
	}
	newRequest := client.HTTPClient.R()
	if client.userAgent != "" {
		newRequest.SetHeader("User-Agent", client.userAgent)
	} else {
		newRequest.SetHeader("User-Agent", browser.Random())
	}
	if client.cookie != "" {
		newRequest.SetHeader("Cookie", client.cookie)
	}

	return &APP{Request: newRequest, Client: client}
}
func (client *Client) SetDebug() *Client {
	client.debug = true
	return client
}
func (client *Client) SetBaseURL(url string) *Client {
	client.baseURL = url
	return client
}
func (client *Client) SetUserAgent(userAgent string) *Client {
	client.userAgent = userAgent
	return client

}
func (client *Client) SetRetryCount(count int) *Client {
	client.retryCount = count
	return client
}
func (client *Client) SetCookie(cookie string) *Client {
	client.cookie = cookie
	return client
}
