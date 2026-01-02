# OCPP_VERSION variable for examples and performance tests (default: 1.6)
# Supported versions: 1.6, 2.0.1
OCPP_VERSION ?= 1.6

# Derive directory path and service name from OCPP_VERSION
OCPP_DIR = example/$(OCPP_VERSION)
OCPP_SERVICE = $(if $(filter 1.6,$(OCPP_VERSION)),central-system,csms)

.PHONY: test example example-observability example-ocpp-201 example-ocpp-16 example-ocpp16-observability perf-tests perf-tests-ci perf-tests-ocpp16 perf-tests-ocpp201 perf-tests-ocpp21 perf-tests-ocpp16-ci perf-tests-ocpp201-ci perf-tests-ocpp21-ci integration-tests unit-tests benchmarks

integration-tests:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit --exit-code-from integration_test

unit-tests:
	docker compose -f docker-compose.test.yaml up unit_test --abort-on-container-exit --exit-code-from unit_test

# Run benchmarks for ocppj, ws, and internal packages
benchmarks:
	docker compose -f docker-compose.test.yaml up benchmarks --abort-on-container-exit --exit-code-from benchmarks

# Generic example target
# Usage: make example OCPP_VERSION=2.0.1
example:
	docker compose -f $(OCPP_DIR)/docker-compose.yml up --build

# Generic example with observability enabled
# Usage: make example-observability OCPP_VERSION=2.0.1
example-observability:
	METRICS_ENABLED=true docker compose -f $(OCPP_DIR)/docker-compose.yml -f example/docker-compose.observability.yaml up --build

# Version-specific targets (for backward compatibility)
example-ocpp-201:
	$(MAKE) example OCPP_VERSION=2.0.1

example-ocpp-16:
	$(MAKE) example OCPP_VERSION=1.6

example-ocpp16-observability:
	$(MAKE) example-observability OCPP_VERSION=1.6

# Generic performance tests target
# Usage: make perf-tests OCPP_VERSION=2.0.1
perf-tests:
	docker compose -f $(OCPP_DIR)/docker-compose.yml \
	               -f $(OCPP_DIR)/docker-compose.k6.yml \
	               -f example/docker-compose.observability.yaml up --build

# Generic CI performance tests target
# Usage: make perf-tests-ci OCPP_VERSION=2.0.1
perf-tests-ci:
	docker compose -f $(OCPP_DIR)/docker-compose.yml \
	               -f $(OCPP_DIR)/docker-compose.k6-ci.yml up $(OCPP_SERVICE) k6 --build --abort-on-container-exit

# Version-specific targets (for backward compatibility)
perf-tests-ocpp16:
	$(MAKE) perf-tests OCPP_VERSION=1.6

perf-tests-ocpp201:
	$(MAKE) perf-tests OCPP_VERSION=2.0.1

perf-tests-ocpp16-ci:
	$(MAKE) perf-tests-ci OCPP_VERSION=1.6

perf-tests-ocpp201-ci:
	$(MAKE) perf-tests-ci OCPP_VERSION=2.0.1