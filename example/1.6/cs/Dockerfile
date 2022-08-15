############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

ENV GO111MODULE on
WORKDIR $GOPATH/src/github.com/lorenzodonini/ocpp-go
COPY . .
# Fetch dependencies.
RUN go mod download
# Build the binary.
RUN go build -ldflags="-w -s" -o /go/bin/central_system example/1.6/cs/*.go

############################
# STEP 2 build a small image
############################
FROM alpine

COPY --from=builder /go/bin/central_system /bin/central_system

# Ports on which the service may be exposed.
EXPOSE 8887
EXPOSE 443

CMD [ "central_system" ]
