package sanntid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestArrivalDataUrl(t *testing.T) {
	expected := "http://reisapi.ruter.no/stopvisit/getdepartures/12345"
	result := arrivalDataUrl(12345)

	if expected != result {
		t.Errorf(
			"Expected URL == %q (got: %q)",
			expected,
			result)
	}
}

func TestRequestArrivalData(t *testing.T) {
	exampleText := "Ruter API lol"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, exampleText)
	}))
	defer ts.Close()

	expected := []byte(exampleText)
	result, _ := requestArrivalData(ts.URL)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf(
			"Expected result == %q (got: %q)",
			expected,
			result)
	}
}

func TestParseArrivalData(t *testing.T) {
	exampleContent := []byte(`[
	{
		"RecordedAtTime":"2015-02-27T12:29:41.618+01:00",
		"MonitoringRef":"3010536",
		"MonitoredVehicleJourney":{
			"LineRef":"20",
			"DirectionRef":"2",
			"FramedVehicleJourneyRef":{
				"DataFrameRef":
				"2015-02-27",
				"DatedVehicleJourneyRef":"5803"
			},
			"PublishedLineName":"20",
			"DirectionName":"2",
			"OperatorRef":"Unibuss",
			"OriginName":"Galgeberg (i Jordalgata)",
			"OriginRef":"3010640",
			"DestinationRef":3012501,
			"DestinationName":"Skøyen",
			"OriginAimedDepartureTime":"0001-01-01T00:00:00",
			"DestinationAimedArrivalTime":"0001-01-01T00:00:00",
			"Monitored":true,
			"InCongestion":false,
			"Delay":"PT85S",
			"TrainBlockPart":null,
			"BlockRef":"2010",
			"VehicleRef":"101047",
			"VehicleMode":0,
			"VehicleJourneyName":"20676",
			"MonitoredCall":{
				"VisitNumber":5,
				"VehicleAtStop":false,
				"DestinationDisplay":"Skøyen",
				"AimedArrivalTime":"2015-02-27T12:29:00+01:00",
				"ExpectedArrivalTime":"2015-02-27T12:30:25+01:00",
				"AimedDepartureTime":"2015-02-27T12:29:00+01:00",
				"ExpectedDepartureTime":"2015-02-27T12:30:25+01:00",
				"DeparturePlatformName":"2"
			},
			"VehicleFeatureRef":null
		},
		"Extensions":{
			"IsHub":false,
			"OccupancyData":{
				"OccupancyAvailable":true,
				"OccupancyPercentage":20
			},
			"Deviations":[],
			"LineColour":"E60000"
		}
	}
	]`)
	expected := sanntidArrivalData{
		sanntidMonitoredVehicleJourney{
			"Skøyen",
			sanntidMonitoredCall{
				"2015-02-27T12:30:25+01:00",
				"2",
				"Skøyen",
			},
			"20",
			0,
			sanntidDirection(2),
		},
	}

	result := parseArrivalData(exampleContent)[0]

	if !reflect.DeepEqual(expected, result) {
		t.Errorf(
			"Expected result == %q (got: %q)",
			expected,
			result)
	}
}
