package main

import (
	"fmt"
	"os"
	"strconv"
)

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
