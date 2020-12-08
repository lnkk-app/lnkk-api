package shortener

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/html/charset"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"
)

const (
	// DatastoreGeoLocation collection GEO_LOCATION
	DatastoreGeoLocation string = "GEOLOCATION"
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

	// GeoLocation records a IP's geo location
	GeoLocation struct {
		IP          string `json:"ip"`
		ISP         string `json:"isp"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		CountryName string `json:"country_name"`
		Latitude    string `json:"latitude"`
		Longitude   string `json:"longitude"`
	}
)

// LookupGeoLocation looks up the IP's geolocation
func LookupGeoLocation(ip string) (*LocationType, error) {

	url := fmt.Sprintf("http://api.geoiplookup.net/?query=%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("StatusCode: %v", resp.StatusCode))
		return nil, err
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

// CreateGeoLocation looks up the IP's geolocation if it is unknown
func CreateGeoLocation(ctx context.Context, ip string) error {

	var loc GeoLocation
	k := geoLocationKey(ip)

	if err := platform.DataStore().Get(ctx, k, &loc); err != nil {
		// assuming the location is unknown
		l, err := LookupGeoLocation(ip)
		if err != nil {
			platform.ReportError(err)
			return err
		}

		if _, err := platform.DataStore().Put(ctx, k, l.asInternal()); err != nil {
			platform.ReportError(err)
			return err
		}
	}

	return nil
}

// asInternal converts the geolocation into the internal DS struct
func (r *LocationType) asInternal() *GeoLocation {
	location := GeoLocation{
		IP:          r.IP,
		ISP:         r.Isp,
		City:        r.City,
		CountryCode: r.Countrycode,
		CountryName: r.Countryname,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
	}
	return &location
}

// Anonimize the IP to be GDPR compliant
func anonimizeIP(ip string) string {
	if strings.ContainsRune(ip, 58) {
		// IPv6
		parts := strings.Split(ip, ":")
		return fmt.Sprintf("%s:%s:%s:0000:0000:0000:0000:0000", parts[0], parts[1], parts[2])
	}
	// IPv4
	parts := strings.Split(ip, ".")
	return fmt.Sprintf("%s.%s.%s.0", parts[0], parts[1], parts[2])
}

func geoLocationKey(ip string) *datastore.Key {
	return datastore.NameKey(DatastoreGeoLocation, ip, nil)
}
