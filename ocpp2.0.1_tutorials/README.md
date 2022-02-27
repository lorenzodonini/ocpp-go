This tutorial is a step-by-step guide to help you getting familiar with [OCPP 2.0.1](https://www.openchargealliance.org/protocols/ocpp-201/) architecture and the OCPP-Go library.

### Prerequisites
Install Go, VSCode, and wscat:
```
brew install go
brew install --cask visual-studio-code
npm install -g wscat
```
If you're new to Go, you MUST read the learning resources
- [Getting Started](https://go.dev/doc/tutorial/getting-started)
- [A Tour of Go](https://go.dev/tour/)

### Hello OCPP
Let's first create a Go project having OCPP-Go as a dependency:
```
cd ~                # go to Home directory
mkdir hello         # make the directories
cd hello            # change directory
go mod init hello   # Initializes the sevencms/hello/go.mod file

# Add the OCPP library dependencies
go get github.com/lorenzodonini/ocpp-go@master
```

From VSCode (`File > Open Folder`) open the directory `hello` and create a new file named `hello.go`. The IDE will suggest the installation of the official Go extension and other utility tools, accept them all. Then paste the contents of [this example](01-websocket-connection/main.go) into the file.

## Test
```
# Resolve the module dependencies
% go mod tidy

# Start the CSMS
% go run .
Starting CSMS on port 7777

# Open another tab and simulate a WebSocket client connection
% wscat -s ocpp2.0.1 -c ws://localhost:7777/cs001
Connected (press CTRL+C to quit)
>
```