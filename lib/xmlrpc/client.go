// Copyright 2020, Shulhan <m.shulhan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	libhttp "github.com/shuLhan/share/lib/http"
	libnet "github.com/shuLhan/share/lib/net"
)

const (
	requestHeader = "POST %s HTTP/1.1\r\n" +
		"User-agent: lib/xmlrpc (go)\r\n" +
		"Host: %s\r\n" +
		"Content-Type: text/xml\r\n" +
		"Content-Length: %d\r\n" +
		"\r\n" +
		"%s" +
		"%s"

	defaultTimeout = 60 * time.Second
)

//
// Client for XML-RPC.
//
type Client struct {
	conn    net.Conn
	timeout time.Duration
	url     *url.URL
}

//
// NewClient create and initialize new connection to RPC server.
//
func NewClient(url *url.URL, timeout time.Duration) (client *Client, err error) {
	if url == nil {
		return nil, nil
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}

	host, ip, port := libnet.ParseIPPort(url.Host, 0)

	client = &Client{
		url:     url,
		timeout: timeout,
	}

	if url.Scheme == schemeIsHTTPS {
		var insecure bool
		if ip != nil {
			insecure = true
		}
		if port == 0 {
			host += ":443"
		}

		config := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: insecure, //nolint: gosec
		}

		client.conn, err = tls.Dial("tcp", host, config)
	} else {
		if port == 0 {
			host += ":80"
		}

		client.conn, err = net.Dial("tcp", host)
	}
	if err != nil {
		return nil, fmt.Errorf("NewClient: Dial: %w", err)
	}

	return client, nil
}

//
// Close the client connection.
//
func (cl *Client) Close() (err error) {
	cl.url = nil
	return cl.conn.Close()
}

//
// Send the RPC method with parameters to the server.
//
func (cl *Client) Send(req Request) (resp Response, err error) { //nolint: interfacer
	var buf bytes.Buffer

	xmlbin, _ := req.MarshalText()

	fmt.Fprintf(&buf, requestHeader, cl.url.Path, cl.url.Host,
		len(xml.Header)+len(xmlbin), xml.Header, xmlbin)

	_, err = cl.conn.Write(buf.Bytes())
	if err != nil {
		return resp, err
	}

	xmlbin, err = ioutil.ReadAll(cl.conn)
	if err != nil {
		return resp, err
	}

	httpRes, resBody, err := libhttp.ParseResponseHeader(xmlbin)
	if err != nil {
		return resp, fmt.Errorf("Send: %w", err)
	}

	if len(resBody) > 0 {
		err = resp.UnmarshalText(resBody)
		if err != nil {
			return resp, err
		}
	}
	if !resp.IsFault {
		if httpRes.StatusCode != 200 {
			resp.IsFault = true
		}
	}

	return resp, nil
}
