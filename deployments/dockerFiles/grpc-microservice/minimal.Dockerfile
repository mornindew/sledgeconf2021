FROM mornindew/base-demo-protobuf:latest as builder

ENV GO111MODULE=off

RUN apk add git
RUN apk --update add ca-certificates
WORKDIR /go/src


RUN git clone https://github.com/mornindew/sledgeconf2021.git


## TODO: Paramaterize the root folder name
# Gets the necessary dependencies that are called by the code
RUN cd /go/src/sledgeconf2021/ && go get ./...
#COMPILE PARAMETERS TO TELL THE COMPLIER TO STATICALLY LINK THE RUNTIME LIBRARIES INTO THE BINARY
RUN cd /go/src/sledgeconf2021/pkg/grpc-service/main && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app/server


# Buidling the final container from the minimilist image
FROM scratch

COPY --from=builder /app/ /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
EXPOSE 50051
ENTRYPOINT ["./server"]