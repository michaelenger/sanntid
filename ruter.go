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

// Get the arrival data for a specific location ID
func GetArrivalData(locationID int) ([]sanntidArrivalData, error) {
	data, err := requestArrivalData(arrivalDataUrl(locationID))
	if err != nil {
		return nil, err
	}

	return parseArrivalData(data), nil
}

// Construct the arrival data URL
func arrivalDataUrl(locationID int) string {
	return fmt.Sprintf("http://reisapi.ruter.no/stopvisit/getdepartures/%d", locationID)
}

// RequestArrivalData retrieves information about the upcoming arrivals for
// a given location based on its locationId.
func requestArrivalData(url string) ([]byte, error) {
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
