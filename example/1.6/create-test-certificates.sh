#!/bin/bash

mkdir -p certs/cs
mkdir -p certs/cp
cd certs
# Create CA
openssl req -new -x509 -nodes -days 120 -extensions v3_ca -keyout ca.key -out ca.crt -subj "/CN=ocpp-go-example"
# Generate self-signed central-system certificate
openssl genrsa -out cs/central-system.key 4096
openssl req -out cs/central-system.csr -key cs/central-system.key -new -subj "/CN=central-system"
openssl x509 -req -in cs/central-system.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out cs/central-system.crt -days 120
# Generate self-signed charge-point certificate
openssl genrsa -out cp/charge-point.key 4096
openssl req -new -out cp/charge-point.csr -key cp/charge-point.key -subj "/CN=charge-point"
openssl x509 -req -in cp/charge-point.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out cp/charge-point.crt -days 120
