package noaaclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"

	"github.com/mornindew/sledgeconf2021/pkg/utils"
)

//Define the struct
type NoaaClient struct {
	client          *http.Client
	datum           Datum
	preferredMetric MeasurementUnit
	timeZone        string
	format          string
	application     string
}

//Constructor
func NewNoaaClient(datum Datum, preferredMetric string) *NoaaClient {
	//Construct the object to return
	noaaClientToReturn := &NoaaClient{timeZone: "gmt", format: "json", application: "sledgeconf", datum: datum}
	//Assumes metric
	if preferredMetric != English.String() {
		noaaClientToReturn.preferredMetric = English
	} else {
		noaaClientToReturn.preferredMetric = Metric
	}
	noaaClientToReturn.client = &http.Client{Timeout: time.Duration(3) * time.Second}
	return noaaClientToReturn
}

//Public Methods on the client
func (object *NoaaClient) RetreiveDataVariable(startDate, endDate *time.Time, dataProduct DataProduct, stationID *string) (*sledgconf_demo_proto_v1.ProductDataValues, error) {
	//Precondition
	if startDate == nil || endDate == nil || stationID == nil {
		return nil, customerrors.PreconditionError{Msg: "Empty Mandatory Values"}
	}

	url := object.constructURL(startDate, endDate, dataProduct, stationID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Handle the status codes
	if resp.StatusCode == 200 {
		//Parse it into object
		return object.parse200Response(&resp.Body), nil
	} else if resp.StatusCode == 400 {
		//They use 400 to handle when a station doesn't have those values.  Will return an empty response body
		//They should use 404
		return &sledgconf_demo_proto_v1.ProductDataValues{}, nil
	}
	//If we got here then something went wrong - for simplicity going to genericze to internal server errors

	return nil, object.parseErrorResponse(&resp.Body)

}

//Internal methods

//constructURL - internal function to build the URL.  This is NOT nil safe as it is private and we assume the public method is checking nil values
func (object *NoaaClient) constructURL(startDate, endDate *time.Time, dataProduct DataProduct, stationID *string) string {
	params := "begin_date=" + utils.ConvertTimeToyyyyMMdd(*startDate) + "&end_date=" + utils.ConvertTimeToyyyyMMdd(*endDate) + "&station=" + *stationID + "&product=" + dataProduct.String() + "&time_zone=" + object.timeZone + "&application=" + object.application + "&format=" + object.format + "&units=" + object.preferredMetric.String()
	switch dataProduct {
	case WaterLevel:
		params = params + "&datum=" + object.datum.String()
		break
	case WaterTemperature:
		params = params + "&datum=" + object.datum.String()
		break
	}
	baseUrl := "https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?"
	return baseUrl + params
}

//this function blows - but not all 400's are the same and we need to differentiate based on the message
func (object *NoaaClient) parse200Response(response *io.ReadCloser) *sledgconf_demo_proto_v1.ProductDataValues {
	//Object to return
	successObject := &sledgconf_demo_proto_v1.ProductDataValues{}

	//Have to parse out the error message
	body, err := ioutil.ReadAll(*response)
	if err != nil {
		fmt.Println("Error Reading the response!  Sending back nil " + err.Error())
		return successObject
	}
	//Marshal the data
	err = utils.MarshalDataToInterface(body, successObject)
	if err != nil {
		//If we cannot parse then we should log but we did get a 200 so return an empty
		fmt.Println("Error Parsing the body! " + err.Error())
	}

	//Since we got here it appears to be a valid 400
	return successObject
}

//generic method to handle error marshalling.   This API really only returns 200, 400, or 500.  Since we handle 400 as a 200 then everythign gets parsed as a 500
func (object *NoaaClient) parseErrorResponse(response *io.ReadCloser) error {
	//Have to parse out the error message
	body, err := ioutil.ReadAll(*response)
	if err != nil {
		//For some reason we couldn't parse so just assuming a bad format
		return customerrors.InternalServerError{Msg: "Bad Format Error"}
	}
	//Marshal the data
	errorObject := &ErrorResponse{}
	err = utils.MarshalDataToInterface(body, errorObject)
	if err != nil {
		//Assume that it is a critical 400 if we cannot parse it
		return customerrors.InternalServerError{Msg: "Error Parsing the body"}
	}
	//Since we got here it appears to be a valid 500
	return customerrors.InternalServerError{Msg: errorObject.Error.Message}
}
