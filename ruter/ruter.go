package ruter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Internal types

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
	DirectionRef      int `json:",string"`
}

type sanntidArrivalData struct {
	MonitoredVehicleJourney sanntidMonitoredVehicleJourney
}

type sanntidStopData struct {
	Name      string
	PlaceType string
	ID        int
}

// External types

// Stop which can receive arrivals
type Stop struct {
	Name string
	ID   int
}

// Direction of the arrival
type Direction int

const (
	DirUnknown Direction = iota
	DirUp
	DirDown
)

// Public transportation line
type Line struct {
	Name        string
	Destination string
	Direction   Direction
}

// Arrival at a stop
type Arrival struct {
	Line                Line
	ExpectedArrivalTime time.Time
}

func requestData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseArrivalData(content []byte) []Arrival {
	var data []sanntidArrivalData
	var arrivals []Arrival

	timeLayout := "2006-01-02T15:04:05-07:00"

	json.Unmarshal(content, &data)
	for _, v := range data {
		arrivalTime, _ := time.Parse(timeLayout, v.MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime)

		line := Line{
			v.MonitoredVehicleJourney.PublishedLineName,
			v.MonitoredVehicleJourney.DestinationName,
			Direction(v.MonitoredVehicleJourney.DirectionRef),
		}
		arrival := Arrival{
			line,
			arrivalTime,
		}

		arrivals = append(arrivals, arrival)
	}

	return arrivals
}

func parseStopData(content []byte) []Stop {
	var data []sanntidStopData
	var stops = make([]Stop, 0)
	json.Unmarshal(content, &data)

	for _, v := range data {
		if v.PlaceType == "Stop" {
			stops = append(stops, Stop{v.Name, v.ID})
		}
	}

	return stops
}

// Get the arrival data for a specific location ID
func GetArrivals(locationID int) ([]Arrival, error) {
	url := fmt.Sprintf("https://reisapi.ruter.no/stopvisit/getdepartures/%d", locationID)
	data, err := requestData(url)
	if err != nil {
		return nil, err
	}

	return parseArrivalData(data), nil
}

// Get a stop based on its name
func GetStop(name string) (Stop, error) {
	var stop Stop

	url := fmt.Sprintf("https://reisapi.ruter.no/place/getplaces/%s", name)
	data, err := requestData(url)
	if err != nil {
		return stop, err
	}

	stops := parseStopData(data)
	if len(stops) == 0 {
		return stop, fmt.Errorf("Unable to find stop with search text: %s", name)
	}

	return stops[0], nil
}
