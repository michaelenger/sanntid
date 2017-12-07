package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRequestData(t *testing.T) {
	exampleText := "Ruter API lol"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, exampleText)
	}))
	defer ts.Close()

	expected := []byte(exampleText)
	result, _ := requestData(ts.URL)

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

func TestParsePlaceData(t *testing.T) {
	exampleContent := []byte(`[
	{
		"Lines":[
			{
				"ID":1,
				"Name":"1",
				"Transportation":8,
				"LineColour":"EC700C"
			},
			{
				"ID":2,
				"Name":"2",
				"Transportation":8,
				"LineColour":"EC700C"
			},
			{
				"ID":3,
				"Name":"3",
				"Transportation":8,
				"LineColour":"EC700C"
			},
			{
				"ID":4,
				"Name":"4",
				"Transportation":8,
				"LineColour":"EC700C"
			},
			{
				"ID":5,
				"Name":"5",
				"Transportation":8,
				"LineColour":"EC700C"
			}
		],
		"X":595831,
		"Y":6644886,
		"Zone":"1",
		"ShortName":"MJ",
		"IsHub":false,
		"ID":3010200,
		"Name":"Majorstuen [T-bane]",
		"District":"Oslo",
		"DistrictID":null,
		"PlaceType":"Stop"
	},
	{
		"Lines":[
			{
				"ID":11,
				"Name":"11",
				"Transportation":7,
				"LineColour":"0B91EF"
			},
			{
				"ID":12,
				"Name":"12",
				"Transportation":7,
				"LineColour":"0B91EF"
			},
			{
				"ID":19,
				"Name":"19",
				"Transportation":7,
				"LineColour":"0B91EF"
			},
			{
				"ID":3902,
				"Name":"N2",
				"Transportation":2,
				"LineColour":"E60000"
			},
			{
				"ID":3911,
				"Name":"N11",
				"Transportation":2,
				"LineColour":"E60000"
			},
			{
				"ID":3912,
				"Name":"N12",
				"Transportation":2,
				"LineColour":"E60000"
			}
		],
		"X":595929,
		"Y":6644771,
		"Zone":"1",
		"ShortName":"MAJK",
		"IsHub":false,
		"ID":3010201,
		"Name":"Majorstuen (i Kirkeveien)",
		"District":"Oslo",
		"DistrictID":null,
		"PlaceType":"Stop"
	}
	]`)
	expected := sanntidPlaceData{
		"Majorstuen [T-bane]",
		"Stop",
		3010200,
	}

	result := parsePlaceData(exampleContent)

	if !reflect.DeepEqual(expected, result[0]) {
		t.Errorf(
			"Expected result == %q (got: %q)",
			expected,
			result)
	}
}
