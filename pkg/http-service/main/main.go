package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	noaaclient "github.com/mornindew/sledgeconf2021/pkg/noaa-client"
	tidesandcurrents "github.com/mornindew/sledgeconf2021/pkg/tides-and-currents"
)

func main() {

	http.HandleFunc("/station/", reverseRequestHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("Error Starting Server: " + err.Error())
	}
}

func reverseRequestHandler(w http.ResponseWriter, req *http.Request) {
	//Decode the request
	fmt.Println("CALLING REQUEST")
	//Get the stationID from the URL - a framework could help here.  YUK
	urlPathSplits := strings.Split(req.URL.Path, "/")
	//Station is in the 2nd slot
	stationID := urlPathSplits[2]

	//only do this because the function takes an array
	arrayOfStationIDs := make([]string, 0)
	arrayOfStationIDs = append(arrayOfStationIDs, stationID)
	fmt.Println(stationID)
	values := req.URL.Query()
	startTimeEpoch := values.Get("startTime")
	startTimeEpochInt, err := strconv.ParseInt(startTimeEpoch, 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert the startTime to a valid time", http.StatusBadRequest)
		return
	}
	startTime := time.Unix(startTimeEpochInt, 0)

	endTimeEpoch := values.Get("endTime")
	endTimeEpochInt, err := strconv.ParseInt(endTimeEpoch, 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert the end time to a valid time", http.StatusBadRequest)
		return
	}
	endTime := time.Unix(endTimeEpochInt, 0)
	//Convert the Datum
	//Datum is in the third
	datum := urlPathSplits[3]
	datumEnum, err := noaaclient.ConvertStringDatumToEnum(datum)
	if err != nil {
		http.Error(w, "Unable to convert the end datum to the enum", http.StatusBadRequest)
		return
	}
	//Get the preferred Metric - I don't even bother checking as it defaults to metric and nils are impossible
	preferredMetric := values.Get("preferredMetric")

	stations, err := tidesandcurrents.GetStationDataAsync(arrayOfStationIDs, &startTime, &endTime, datumEnum, preferredMetric)
	if err != nil {
		switch err.(type) {
		case customerrors.PreconditionError:
			http.Error(w, "Invalid Data", http.StatusBadRequest)
		case customerrors.InvalidData:
			http.Error(w, "Invalid Data", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//Convert the struct to JSON
	js, err := json.Marshal(stations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
