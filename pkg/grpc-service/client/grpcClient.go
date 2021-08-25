package grpcclient

import (
	"context"
	"log"
	"time"

	customerrors "github.com/mornindew/sledgeconf2021/pkg/custom-errors"
	sledgconf_demo_proto_v1 "github.com/mornindew/sledgeconf2021/pkg/grpc-service/genProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//GrpcServiceClient - main struct that is used to call the server
type GrpcServiceClient struct {
	userConn sledgconf_demo_proto_v1.ExampleReddiyoGRPCServiceClient
	timeout  time.Duration
}

/*
ConstructClient - Constructor to return a client

Typically the constructor doesn't take any vars but this one does so that I can use the same client on different service locations.

Normally I would use an Environment Variable to pull the DNS location for the client
*/
func ConstructClient(location string) (*GrpcServiceClient, error) {
	//Precondition
	if location == "" {
		return nil, customerrors.PreconditionError{Msg: "Location Not Set"}
	}
	//Dial the server - this will open the pipe to the server to be multiplexed
	//This is where you would configure your TLS || Client Load Balancing
	connection, err := grpc.Dial(location, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Could not get the connection. " + err.Error())
		return nil, customerrors.ClientConstructionError{Msg: "Could not get a client"}
	}

	constructedClient := &GrpcServiceClient{
		userConn: sledgconf_demo_proto_v1.NewExampleReddiyoGRPCServiceClient(connection),
		timeout:  time.Second,
	}
	return constructedClient, nil
}

//GrpcServiceClient - methods

/*
ReverseArray - call into the grpcService and will reverse the array
*/
func (client *GrpcServiceClient) GetDataFromStations(arrayToReverse *[]string) (*map[string]*sledgconf_demo_proto_v1.Station, error) {
	//Precondition Check - I do a precondition check in the client to avoid making a call to the server for anything that isn't well constructed
	if arrayToReverse == nil || len(*arrayToReverse) == 0 {
		return nil, customerrors.PreconditionError{Msg: "Empty Array"}
	}

	//sets up the cancel function and the context
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	//Construct the protobuf params
	//We don't let the protobuf structs leak out of the client or the server layer
	request := &sledgconf_demo_proto_v1.GetDataFromStationsRequest{
		ArrayOfStationIDs: *arrayToReverse,
	}
	//Make the call to the server
	response, err := client.userConn.GetDataFromStations(ctx, request)
	if err != nil {
		//Handle error Cases - thse could be GRPC, Connection, or custom
		statusCode, ok := status.FromError(err)
		if !ok {
			//This happens if we cannot get the status - will throw the generic
			return nil, customerrors.InternalServerError{Msg: err.Error()}
		}
		switch statusCode.Code() {
		case codes.FailedPrecondition:
			return nil, customerrors.PreconditionError{Msg: err.Error()}
		case codes.InvalidArgument:
			return nil, customerrors.InvalidData{Msg: err.Error()}
		default:
			return nil, customerrors.InternalServerError{Msg: err.Error()}
		}
	}
	return &response.MapOfStationData, nil
}
