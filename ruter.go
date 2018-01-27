package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// sanntidDirection defines the direction of the vehicle. It is either,
// 0 (undefined (?)), 1 or 2.
type sanntidDirection int

type sanntidMonitoredCall struct {
	ExpectedArrivalTime   string
	DeparturePlatformName string
	DestinationDisplay    string
}

type sanntidMonitoredVehicleJourney struct {
	DestinationName   string
	MonitoredCall     sanntidMonitoredCall
	PublishedLineName string
	VehicleMode       int
	DirectionRef      sanntidDirection `json:",string"`
}

// ArrivalData cointains the parsed data returned from a request to
// Ruter's API.
type sanntidArrivalData struct {
	MonitoredVehicleJourney sanntidMonitoredVehicleJourney
}

type sanntidPlaceData struct {
	Name      string
	PlaceType string
	ID        int
}

// RequestArrivalData retrieves information about the upcoming arrivals for
// a given location based on its locationId.
func requestData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseArrivalData(content []byte) []sanntidArrivalData {
	var data []sanntidArrivalData

	json.Unmarshal(content, &data)

	return data
}

func parsePlaceData(content []byte) []sanntidPlaceData {
	var data []sanntidPlaceData

	json.Unmarshal(content, &data)

	return data
}

// Get the arrival data for a specific location ID
func GetArrivalData(locationID int) ([]sanntidArrivalData, error) {
	url := fmt.Sprintf("https://reisapi.ruter.no/stopvisit/getdepartures/%d", locationID)
	data, err := requestData(url)
	if err != nil {
		return nil, err
	}

	return parseArrivalData(data), nil
}

// Get a place based on a name
func GetPlace(name string) (sanntidPlaceData, error) {
	var place sanntidPlaceData

	url := fmt.Sprintf("https://reisapi.ruter.no/place/getplaces/%s", name)
	data, err := requestData(url)
	if err != nil {
		return place, err
	}

	var places []sanntidPlaceData
	json.Unmarshal(data, &places)
	if len(places) > 0 {
		for _, v := range places {
			if v.PlaceType == "Stop" {
				return v, nil
			}
		}
	}

	return place, fmt.Errorf("Unable to find place with search text: %s", name)
}
