package tidesandcurrents

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	noaaclient "github.com/mornindew/sledgeconf2021/pkg/noaa-client"
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

//GetTidesAndCurrentsSyncTest - Will test the syncronous getting of the tides and currents
func TestGetTidesAndCurrentsSync(t *testing.T) {

	//Convert the station list to a map
	stationIDs := make([]string, 0)
	for key := range testingConstants {
		stationIDs = append(stationIDs, key)
	}
	//Put all the stations in the map
	//Get the last 24 hours worth of data
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	//Loop through the stations around Providence

	//Time the call
	callStartTime := time.Now()
	stationData, err := GetStationDataSync(stationIDs, &startTime, &endTime, noaaclient.CRD, noaaclient.Metric.String())
	if err != nil {
		t.Error(err.Error())
	}
	callEndTime := time.Now()
	callTime := callEndTime.Sub(callStartTime)
	fmt.Println("Call Duration In Millis: " + strconv.FormatInt(callTime.Milliseconds(), 10))
	//Check the length
	if len(*stationData) != len(stationIDs) {
		t.Error("Not the correct sets of data")
	}

	//loop through the station data
	for key, station := range *stationData {
		//Get the name
		name := testingConstants[key]
		fmt.Println("Station Name: " + name)
		//Loop through the product data
		for productKey, productValue := range station.ProductData {
			fmt.Println("	ProductKey: " + productKey)

			b, err := json.Marshal(productValue)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("		" + string(b))
		}
	}

}

//GetTidesAndCurrentsSyncTest - Will test the syncronous getting of the tides and currents
func TestGetTidesAndCurrentsAsync(t *testing.T) {

	//Convert the station list to a map
	stationIDs := make([]string, 0)
	for key := range testingConstants {
		stationIDs = append(stationIDs, key)
	}
	//Put all the stations in the map
	//Get the last 24 hours worth of data
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	//Loop through the stations around Providence

	//Time the call
	callStartTime := time.Now()
	stationData, err := GetStationDataAsync(stationIDs, &startTime, &endTime, noaaclient.CRD, noaaclient.Metric.String())
	if err != nil {
		t.Error(err.Error())
	}
	callEndTime := time.Now()
	callTime := callEndTime.Sub(callStartTime)
	fmt.Println("Call Duration In Millis: " + strconv.FormatInt(callTime.Milliseconds(), 10))
	//Check the length
	if len(*stationData) != len(stationIDs) {
		t.Error("Not the correct sets of data")
	}
	//loop through the station data
	for key, station := range *stationData {
		//Get the name
		name := testingConstants[key]
		fmt.Println("Station Name: " + name)
		//Loop through the product data
		for productKey, productValue := range station.ProductData {
			fmt.Println("	ProductKey: " + productKey)

			b, err := json.Marshal(productValue)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("		" + string(b))
		}
	}
}
