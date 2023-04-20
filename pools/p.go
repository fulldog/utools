package pools

import (
	"bytes"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.mercari.io/go-dnscache"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var httpPool *sync.Pool
var bufferPool *sync.Pool
var httpClient *http.Client

func init() {
	resolver, _ := dnscache.New(5*time.Second, 10*time.Second, zap.NewNop())
	rand.Seed(time.Now().UTC().UnixNano()) // You MUST run in once in your application
	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			//DialContext: (&net.Dialer{
			//	Timeout:   30 * time.Second,
			//	KeepAlive: 30 * time.Second,
			//}).DialContext,
			//https://pkg.go.dev/go.mercari.io/go-dnscache
			DialContext:           dnscache.DialFunc(resolver, nil),
			ForceAttemptHTTP2:     true,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			//https://xujiahua.github.io/posts/20200723-golang-http-reuse/
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     300,
			MaxIdleConns:        150,
			//TLSClientConfig: &tls.Config{
			//	CipherSuites: append(defaultCipherSuites[8:], defaultCipherSuites[:8]...),
			//},
			//https://www.imwzk.com/posts/2021-03-14-why-i-always-get-503-with-golang/
			//defaultCipherSuites := []uint16{0xc02f, 0xc030, 0xc02b, 0xc02c, 0xcca8, 0xcca9, 0xc013, 0xc009, 0xc014, 0xc00a, 0x009c, 0x009d, 0x002f, 0x0035, 0xc012, 0x000a}
		},
	}
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}

	httpPool = &sync.Pool{
		New: func() interface{} {
			client := resty.NewWithClient(httpClient).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
			client.JSONMarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal
			client.JSONUnmarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal
			return client
		},
	}
}

func BufferPoolsGet() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func BufferPoolsPut(x *bytes.Buffer) {
	x.Reset()
	bufferPool.Put(x)
}

func HttpPoolsGet() *resty.Client {
	return httpPool.Get().(*resty.Client)
}

func HttpPoolsPut(c *resty.Client) {
	c.Header = http.Header{}
	c.SetRetryCount(0)
	c.SetDebug(false)
	c.DisableTrace()
	httpPool.Put(c)
}
