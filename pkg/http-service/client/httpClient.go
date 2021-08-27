package httpclient

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
)

type StationDataHttpClient struct {
	serviceName string
}

//CreateClient constructs that will create the client.  It will return the http client to use
func CreateClient(serviceName string) (*StationDataHttpClient, error) {
	if serviceName == "" {
		return nil, customerrors.ClientConstructionError{Msg: "Empty Service Name"}
	}
	client := &StationDataHttpClient{serviceName: serviceName}
	return client, nil
}

//GetDataFromStations - Simple http call to get all the data from a specific station ID.  It will return the Station Struct with any data that was found
func (v *StationDataHttpClient) GetDataFromStation(stationID string, startTime, endTime *time.Time, datum string, metricPreference sledgconf_demo_proto_v1.MetricPreference) (*sledgconf_demo_proto_v1.Station, error) {

	// Preconidtion
	if stationID == "" || startTime == nil || endTime == nil || datum == "" {
		return nil, customerrors.PreconditionError{Msg: "Missing Mandatory Values"}
	}
	//Build the URL
	fullPath := "http://" + v.serviceName + "/station/" + stationID + "/" + datum
	base, err := url.Parse(fullPath)
	if err != nil {
		return nil, customerrors.InternalServerError{Msg: "Error Calling Service: " + err.Error()}
	}

	// Query params
	params := url.Values{}
	startTimeString := strconv.FormatInt(startTime.Unix(), 10)
	params.Add("startTime", startTimeString)
	endTimeString := strconv.FormatInt(endTime.Unix(), 10)
	params.Add("endTime", endTimeString)
	params.Add("preferredMetric", metricPreference.String())
	base.RawQuery = params.Encode()

	response, err := http.Get(base.String())
	if err != nil {
		return nil, customerrors.InternalServerError{Msg: "Error Calling Service: " + err.Error()}
	}
	defer response.Body.Close()
	//Handle non 200
	if response.StatusCode > 299 {
		//Switch on the big ones that this code handles
		switch response.StatusCode {
		case 400:
			return nil, customerrors.BadRequest{Msg: "Bad Request"}
		case 404:
			return nil, customerrors.NotFoundError{Msg: "Entity Not Found"}
		}
		return nil, customerrors.HTTPError{Msg: "Non-200 Error", Code: response.StatusCode}
	}

	//Marshal it over to our object
	objectToParse := make(map[string]*sledgconf_demo_proto_v1.Station, 0)
	json.NewDecoder(response.Body).Decode(&objectToParse)
	//Get the station ID that was requested
	val, ok := objectToParse[stationID]
	if !ok {
		//This is weird since we didn't get an error so we will return an empty struct
		return &sledgconf_demo_proto_v1.Station{}, nil
	}
	return val, nil
}
