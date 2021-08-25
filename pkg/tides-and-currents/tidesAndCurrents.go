package tidesandcurrents

import (
	"sync"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
	noaaclient "github.com/mornindew/sledgeconf2021/pkg/noaa-client"
)

/*
Function to get the tides and currents Syncronously

Returns a map with each station and all available data

Errors:
	PreconditionError - missing mandatory data
	InvalidData - incorrect data
	InternalServerError - unhandled error
*/
func GetStationDataSync(stationIDs []string, startDate, endDate *time.Time, datum noaaclient.Datum, preferredMetric string) (*map[string]*sledgconf_demo_proto_v1.Station, error) {
	//Precondition check
	if len(stationIDs) == 0 || startDate == nil || endDate == nil {
		return nil, customerrors.PreconditionError{Msg: "Missing Mandatory Data"}
	}
	//check that start date is after end date
	if !endDate.After(*startDate) {
		return nil, customerrors.InvalidData{Msg: "The End Date is Not After the Start Date", InternalErrorCode: 1156}
	}
	mapToReturnOfAllStations := make(map[string]*sledgconf_demo_proto_v1.Station)
	//Setup the client to call Noaa
	client := noaaclient.NewNoaaClient(datum, preferredMetric)

	//Loop through all the station IDs
	for _, val := range stationIDs {
		//Construct the station object
		stationData := &sledgconf_demo_proto_v1.Station{}
		//Construc the station data
		stationData.ProductData = make(map[string]*sledgconf_demo_proto_v1.ProductDataValues)
		//Loop through all the data and aggregate what is available
		//Sort of tricky loop but it is iterating through the enumeration list
		for productEnum := noaaclient.DataProduct(0); productEnum < noaaclient.MaximumLimit; productEnum++ {
			stationProductData, err := client.RetreiveDataVariable(startDate, endDate, productEnum, &val)
			if err != nil {
				return nil, err
			}
			//Append the product data
			stationData.ProductData[productEnum.String()] = stationProductData
		}
		//Append it to the full station map response
		mapToReturnOfAllStations[val] = stationData
	}
	return &mapToReturnOfAllStations, nil
}

/*
Function to get the tides and currents Asyncronously

Returns a map with each station and all available data

Errors:
	PreconditionError - missing mandatory data
	InvalidData - incorrect data
	InternalServerError - unhandled error
*/
func GetStationDataAsync(stationIDs []string, startDate, endDate *time.Time, datum noaaclient.Datum, preferredMetric string) (*map[string]*sledgconf_demo_proto_v1.Station, error) {
	//Precondition check
	if len(stationIDs) == 0 || startDate == nil || endDate == nil {
		return nil, customerrors.PreconditionError{Msg: "Missing Mandatory Data"}
	}
	//check that start date is after end date
	if !endDate.After(*startDate) {
		return nil, customerrors.InvalidData{Msg: "The End Date is Not After the Start Date", InternalErrorCode: 1156}
	}
	mapToReturnOfAllStations := make(map[string]*sledgconf_demo_proto_v1.Station)
	//Setup the waitgroup
	var wg sync.WaitGroup
	wg.Add(len(stationIDs))
	//Setup an err chan
	dataChan := make(chan interface{})
	//Loop through all the station IDs
	for _, val := range stationIDs {
		go func(stationID string) {
			defer wg.Done()
			//get all the data from each station
			stationData, err := getAllDataFromStationAsync(datum, preferredMetric, startDate, endDate, stationID)
			if err != nil {
				dataChan <- err
				return
			}
			//Append it to the full station map response
			dataChan <- stationData
		}(val)
	}

	//Separate go routine to "wait" and close the error chan
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	//Loop through all the errors until the err Channel Closes
	for data := range dataChan {
		switch data.(type) {
		case *sledgconf_demo_proto_v1.Station:
			//Now that we know it is a product data value then we add it to the map
			castedValue := data.(*sledgconf_demo_proto_v1.Station)
			mapToReturnOfAllStations[castedValue.StationID] = castedValue
		case error:
			err := data.(error)
			return nil, err
		default:
			//Log unhandled chan data
		}
	}
	return &mapToReturnOfAllStations, nil
}

//internal function to get all the data from a station
func getAllDataFromStationAsync(datum noaaclient.Datum, preferredMetric string, startDate, endDate *time.Time, stationID string) (*sledgconf_demo_proto_v1.Station, error) {
	//Setup the waitgroup
	var wg sync.WaitGroup
	//Setup an err chan
	dataChan := make(chan interface{})

	//Cast the max limit over to an int
	wg.Add(int(noaaclient.MaximumLimit))
	stationData := &sledgconf_demo_proto_v1.Station{StationID: stationID}
	stationData.ProductData = make(map[string]*sledgconf_demo_proto_v1.ProductDataValues)
	//Sort of tricky loop but it is iterating through the enumeration list
	for productEnum := noaaclient.DataProduct(0); productEnum < noaaclient.MaximumLimit; productEnum++ {
		go func(goRoutineProductEnum noaaclient.DataProduct) {
			//Important defer wg.done - to tell the wg when done
			defer wg.Done()
			client := noaaclient.NewNoaaClient(datum, preferredMetric)
			stationProductData, err := client.RetreiveDataVariable(startDate, endDate, goRoutineProductEnum, &stationID)
			if err != nil {
				dataChan <- err
				//I don't have to always return but in this case I will
				return
			}
			//Set the enum - since it doesn't come from teh webserivce and the chan will need it
			stationProductData.DataType = goRoutineProductEnum.ConvertToGrpcEnum()
			//We send the data through the data chan to avoid concurrency issues with map writes
			dataChan <- stationProductData
		}(productEnum)
	}

	//Separate go routine to "wait" and close teh error chan
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	//Loop through all the errors until the err Channel Closes
	for data := range dataChan {
		//I am switching on the type but also could switch on the error
		switch data.(type) {
		case *sledgconf_demo_proto_v1.ProductDataValues:
			//Now that we know it is a product data value then we add it to the map
			castedValue := data.(*sledgconf_demo_proto_v1.ProductDataValues)
			stationData.ProductData[castedValue.DataType.String()] = castedValue
		case error:
			err := data.(error)
			return nil, err
		default:
			//Log unhandled chan data
		}
	}
	return stationData, nil
}
