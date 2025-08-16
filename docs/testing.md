# Testing

## Running the tests

The tests are run in a Docker compose deployment to ensure a consistent environment across different systems.
Everything is set up for you, so you don't need to worry about dependencies or configurations. To run the tests, execute
the following command:

```sh
make test
```

## Generating mocks

For generating mocks, the `mockery` tool is used. For `mockery` installation, follow the instructions on
the [official docs](https://vektra.github.io/mockery/latest/).

**Note**: Mocks are not checked in and are instead generated on-the-fly as part of the CI pipeline. If your local tests fail, it may be due to your mocks being out-of-date. This can be fixed by re-generating the mocks. It is recommended to run the tests locally before pushing your changes.

When adding new interfaces and needing to generate mocks, you should:

1. Add the package/interface to the `.mockery.yaml` file. Be mindful of the naming conventions. The mocks should mostly
   be generated in the same package as the interface, be snake case and suffixed with `_mock.go`. Following the existing
   examples should give you a good idea of how to proceed.

2. Run the following command:
   ```sh
   mockery --config .mockery.yaml
   ```
   Alternatively, you may generate the mocks via make target:
   ```sh
   make mocks
   ```

## Toxiproxy

For testing the resilience of the library against network issues, some tests use Toxiproxy. 
If you wish to run the network tests locally (without docker compose) you need to:
- Install [toxiproxy](https://github.com/Shopify/toxiproxy) for your platform
- Start a local toxiproxy - `toxiproxy-server -port 8474 -host localhost`
- Run the tests - `go fmt ./... && go vet ./... && go test -v -count=1 -failfast ./...`
