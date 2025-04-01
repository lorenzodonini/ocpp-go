## Basic Authentication
The OCPP 2.0.1 "Part2/Chapter 1: Security Functional Block" describes the security requirements for the protocol, supporting three security profiles:

1. Unsecured Transport with Basic Authentication (HTTP Basic Authentication)
2. TLS with Basic Authentication
3. TLS with Client Side Certificates

The Unsecured Transport with Basic Authentication Profile should only be used in trusted networks, for instance in networks where there is a VPN between the CSMS and the Charging Station. For field operation it is highly recommended to use a security profile with TLS. Anyways, it's the easiest profile to implement and test:

### OCPP-Go implementation
Reference: https://github.com/lorenzodonini/ocpp-go/blob/b37db1468a177770f97b43d9f083033eef02bf77/ws/websocket.go#L165-L169

```
func init() {
  ...
  server = ws.NewServer()
  server.SetBasicAuthHandler(func(user string, pass string) bool {
    ok := authenticate(user, pass) // ... check for user and pass correctness
    return ok
  })

  csms = ocpp2.NewCSMS(nil, server)
  ...
}

func authenticate(user string, pass string) bool {
  return user == "cs001" && pass == "s3cr3t"
}
```

### Running the example
```
# Start CSMS
% cd 02-websockets-basic-auth
% go run .

# WebSocket client with wrong password
% wscat -s ocpp2.0.1 --auth "cs001:wrongpass" -c ws://localhost:7777/cs001
error: Unexpected server response: 401

# WebSocket client with good password
% wscat -s ocpp2.0.1 --auth "cs001:s3cr3t" -c ws://localhost:7777/cs001
Connected (press CTRL+C to quit)
>

# Server Logs
INFO[2022-02-27T20:35:57-03:00] CSMS started on port 7777
INFO[2022-02-27T20:35:57-03:00] listening on tcp network :7777						logger=websocket
ERRO[2022-02-27T20:36:02-03:00] basic auth failed: credentials invalid		logger=websocket
INFO[2022-02-27T20:36:29-03:00] Charging Station connected								client=cs001
```