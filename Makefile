test:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit

mocks:
	docker run --rm -v $(PWD):/src -w /src vektra/mockery:v3.5.3 --config .mockery.yaml

.PHONY: test mocks
