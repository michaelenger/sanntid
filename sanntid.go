package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

// RequestArrivalData retrieves information about the upcoming arrivals for
// a given location based on its locationId.
func requestArrivalData(locationID int) ([]sanntidArrivalData, error) {
	var data []sanntidArrivalData

	url := fmt.Sprintf("http://reisapi.ruter.no/stopvisit/getdepartures/%d", locationID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &data)

	return data, err
}

const (
	// DirAny will give you Line in any direction.
	DirAny sanntidDirection = iota

	// DirUp will give you Line in only one direction.
	DirUp

	// DirDown will give you Line in only one direction, reverse of DirUp.
	DirDown
)

type Line struct {
	Name        string
	Destination string
	Direction   sanntidDirection
}

type Arrival struct {
	Line                Line
	ExpectedArrivalTime string
	Platform            string
}

func GetArrivals(locationId int, direction sanntidDirection) ([]Arrival, error) {
	var arrivals []Arrival

	data, err := requestArrivalData(locationId)

	if err == nil {
		for i := 0; i < len(data); i++ {
			lineDir := data[i].MonitoredVehicleJourney.DirectionRef
			if direction == DirAny || direction == lineDir {
				line := Line{
					data[i].MonitoredVehicleJourney.PublishedLineName,
					data[i].MonitoredVehicleJourney.DestinationName,
					lineDir,
				}
				arrival := Arrival{
					line,
					data[i].MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime,
					data[i].MonitoredVehicleJourney.MonitoredCall.DeparturePlatformName,
				}

				arrivals = append(arrivals, arrival)
			}
		}
	}

	return arrivals, err
}

func main() {
	args := os.Args[1:]

	if len(args) >= 1 {
		locationID, err := strconv.ParseInt(args[0], 10, 0)
		if err == nil {
			arrivals, err := GetArrivals(int(locationID), DirAny)

			if err == nil {
				for i := 0; i < len(arrivals); i++ {
					fmt.Printf(
						"%s %s - %s \n",
						arrivals[i].Line.Name,
						arrivals[i].Line.Destination,
						arrivals[i].ExpectedArrivalTime,
					)
				}
			}
		}
		if err != nil {
			fmt.Printf("Error: %q\n", err)
		}
	} else {
		fmt.Println("Error: Missing location ID")
	}
}
