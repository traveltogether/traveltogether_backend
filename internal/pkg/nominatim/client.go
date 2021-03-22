package nominatim

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ServiceUnavailable = errors.New("nominatim service is unavailable")
)

var client = http.Client{
	Transport: &customTransport{
		roundTripper: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "nominatim.openstreetmap.org",
			},
		},
	},
	Timeout: 5 * time.Second,
}

type customTransport struct {
	roundTripper http.RoundTripper
}

func (transport *customTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Add("User-Agent", "TravelTogether/1.0 (student project)")
	return transport.roundTripper.RoundTrip(request)
}

func GetAddress(lat float32, lon float32) (*OSMAddress, error) {
	url := getNextUrl()
	if url == nil {
		return nil, ServiceUnavailable
	}

	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://%s/reverse?format=jsonv2&lat=%f&lon=%f", *url, lat, lon), nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusTooManyRequests {
			banUrl(*url)
		}
		return GetAddress(lat, lon)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	address := &OSMResponse{}
	err = json.Unmarshal(body, address)
	if err != nil {
		return nil, err
	}

	return &(address.Address), nil
}
