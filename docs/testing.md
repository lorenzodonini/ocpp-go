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

**Note**: Mock generation is also run as part of the CI pipeline, so you can check if the mocks are up-to-date by
running the tests. However, It is recommended to run the tests locally before pushing your changes.

When adding new interfaces and needing to generate mocks, you should:

1. Add the package/interface to the `.mockery.yaml` file. Be mindful of the naming conventions. The mocks should mostly
   be generated in the same package as the interface, be snake case and suffixed with `_mock.go`. Following the existing
   examples should give you a good idea of how to proceed.

2. Run the following command:
   ```sh
   mockery 
   ```

