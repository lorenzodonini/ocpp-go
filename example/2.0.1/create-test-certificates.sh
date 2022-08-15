#!/bin/bash

mkdir -p certs/csms
mkdir -p certs/chargingstation
cd certs
# Create CA
openssl req -new -x509 -nodes -sha256 -days 120 -extensions v3_ca -keyout ca.key -out ca.crt -subj "/CN=ocpp-go-example"
# Generate self-signed CSMS certificate
openssl genrsa -out csms/csms.key 4096
openssl req -new -out csms/csms.csr -key csms/csms.key -config ../openssl-csms.conf -sha256
openssl x509 -req -in csms/csms.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out csms/csms.crt -days 120 -extensions req_ext -extfile ../openssl-csms.conf -sha256
# Generate self-signed charging-station certificate
openssl genrsa -out chargingstation/charging-station.key 4096
openssl req -new -out chargingstation/charging-station.csr -key chargingstation/charging-station.key -config ../openssl-chargingstation.conf -sha256
openssl x509 -req -in chargingstation/charging-station.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out chargingstation/charging-station.crt -days 120 -extensions req_ext -extfile ../openssl-chargingstation.conf -sha256
