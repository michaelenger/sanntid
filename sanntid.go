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
	VehicleMode int
	Direction   sanntidDirection
}

type Arrival struct {
	Line                Line
	ExpectedArrivalTime string
	Platform            string
}

func getArrivals(locationID int, direction sanntidDirection) ([]Arrival, error) {
	var arrivals []Arrival

	data, err := GetArrivalData(locationID)

	if err == nil {
		for i := 0; i < len(data); i++ {
			lineDir := data[i].MonitoredVehicleJourney.DirectionRef
			if direction == DirAny || direction == lineDir {
				line := Line{
					data[i].MonitoredVehicleJourney.PublishedLineName,
					data[i].MonitoredVehicleJourney.DestinationName,
					data[i].MonitoredVehicleJourney.VehicleMode,
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

func vehicleType(mode int) string {
	switch mode {
	case 0: // bus
		return "🚌"
	case 2: // train
		return "🚄"
	case 3: // tram
		return "🚋"
	case 4: // metro
		return "🚈"
	default:
		return "❓"
	}
}

func showArrivals(locationID int) {
	arrivals, err := getArrivals(locationID, DirAny)

	if err == nil {
		for i := 0; i < len(arrivals); i++ {
			fmt.Printf(
				"%s  %s %s - %s \n",
				vehicleType(arrivals[i].Line.VehicleMode),
				arrivals[i].Line.Name,
				arrivals[i].Line.Destination,
				arrivals[i].ExpectedArrivalTime,
			)
		}
	}
}

func main() {
	args := os.Args[1:]

	if len(args) >= 1 {
		locationID, err := strconv.ParseInt(args[0], 10, 0)
		if err == nil {
			showArrivals(int(locationID))
		} else {
			place, err := GetPlace(args[0])

			if err == nil {
				showArrivals(place.ID)
			} else {
				fmt.Printf("Error: %q\n", err)
			}
		}
	} else {
		fmt.Println("Error: Missing location ID")
	}
}
