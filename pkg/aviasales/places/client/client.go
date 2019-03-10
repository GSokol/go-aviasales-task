//
// client.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package client

import (
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Client struct {
	transport   *fasthttp.Client
	uriTemplate string
	timeout     time.Duration
}

var fastJson = jsoniter.ConfigCompatibleWithStandardLibrary

func NewClient(opt ...OptionFunc) *Client {
	options := defaultOptions()

	for _, o := range opt {
		o(&options)
	}

	return newClient(options)
}

func newClient(options Options) *Client {
	uriTemplate := fmt.Sprintf("%s?term=%%s&locale=%%s%%s", options.host)
	return &Client{
		transport:   &fasthttp.Client{},
		uriTemplate: uriTemplate,
		timeout:     options.timeout,
	}
}

func (c *Client) Get(resp *[]Response, term, locale []byte, types ...placeType) error {
	rawReq := fasthttp.AcquireRequest()
	rawReq.URI().Update(fmt.Sprintf(
		c.uriTemplate,
		term,
		locale,
		buildTypes(types),
	))

	rawResp := fasthttp.AcquireResponse()

	if err := c.transport.DoTimeout(rawReq, rawResp, c.timeout); err != nil {
		return err
	}

	if statusCode := rawResp.StatusCode(); statusCode != fasthttp.StatusOK {
		return fmt.Errorf("request error: %d", statusCode)
	}

	fastJson.Unmarshal(rawResp.Body(), resp)

	fasthttp.ReleaseRequest(rawReq)
	fasthttp.ReleaseResponse(rawResp)
	return nil
}

func buildTypes(types []placeType) string {
	parts := make([]string, len(types))
	for i, part := range types {
		parts[i] = string(part)
	}
	return "&type[]=" + strings.Join(parts, "&type[]=")
}
