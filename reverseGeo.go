package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	
)

type GeocoderResponse struct {
	Response struct {
		MetaInfo struct {
			TimeStamp string `json:"TimeStamp"`
		} `json:"MetaInfo"`
		View []struct {
			Result []struct {
				MatchLevel string `json:"MatchLevel"`
				Location   struct {
					Address struct {
						Label       string `json:"Label"`
						Country     string `json:"Country"`
						State       string `json:"State"`
						County      string `json:"Country"`
						City        string `json:"City"`
						District    string `json:"District"`
						Street      string `json:"Street"`
						HouseNumber string `json:"HouseNumber"`
						PostalCode  string `json:"PostalCode"`
					} `json:"Address"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}
type Position struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
type Geocoder struct {
	ApiKey string `json:"apiKey"`
}

func (geocoder *Geocoder) reverse(position Position) (GeocoderResponse, error) {
        endpoint, _ := url.Parse("https://reverse.geocoder.ls.hereapi.com/6.2/reversegeocode.json")
	queryParams := endpoint.Query()
	queryParams.Set("apiKey", geocoder.ApiKey)
	queryParams.Set("mode", "retrieveAddresses")
	queryParams.Set("prox", position.Latitude+","+position.Longitude)
	endpoint.RawQuery = queryParams.Encode()
	response, err := http.Get(endpoint.String())
	if err != nil {
		return GeocoderResponse{}, err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var geocoderResponse GeocoderResponse
		json.Unmarshal(data, &geocoderResponse)
		return geocoderResponse,nil
	}

}

func main() {
	latitude := flag.String("lat", "-1.28012", "Latitude")
	longitude := flag.String("lng", "36.87314", "Longitude")
	flag.Parse()
	geocoder := Geocoder{ApiKey: "SycGZkK07f9qp6_TRJDRXsuJRlpuOLIucAABWISxbxA"}
	result, err := geocoder.reverse(Position{Latitude: *latitude, Longitude: *longitude})
	if err != nil {
		fmt.Printf("the HTTP request failed with error %s\n", err)
		return
	}
	if len(result.Response.View) > 0 && len(result.Response.View[0].Result) > 0 {
		data, _ := json.Marshal(result.Response.View[0].Result[0])
		fmt.Println(string(data))
	}
}
