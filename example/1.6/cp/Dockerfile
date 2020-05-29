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
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/charge_point example/1.6/cp/*.go

############################
# STEP 2 build a small image
############################
FROM alpine

COPY --from=builder /go/bin/charge_point /bin/charge_point

# Add CA certificates
# It currently throws a warning on alpine: WARNING: ca-certificates.crt does not contain exactly one certificate or CRL: skipping.
# Ignore the warning.
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && update-ca-certificates

CMD [ "charge_point" ]
