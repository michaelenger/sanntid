package ruter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRequestData(t *testing.T) {
	exampleText := "Ruter API lol"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, exampleText)
	}))
	defer ts.Close()

	result, _ := requestData(ts.URL)
	expected := []byte(exampleText)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result == %q (got: %q)", expected, result)
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
			"FramedVehicleJourneyRef":{"DataFrameRef":"2015-02-27","DatedVehicleJourneyRef":"5803"},
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
			"MonitoredCall":{"VisitNumber":5,"VehicleAtStop":false,"DestinationDisplay":"Skøyen","AimedArrivalTime":"2015-02-27T12:29:00+01:00","ExpectedArrivalTime":"2015-02-27T12:30:25+01:00","AimedDepartureTime":"2015-02-27T12:29:00+01:00","ExpectedDepartureTime":"2015-02-27T12:30:25+01:00","DeparturePlatformName":"2"},
			"VehicleFeatureRef":null
		},
		"Extensions":{
			"IsHub":false,
			"OccupancyData":{"OccupancyAvailable":true,"OccupancyPercentage":20},
			"Deviations":[],
			"LineColour":"E60000"
		}
	}
	]`)

	result := parseArrivalData(exampleContent)
	arrivalTime, _ := time.Parse("2006-01-02T15:04:05-07:00", "2015-02-27T12:30:25+01:00")
	expected := []Arrival{
		Arrival{
			Line{
				"20",
				"Skøyen",
				Direction(2),
			},
			arrivalTime,
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result == %q (got: %q)", expected, result)
	}
}

func TestParseArrivalIncorrectData(t *testing.T) {
	exampleContent := []byte(`{
		"Explosions": "🔥"
	}
	]`)

	result := parseArrivalData(exampleContent)
	var expected []Arrival

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result == %q (got: %q)", expected, result)
	}
}

func TestParseStopData(t *testing.T) {
	exampleContent := []byte(`[
		{
			"Lines":[
				{"ID":1,"Name":"1","Transportation":8,"LineColour":"EC700C"},
				{"ID":2,"Name":"2","Transportation":8,"LineColour":"EC700C"},
				{"ID":3,"Name":"3","Transportation":8,"LineColour":"EC700C"},
				{"ID":4,"Name":"4","Transportation":8,"LineColour":"EC700C"},
				{"ID":5,"Name":"5","Transportation":8,"LineColour":"EC700C"}
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
				{"ID":11,"Name":"11","Transportation":7,"LineColour":"0B91EF"},
				{"ID":12,"Name":"12","Transportation":7,"LineColour":"0B91EF"},
				{"ID":19,"Name":"19","Transportation":7,"LineColour":"0B91EF"},
				{"ID":3902,"Name":"N2","Transportation":2,"LineColour":"E60000"},
				{"ID":3911,"Name":"N11","Transportation":2,"LineColour":"E60000"},
				{"ID":3912,"Name":"N12","Transportation":2,"LineColour":"E60000"}
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
		},
		{
			"Stops":[
				{"Lines":[],"X":595831,"Y":6644886,"Zone":"1","ShortName":"MJ","IsHub":false,"ID":3010200,"Name":"Majorstuen [T-bane]","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595929,"Y":6644771,"Zone":"1","ShortName":"MAJK","IsHub":false,"ID":3010201,"Name":"Majorstuen (i Kirkeveien)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595781,"Y":6644810,"Zone":"1","ShortName":"MAJS","IsHub":false,"ID":3010202,"Name":"Majorstuen (i Sørkedalsveien)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595891,"Y":6644844,"Zone":"1","ShortName":"MAJV","IsHub":false,"ID":3010203,"Name":"Majorstuen (i Valkyriegata)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595997,"Y":6644825,"Zone":"1","ShortName":"Majk1","IsHub":false,"ID":3010206,"Name":"Majorstuen (ved Ole Vigs gate)","District":"Oslo","DistrictID":"","PlaceType":"Stop"}
			],
			"Center":{"X":595886,"Y":6644827},
			"ID":1000022189,
			"Name":"Majorstuen (område)",
			"District":"Oslo",
			"DistrictID":null,
			"PlaceType":"Area"
		}
	]`)

	result := parseStopData(exampleContent)
	expected := []Stop{
		Stop{
			"Majorstuen [T-bane]",
			3010200,
		},
		Stop{
			"Majorstuen (i Kirkeveien)",
			3010201,
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result == %q (got: %q)", expected, result)
	}
}

func TestParseStopDataNoStops(t *testing.T) {
	exampleContent := []byte(`[
		{
			"Stops":[
				{"Lines":[],"X":595831,"Y":6644886,"Zone":"1","ShortName":"MJ","IsHub":false,"ID":3010200,"Name":"Majorstuen [T-bane]","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595929,"Y":6644771,"Zone":"1","ShortName":"MAJK","IsHub":false,"ID":3010201,"Name":"Majorstuen (i Kirkeveien)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595781,"Y":6644810,"Zone":"1","ShortName":"MAJS","IsHub":false,"ID":3010202,"Name":"Majorstuen (i Sørkedalsveien)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595891,"Y":6644844,"Zone":"1","ShortName":"MAJV","IsHub":false,"ID":3010203,"Name":"Majorstuen (i Valkyriegata)","District":"Oslo","DistrictID":"","PlaceType":"Stop"},
				{"Lines":[],"X":595997,"Y":6644825,"Zone":"1","ShortName":"Majk1","IsHub":false,"ID":3010206,"Name":"Majorstuen (ved Ole Vigs gate)","District":"Oslo","DistrictID":"","PlaceType":"Stop"}
			],
			"Center":{"X":595886,"Y":6644827},
			"ID":1000022189,
			"Name":"Majorstuen (område)",
			"District":"Oslo",
			"DistrictID":null,
			"PlaceType":"Area"
		}
	]`)

	result := parseStopData(exampleContent)

	if len(result) != 0 {
		t.Errorf("Expected len(result) == 0 (got: %d)", len(result))
	}
}

func TestParseStopDataIncorrectData(t *testing.T) {
	exampleContent := []byte(`LOLWUT`)

	result := parseStopData(exampleContent)

	if len(result) != 0 {
		t.Errorf("Expected len(result) == 0 (got: %d)", len(result))
	}
}
