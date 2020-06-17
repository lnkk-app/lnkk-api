package shortener

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"golang.org/x/net/html/charset"
)

type (
	// LookupResult is the struct returned by a lookup on geoiplookup.net
	LookupResult struct {
		XMLName xml.Name    `xml:"ip"`
		Text    string      `xml:",chardata"`
		Results ResultsType `xml:"results"`
	}

	// ResultsType container for the location
	ResultsType struct {
		XMLName xml.Name     `xml:"results"`
		Text    string       `xml:",chardata"`
		Result  LocationType `xml:"result"`
	}

	// LocationType holds the geo data
	LocationType struct {
		Text        string `xml:",chardata"`
		IP          string `xml:"ip"`
		Host        string `xml:"host"`
		Isp         string `xml:"isp"`
		City        string `xml:"city"`
		Countrycode string `xml:"countrycode"`
		Countryname string `xml:"countryname"`
		Latitude    string `xml:"latitude"`
		Longitude   string `xml:"longitude"`
	}
)

// lookupGeoLocation looks up the IP's geolocation
func lookupGeoLocation(ip string) (*LocationType, error) {

	url := fmt.Sprintf("http://api.geoiplookup.net/?query=%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	var location LookupResult

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&location)

	if err != nil {
		return nil, err
	}

	return &location.Results.Result, nil
}
