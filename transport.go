/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : transport.go

* Purpose :

* Creation Date : 03-26-2016

* Last Modified : Sat Mar 26 18:39:28 2016

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package htest

import (
	"crypto/tls"
	// 	"fmt"
	"net"
	"net/http"
	// 	"net/url"
	"time"
)

type HTTransport struct {
	Transport http.RoundTripper
	config    *Config
}

func NewHTTransport(c *Config) *HTTransport {
	var keepalive time.Duration
	if c.Request.KeepAlive {
		keepalive = time.Duration(30 * time.Second)
	}
	// 	urlProxy, _ := url.Parse(fmt.Sprintf("%s://%s", c.Request.Scheme, c.Request.testIp))

	var tr HTTransport
	tr.Transport = &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: c.Request.SkipTls},
		DisableCompression: !c.Request.Compression,
		DisableKeepAlives:  !c.Request.KeepAlive,
		Dial: (&net.Dialer{
			KeepAlive: keepalive,
		}).Dial,
		// 		Proxy: http.ProxyURL(urlProxy),
	}
	tr.config = c

	return &tr
}

func (c *HTTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := c.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	t1 := time.Now()
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		// 		urlPath, err := url.Parse(resp.Header.Get("Location"))
		// 		if err != nil {
		// 			Logger.Error(err.Error())
		// 		}
		// 		urlProxy, _ := url.Parse(fmt.Sprintf("%s://%s", urlPath.Scheme, c.config.Request.testIp))
		// 		var keepalive time.Duration
		// 		if c.config.Request.KeepAlive {
		// 			keepalive = time.Duration(30 * time.Second)
		// 		}
		// 		c.Transport = &http.Transport{
		// 			TLSClientConfig:    &tls.Config{InsecureSkipVerify: c.config.Request.SkipTls},
		// 			DisableCompression: !c.config.Request.Compression,
		// 			DisableKeepAlives:  !c.config.Request.KeepAlive,
		// 			Dial: (&net.Dialer{
		// 				KeepAlive: keepalive,
		// 			}).Dial,
		// 			Proxy: http.ProxyURL(urlProxy),
		// 		}
		Logger.Warning("not support 301/302 verify", c.config.Request.testIp, c.config.Request.Hostname, time.Since(t1), resp.Header.Get("Location"))
	}
	return
}
