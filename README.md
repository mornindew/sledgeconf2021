# Tides and Currents Working Example

This is a basic sample repo shows how we use Golang and how we build our microservices at Reddiyo.   The intended goal is to walk though how Go has worked for us (and not worked) and then provide a working example that can be built into a scratch docker container.  It isn't production ready but is a decent starting point for GRPC (and HTTP) microservices.

This code is built to pull and aggregate all the data from NOAA weather stations on the coast.  Specifically all the tides and currents data.

![NOAA](docs/images/tidesandCurrents.jpeg)

<https://tidesandcurrents.noaa.gov/>

<https://api.tidesandcurrents.noaa.gov/api/prod/>

The client will allow you to put in one or many station IDs and retrieve all the available data for a time window.

Slides from the presentation are stored in github in the in the docs/presentation.pdf

## Technology Used

* **Golang**
* **GRPC**
* **Protobuf**
* **Docker** (and docker-compose)

## Tech Assumptions

* [docker](https://docs.docker.com/v17.12/install/)
* [docker-compose](https://docs.docker.com/compose/install/)
* [Golang](https://golang.org/doc/install) - if you want to customize
* [GRPC](https://grpc.io/docs/quickstart/go/) - if you want to customize

## Repository Structure

```project
|   README.md
|   LICENSE
|   install.sh
|
└─── api - folder that stores API Definitions.  In this case it is the proto file for GRPC Microservices
|
└─── deployments - all the K8s files, docker files, or terraform files needed to build and deploy into GKE
|   |
|   └─── dockerFiles - each of the docker files and the docker compose files  
|
└─── docs - any supporting docs (e.g. Images)
|
└─── pkg - golang public packages
|   |
|   |─── grpc-service - the grpc microservice
|   |   └─── client - the client code that is used to make GRPC requests 
|   |   └─── genProto - any generated Protobuf files - these are used to store data and pass data from client to server
|   |   └─── main - the main server code that runs the GRPC endpoints
|   |
|   |─── http-service - the http microservice
|   |   └─── client - the client code
|   |   └─── main - the main server code
|   |
|   |─── noaa-client - a client package to talking to the Noaa Servers and pulling station data
|   |
|   └─── station - the package that handles knowing how to request data from Noaa, validate data, and concatonate the data
|   |
|   └─── utils - just a basic utilities package to be used across all code
|   |
|   └─── custom-errors - package that allows sharing of all custom errors
```

## Build Notes and Commands

### Compile Protobuf

Below is the command that is used to compile proto file into GoLang.  This will need to be run from the root.

```protoc -I api/v1/proto  api/v1/proto/demo.proto --go_out=plugins=grpc:pkg/grpc-service/genProto```

### Curl Docker Image

If you are running the HTTP docker container and want to curl it then use this command

```curl -X GET -H "Content-type: application/json" 'http://localhost:8888/station/8452314/CRD?endTime=1629937365&preferredMetric=English&startTime=1629850965'```

### Docker Build

You can build your own docker files from the source.   It is easiest to use docker-compose.  You can use the docker-compose.yml to set your params and then pass into the docker file.  Docker Files are located at "deployments/dockerFiles".  All docker builds are "from scratch" and should only be about 10MB

Note: If you fork the repository you need to update the docker file to pull from your repository.

Command to Build and Run
```docker-compose up```

### Sample Call To NOAA API

Below is an example call that you can drop into your browser to see what some of the data looks like.

<https://api.tidesandcurrents.noaa.gov/api/prod/datagetter?begin_date=20210808%2015:00&end_date=20210820%2015:06&station=8454000&product=water_temperature&units=english&time_zone=gmt&application=ports_screen&format=json>
