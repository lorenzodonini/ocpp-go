integration-tests:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit --exit-code-from integration_test

unit-tests:
	docker compose -f docker-compose.test.yaml up unit_test --abort-on-container-exit --exit-code-from unit_test
