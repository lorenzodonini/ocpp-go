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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/charge_point example/cp1.6/charge_point_sim.go

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /go/bin/charge_point /bin/charge_point

CMD [ "charge_point" ]
