package ruter

import (
	"fmt"
	"github.com/franela/goreq"
)

type sanntidMonitoredCall struct {
	ExpectedArrivalTime string
	DeparturePlatformName string
	DestinationDisplay string
}

type sanntidMonitoredVehicleJourney struct {
	DestinationName string
	MonitoredCall sanntidMonitoredCall
	PublishedLineName string
	VehicleMode int
}

type sanntidArrivalData struct {
	MonitoredVehicleJourney sanntidMonitoredVehicleJourney
}

func RequestArrivalData(locationId int) ([]sanntidArrivalData, error) {
	var data []sanntidArrivalData

	url := fmt.Sprintf("http://reisapi.ruter.no/stopvisit/getdepartures/%d", locationId)
	res, err := goreq.Request{ Uri: url }.Do()
	if err == nil {
		res.Body.FromJsonTo(&data)
	}

	return data, err
}
