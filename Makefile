.PHONY: test-docker test

test-docker:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit

test:
	go test -race -v -covermode=atomic -coverprofile=coverage.out ./ocppj \
    go test -race -v -covermode=atomic -coverprofile=ocpp16.out -coverpkg=github.com/lorenzodonini/ocpp-go/ocpp1.6/... github.com/lorenzodonini/ocpp-go/ocpp1.6_test \
    go test -race -v -covermode=atomic -coverprofile=ocpp201.out -coverpkg=github.com/lorenzodonini/ocpp-go/ocpp2.0.1/... github.com/lorenzodonini/ocpp-go/ocpp2.0.1_test
