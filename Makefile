test:
	docker compose -f docker-compose.test.yaml up toxiproxy integration_test --abort-on-container-exit
