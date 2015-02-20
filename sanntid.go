package main

import (
	"github.com/michaelenger/sanntid/ruter"
)

type Line struct {
	Name string
	Destination string
}

type Arrival struct {
	Line Line
	ExpectedArrivalTime string
	Platform string
}

func GetArrivals(locationId int) ([]Arrival, error) {
	var arrivals []Arrival

	data, err := ruter.RequestArrivalData(locationId)

	if err == nil {
		for i,j := 0,0; i < len(data); i,j = i+1,j+1 {
			line := Line{
				data[i].MonitoredVehicleJourney.PublishedLineName,
				data[i].MonitoredVehicleJourney.DestinationName,
			}
			arrival := Arrival{
				line,
				data[i].MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime,
				data[i].MonitoredVehicleJourney.MonitoredCall.DeparturePlatformName,
			}

			arrivals = append(arrivals, arrival)
		}
	}

	return arrivals, err
}
