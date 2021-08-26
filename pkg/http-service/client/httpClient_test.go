package httpclient

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
)

//Created a map (for human readability)
var testingConstants = map[string]string{
	"8454000": "Providence",
	"8452944": "Conimicut Light",
	"8453662": "Providence Visibility",
	"8452951": "Potter Cove",
	"8447412": "Fall River Visibility",
	"8447387": "Borden Flats Light at Fall River",
	"8447386": "Fall River",
	"8452314": "Sandy Point Visibility",
	"8454049": "Quonset Point",
}

//TestGetStationDataOverHTTP - simple client function that will call the httpService and reverse the arrays
func TestGetStationDataOverHTTP(t *testing.T) {

	serviceName := "localhost:8888"
	client, err := CreateClient(serviceName)
	if err != nil {
		t.Error("error creating a client")
		return
	}

	//Convert the station list to a map
	stationIDs := make([]string, 0)
	for key := range testingConstants {
		stationIDs = append(stationIDs, key)
	}

	//Put all the stations in the map
	//Get the last 24 hours worth of data
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	//Create a Map to put the data in
	mapOfStations := make(map[string]*sledgconf_demo_proto_v1.Station)
	//Time the call
	callStartTime := time.Now()
	//Loop through the stations
	for _, statID := range stationIDs {
		//stationData, err := GetStationDataSync(stationIDs, &startTime, &endTime, noaaclient.CRD, noaaclient.Metric.String())
		stationData, err := client.GetDataFromStation(statID, &startTime, &endTime, "CRD", sledgconf_demo_proto_v1.MetricPreference_English)
		if err != nil {
			t.Error(err.Error())
			return
		}
		mapOfStations[statID] = stationData
	}
	callEndTime := time.Now()
	callTime := callEndTime.Sub(callStartTime)
	fmt.Println("Call Duration In Millis: " + strconv.FormatInt(callTime.Milliseconds(), 10))

	if len(mapOfStations) == 0 {
		t.Error("Empty Map")
	}
}
