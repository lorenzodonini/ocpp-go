.PHONY: test example-ocpp-201 example-ocpp-16 perf-tests-ocpp16 perf-tests-ocpp201 perf-tests-ocpp16-ci perf-tests-ocpp201-ci

test:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit

example-ocpp-201:
	docker compose -f example/2.0.1/docker-compose.yml up --build

example-ocpp-16:
	docker compose -f example/1.6/docker-compose.yml up --build

# Run the example with LGTM stack and observability enabled by default:
example-ocpp16-observability:
	METRICS_ENABLED=true docker compose -f example/1.6/docker-compose.yml -f example/docker-compose.observability.yaml up --build

perf-tests-ocpp16-ci:
	docker compose -f example/1.6/docker-compose.yml \
                   -f example/1.6/docker-compose.k6-ci.yml up central-system k6 --build --abort-on-container-exit

perf-tests-ocpp201-ci:
	docker compose -f example/2.0.1/docker-compose.yml \
				   -f  example/2.0.1/docker-compose.k6-ci.yml up csms k6 --build --abort-on-container-exit

perf-tests-ocpp16:
	docker compose -f example/1.6/docker-compose.yml \
				   -f example/1.6/docker-compose.k6.yml \
				   -f example/docker-compose.observability.yaml up --build

perf-tests-ocpp201:
	docker compose -f example/1.6/docker-compose.yml \
	               -f example/2.0.1/docker-compose.k6.yml \
	               -f example/docker-compose.observability.yaml up --build