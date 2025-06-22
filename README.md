# ocpp-go

[![Build Status](https://github.com/xBlaz3kx/ocpp-go/actions/workflows/test.yaml/badge.svg)](https://github.com/xBlaz3kx/ocpp-go/actions/workflows/test.yaml)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/xBlaz3kx/ocpp-go)
[![Coverage Status](https://coveralls.io/repos/github/xBlaz3kx/ocpp-go/badge.svg?branch=master)](https://coveralls.io/github/xBlaz3kx/ocpp-go?branch=master)
[![Go report](https://goreportcard.com/badge/github.com/xBlaz3kx/ocpp-go)](https://goreportcard.com/report/github.com/xBlaz3kx/ocpp-go)

Open Charge Point Protocol implementation in Go.

The library targets modern charge points and central systems, running OCPP version 1.6+.

Given that SOAP will no longer be supported in future versions of OCPP, only OCPP-J is supported in this library.
There are currently no plans of supporting OCPP-S.

> [!NOTE]  
> This library is not affiliated with the Open Charge Alliance (OCA) in any way.

> [!NOTE]  
> This library is a fork of the original [lorenzodonini/ocpp-go](https://github.com/lorenzodonini/ocpp-go) repository.
> There might be some deviations in the API, as the main goal of this fork is to provide a more actively
> maintained version of the library with additional features and improvements, which may not result in stable releases

## Installation

Go version 1.22+ is required.

```sh
go get github.com/xBlaz3kx/ocpp-go@latest
```

You will also need to fetch some dependencies:

```sh
cd <path-to-ocpp-go>
export GO111MODULE=on
go mod download
```

Your application may either act as a **Central System** (server) or as a **Charge Point** (client).

## Features and supported versions

**Note: Releases 0.10.0 introduced breaking changes in some API, due to refactoring. The functionality remains the same,
but naming changed.**

## Roadmap

Planned milestones and features:

- [ ] OCPP 2.1 support (in progress)
- [ ] OCPP 2.1 variable management
- [ ] OCPP 2.0.1 variable management

### Supported versions

-   [x] OCPP 1.6 (documentation available [here](docs/ocpp-1.6.md))
-   [x] OCPP 1.6 Security extension (documentation available [here](docs/ocpp1.6-security-extension.md))
-   [x] OCPP 2.0.1 (examples working, but will need more real-world testing) (documentation
    available [here](docs/ocpp-2.0.1.md))
-   [x] Dedicated package for configuration management
-   [x] OCPP 2.1 (beta) (documentation available [here](docs/ocpp-2.1.md))

### Features

The library offers several advanced features, especially at `websocket` and `ocpp-j` level.

- [x] Automatic message validation
- [x] Verbose logging
- [x] Websocket ping-pong
- [x] Extensive OCPP 1.6 configuration management

#### Automatic message validation

All incoming and outgoing messages are validated by default, using the [validator](gopkg.in/go-playground/validator)
package. Constraints are defined on every request/response struct, as per OCPP specs.

Validation may be disabled at a package level if needed:

```go
ocppj.SetMessageValidation(false)
```

Use at your own risk, as this will disable validation for all messages!

> I will be evaluating the possibility to selectively disable validation for a specific message,
> e.g. by passing message options.

#### Verbose logging

The `ws` and `ocppj` packages offer the possibility to enable verbose logs, via your logger of choice, e.g.:

```go
// Setup your own logger
log = logrus.New()
log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
log.SetLevel(logrus.DebugLevel) // Debug level needed to see all logs
// Pass own logger to ws and ocppj packages
	server := ws.NewServer(ws.WithServerLogger(log.WithField("logger", "websocket")))
	// ... use server
	client := ocppj.NewClient("id", wsClient, nil, nil, log.WithField("logger", "ocppj"))
```

The logger you pass needs to conform to the `logging.Logger` interface.
Commonly used logging libraries, such as zap or logrus, adhere to this interface out-of-the-box.

If you are using a logger, that isn't conform, you can simply write an adapter between the `Logger` interface and your
own logging system.

### Websockets

#### Ping and pong messages

The websocket package supports configuring ping pong for both endpoints.

By default, the client sends a ping every 54 seconds and waits for a pong for 60 seconds, before timing out.
The values can be configured as follows:

```go
cfg := ws.NewClientTimeoutConfig()
cfg.PingPeriod = 10 * time.Second
cfg.PongWait = 20 * time.Second
websocketClient.SetTimeoutConfig(cfg)
```

By default, the server does not send out any pings and waits for a ping from the client for 60 seconds, before timing
out.
To configure the server to send out pings, the `PingPeriod` and `PongWait` must be set to a value greater than 0:

```go
cfg := ws.NewServerTimeoutConfig()
cfg.PingPeriod = 10 * time.Second
cfg.PongWait = 20 * time.Second
websocketServer.SetTimeoutConfig(cfg)
```

To disable sending ping messages, set the `PingPeriod` value to `0`.

#### Websocket compression

You can enable websocket compression on both the client and server side.
To enable compression on the client side, use the following code:

```go
websocketClient := ws.NewClient(
ws.WithClientCompression(true),
)

```

To enable compression on the server side, use the following code:

```go
websocketServer := ws.NewServer(
ws.WithCompression(true),
)

```

## Contributing

Contributions are welcome! Please refer to the [testing](docs/testing.md) guide for instructions on how to run the
tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
