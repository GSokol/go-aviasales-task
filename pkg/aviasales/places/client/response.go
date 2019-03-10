//
// response.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package client

import (
	"encoding/json"
	"sync"
)

type placeType string

const (
	PlaceTypeCity    placeType = "city"
	PlaceTypeAirport placeType = "airport"
)

type Response struct {
	Type            placeType       `json:"type"`
	Cases           json.RawMessage `json:"caces,omitempty"`
	CityCases       json.RawMessage `json:"city_cases,omitempty"`
	CityCode        json.RawMessage `json:"city_code,omitempty"`
	CityName        string          `json:"city_name,omitempty"`
	Code            string          `json:"code"`
	Coordinates     json.RawMessage `json:"coordinates"`
	CountryCases    json.RawMessage `json:"country_cases,omitempty"`
	CountryCode     json.RawMessage `json:"country_code"`
	CountryName     string          `json:"country_name"`
	IndexStrings    json.RawMessage `json:"index_strings"`
	MainAirportName json.RawMessage `json:"main_airport_name,omitempty"`
	Name            string          `json:"name"`
	StateCode       json.RawMessage `json:"state_code,omitempty"`
	Weight          json.RawMessage `json:"weight"`
}

var responsePool sync.Pool

func AcquireResponses() *[]Response {
	v := responsePool.Get()
	if v == nil {
		return &[]Response{}
	}
	return v.(*[]Response)
}

func ReleaseResponses(resp *[]Response) {
	responsePool.Put(resp)
}
