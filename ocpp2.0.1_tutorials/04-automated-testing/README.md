## End-to-end tests
Testing OCPP will usually involve the interaction between a Charging Station and the CSMS. It's possible to test such interations using Mock Objects - see the [ocpp2.0.1_test](https://github.com/lorenzodonini/ocpp-go/tree/1c65522f8ab806fc626f04adbccdebbaed186f2a/ocpp2.0.1_test) directory as a reference. This example demonstrates how to write a simple test using [Goroutines](https://go.dev/tour/concurrency/1).

```
var server = ws.NewServer()
var csms = ocpp2.NewCSMS(nil, server)

var client = ws.NewClient()
var cs = ocpp2.NewChargingStation("cs001", nil, client)

func TestEnd2EndBasicAuth(t *testing.T) {
	go csms.Start(7778, "/{id}")
	defer server.Stop()

  /*
	server.SetBasicAuthHandler(func(user string, pass string) bool {
		return user == "cs001" && pass == "s3cr3t"
	})

	client.SetBasicAuth("cs001", "s3cr3t")
  */

	err := cs.Start("ws://localhost:7778")
	require.Nil(t, err)
}
```

The CSMS must be started in a Goroutine, otherwise it would block the testing thread execution. The `defer` command shuts down the Websocket server at the end of the test. The Charging Station connects successfully and the test passes:

```
go test
PASS
ok  	github.com/lorenzodonini/ocpp-go/ocpp2.0.1_tutorials/04-automated-testing	0.626s
```

Now, if the CSMS Basic Auth handler is enabled by uncommenting only these lines:
```
	server.SetBasicAuthHandler(func(user string, pass string) bool {
		return user == "cs001" && pass == "s3cr3t"
	})
```

then the test fails:
```
% go test
--- FAIL: TestEnd2EndBasicAuth (0.00s)
    basic_auth_test.go:28:
        	Error Trace:	basic_auth_test.go:28
        	Error:      	Expected nil, but got: ws.HttpConnectionError{Message:"websocket: bad handshake", HttpStatus:"401 Unauthorized", HttpCode:401, Details:"Unauthorized\n"}
        	Test:       	TestEnd2EndBasicAuth
FAIL
exit status 1
FAIL	github.com/lorenzodonini/ocpp-go/ocpp2.0.1_tutorials/04-automated-testing	0.572s
```

And after configuring the Charging Station Websocket client to send Basic Auth credentials:
```
  client.SetBasicAuth("cs001", "s3cr3t")
```

the test succeeds:
```
% go test
PASS
ok  	github.com/lorenzodonini/ocpp-go/ocpp2.0.1_tutorials/04-automated-testing	0.593s
```