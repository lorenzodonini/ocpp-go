############################
# STEP 1 build executable binary
############################
FROM golang:stretch AS builder

ENV GO111MODULE on
WORKDIR $GOPATH/src/github.com/lorenzodonini/ocpp-go
COPY . .
# Fetch dependencies.
RUN go mod tidy
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/central_system example/cs/central_system_sim.go

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /go/bin/central_system /bin/central_system

# Port on which the service will be exposed.
EXPOSE 8887

CMD [ "central_system" ]
