package main

import (
	"context"
	"log"
	"net"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
	noaaclient "github.com/mornindew/sledgeconf2021/pkg/noaa-client"
	tidesandcurrents "github.com/mornindew/sledgeconf2021/pkg/tides-and-currents"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server is used to implement the GRPC Service
type server struct{}

func (s *server) GetDataFromStations(ctx context.Context, in *sledgconf_demo_proto_v1.GetDataFromStationsRequest) (*sledgconf_demo_proto_v1.GetDataFromStationsResponse, error) {

	//Precondition check - ensure that there are values in the station ID list
	if len(in.ArrayOfStationIDs) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "empty inputs")
	}
	//Check that enddate is after the start date
	if in.EndTimeEpochInSeconds <= in.StartTimeEpochInSeconds {
		return nil, status.Errorf(codes.InvalidArgument, "The End Date is after the start date")
	}
	//Convert the epoch times to times
	startTime := time.Unix(in.StartTimeEpochInSeconds, 0)
	endTime := time.Unix(in.EndTimeEpochInSeconds, 0)

	//Convert the datum to a valid enum
	datum, err := noaaclient.ConvertStringDatumToEnum(in.Datum)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "The Datum Is Not a valid datum")
	}
	//Get the station data
	mapOfStationData, err := tidesandcurrents.GetStationDataAsync(in.ArrayOfStationIDs, &startTime, &endTime, datum, in.MetricPreference.String())
	//Handle errors
	if err != nil {
		switch err.(type) {
		case customerrors.PreconditionError:
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		case customerrors.InvalidData:
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	//Create the response object
	response := &sledgconf_demo_proto_v1.GetDataFromStationsResponse{MapOfStationData: *mapOfStationData}

	return response, nil
}

func main() {
	//Set up the server to listen - Puke if it cannot
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Not Listening: " + err.Error())
	}
	s := grpc.NewServer()
	//Load the protobuf definition
	sledgconf_demo_proto_v1.RegisterExampleReddiyoGRPCServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Not Listening: " + err.Error())
	}
}
