//
// client.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package storage

import (
	"sync"

	"github.com/GSokol/go-aviasales-task/pkg/aviasales/places/client"
	jsoniter "github.com/json-iterator/go"
)

var responsePool sync.Pool

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type response struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func acquireResponses() *[]response {
	v := responsePool.Get()
	if v == nil {
		return &[]response{}
	}
	return v.(*[]response)
}

func releaseResponses(resp *[]response) {
	(*resp) = (*resp)[:0]
	responsePool.Put(resp)
}

type Client struct {
	client *client.Client
}

func NewClient(c *client.Client) *Client {
	return &Client{client: c}
}

func (c *Client) Get(term, locale []byte) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	clientResps := client.AcquireResponses()
	func() {
		if err = c.client.Get(
			clientResps,
			term,
			locale,
			client.PlaceTypeCity,
			client.PlaceTypeAirport,
		); err != nil {
			return
		}
		resps := acquireResponses()
		for _, clientResp := range *clientResps {
			resp := response{
				Slug:  clientResp.Code,
				Title: clientResp.Name,
			}
			if clientResp.Type == client.PlaceTypeAirport || clientResp.CityName == "" {
				resp.Subtitle = clientResp.CountryName
			} else {
				resp.Subtitle = clientResp.CityName
			}
			(*resps) = append((*resps), resp)
		}

		data, err = json.Marshal(resps)
		releaseResponses(resps)
	}()
	client.ReleaseResponses(clientResps)
	return data, err
}
