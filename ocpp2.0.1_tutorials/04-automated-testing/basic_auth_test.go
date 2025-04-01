package test

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/stretchr/testify/require"
)

var server = ws.NewServer()
var csms = ocpp2.NewCSMS(nil, server)

var client = ws.NewClient()
var cs = ocpp2.NewChargingStation("cs001", nil, client)

func TestEnd2EndBasicAuth(t *testing.T) {
	go csms.Start(7778, "/{id}")
	defer server.Stop()

	server.SetBasicAuthHandler(func(user string, pass string) bool {
		return user == "cs001" && pass == "s3cr3t"
	})
	client.SetBasicAuth("cs001", "s3cr3t")

	err := cs.Start("ws://localhost:7778")
	require.Nil(t, err)
}

func TestEnd2EndFailedAuth(t *testing.T) {
	go csms.Start(7778, "/{id}")
	defer server.Stop()

	server.SetBasicAuthHandler(func(username, password string) bool { return false })

	err := cs.Start("ws://localhost:7778")
	require.NotNil(t, err)
	require.Equal(t, "websocket: bad handshake, http status: 401 Unauthorized", err.Error())
}
