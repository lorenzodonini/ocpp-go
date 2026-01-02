package ocppj_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

// BenchmarkClientDispatcher_SendRequest benchmarks sending requests through the client dispatcher
func BenchmarkClientDispatcher_SendRequest(b *testing.B) {
	endpoint := ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queue := ocppj.NewFIFOClientQueue(1000)
	dispatcher := ocppj.NewDefaultClientDispatcher(queue, nil)
	state := ocppj.NewClientState()
	dispatcher.SetPendingRequestState(state)

	websocketClient := &MockWebsocketClient{}
	websocketClient.On("Write", mock.Anything).Return(nil)
	dispatcher.SetNetworkClient(websocketClient)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	req := newMockRequest("benchmark")
	call, err := endpoint.CreateCall(req)
	if err != nil {
		b.Fatalf("failed to create call: %v", err)
	}
	data, err := call.MarshalJSON()
	if err != nil {
		b.Fatalf("failed to marshal call: %v", err)
	}
	bundle := ocppj.RequestBundle{Call: call, Data: data}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(bundle)
		}
	})
}

// BenchmarkClientDispatcher_CompleteRequest benchmarks completing requests in the client dispatcher
func BenchmarkClientDispatcher_CompleteRequest(b *testing.B) {
	endpoint := ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queue := ocppj.NewFIFOClientQueue(1000)
	dispatcher := ocppj.NewDefaultClientDispatcher(queue, nil)
	state := ocppj.NewClientState()
	dispatcher.SetPendingRequestState(state)

	websocketClient := &MockWebsocketClient{}
	websocketClient.On("Write", mock.Anything).Return(nil)
	dispatcher.SetNetworkClient(websocketClient)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	// Pre-populate queue with requests
	requestIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		req := newMockRequest("benchmark")
		call, err := endpoint.CreateCall(req)
		if err != nil {
			b.Fatalf("failed to create call: %v", err)
		}
		data, err := call.MarshalJSON()
		if err != nil {
			b.Fatalf("failed to marshal call: %v", err)
		}
		bundle := ocppj.RequestBundle{Call: call, Data: data}
		requestIDs[i] = call.UniqueId
		_ = dispatcher.SendRequest(bundle)
	}

	// Wait for requests to be dispatched
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.CompleteRequest(requestIDs[i])
	}
}

// BenchmarkClientDispatcher_SendAndComplete benchmarks the full cycle of sending and completing requests
func BenchmarkClientDispatcher_SendAndComplete(b *testing.B) {
	endpoint := ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queue := ocppj.NewFIFOClientQueue(1000)
	dispatcher := ocppj.NewDefaultClientDispatcher(queue, nil)
	state := ocppj.NewClientState()
	dispatcher.SetPendingRequestState(state)

	websocketClient := &MockWebsocketClient{}
	websocketClient.On("Write", mock.Anything).Return(nil)
	dispatcher.SetNetworkClient(websocketClient)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := newMockRequest("benchmark")
		call, err := endpoint.CreateCall(req)
		if err != nil {
			b.Fatalf("failed to create call: %v", err)
		}
		data, err := call.MarshalJSON()
		if err != nil {
			b.Fatalf("failed to marshal call: %v", err)
		}
		bundle := ocppj.RequestBundle{Call: call, Data: data}
		_ = dispatcher.SendRequest(bundle)
		// Wait a bit for dispatch
		time.Sleep(1 * time.Millisecond)
		dispatcher.CompleteRequest(call.UniqueId)
	}
}

// BenchmarkServerDispatcher_SendRequest benchmarks sending requests through the server dispatcher
func BenchmarkServerDispatcher_SendRequest(b *testing.B) {
	endpoint := ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)
	state := ocppj.NewServerState(nil)
	dispatcher.SetPendingRequestState(state)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	req := newMockRequest("benchmark")
	call, err := endpoint.CreateCall(req)
	if err != nil {
		b.Fatalf("failed to create call: %v", err)
	}
	data, err := call.MarshalJSON()
	if err != nil {
		b.Fatalf("failed to marshal call: %v", err)
	}
	bundle := ocppj.RequestBundle{Call: call, Data: data}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(clientID, bundle)
		}
	})
}

// BenchmarkServerDispatcher_SendRequest_MultipleClients benchmarks sending requests to multiple clients
func BenchmarkServerDispatcher_SendRequest_MultipleClients(b *testing.B) {
	endpoint := ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)
	state := ocppj.NewServerState(nil)
	dispatcher.SetPendingRequestState(state)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("client%d", i)
		clientIDs[i] = clientID
		dispatcher.CreateClient(clientID)
	}

	req := newMockRequest("benchmark")
	call, err := endpoint.CreateCall(req)
	if err != nil {
		b.Fatalf("failed to create call: %v", err)
	}
	data, err := call.MarshalJSON()
	if err != nil {
		b.Fatalf("failed to marshal call: %v", err)
	}
	bundle := ocppj.RequestBundle{Call: call, Data: data}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			clientID := clientIDs[i%numClients]
			_ = dispatcher.SendRequest(clientID, bundle)
			i++
		}
	})
}

// BenchmarkServerDispatcher_CompleteRequest benchmarks completing requests in the server dispatcher
func BenchmarkServerDispatcher_CompleteRequest(b *testing.B) {
	endpoint := ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)
	state := ocppj.NewServerState(nil)
	dispatcher.SetPendingRequestState(state)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	// Pre-populate queue with requests
	requestIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		req := newMockRequest("benchmark")
		call, err := endpoint.CreateCall(req)
		if err != nil {
			b.Fatalf("failed to create call: %v", err)
		}
		data, err := call.MarshalJSON()
		if err != nil {
			b.Fatalf("failed to marshal call: %v", err)
		}
		bundle := ocppj.RequestBundle{Call: call, Data: data}
		requestIDs[i] = call.UniqueId
		_ = dispatcher.SendRequest(clientID, bundle)
	}

	// Wait for requests to be dispatched
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.CompleteRequest(clientID, requestIDs[i])
	}
}

// BenchmarkServerDispatcher_SendAndComplete benchmarks the full cycle of sending and completing requests
func BenchmarkServerDispatcher_SendAndComplete(b *testing.B) {
	endpoint := ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)
	state := ocppj.NewServerState(nil)
	dispatcher.SetPendingRequestState(state)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := newMockRequest("benchmark")
		call, err := endpoint.CreateCall(req)
		if err != nil {
			b.Fatalf("failed to create call: %v", err)
		}
		data, err := call.MarshalJSON()
		if err != nil {
			b.Fatalf("failed to marshal call: %v", err)
		}
		bundle := ocppj.RequestBundle{Call: call, Data: data}
		_ = dispatcher.SendRequest(clientID, bundle)
		// Wait a bit for dispatch
		time.Sleep(1 * time.Millisecond)
		dispatcher.CompleteRequest(clientID, call.UniqueId)
	}
}

// BenchmarkServerDispatcher_ConcurrentClients benchmarks concurrent requests from multiple clients
func BenchmarkServerDispatcher_ConcurrentClients(b *testing.B) {
	endpoint := ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)
	state := ocppj.NewServerState(nil)
	dispatcher.SetPendingRequestState(state)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()
	defer dispatcher.Stop()

	numClients := 100
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("client%d", i)
		clientIDs[i] = clientID
		dispatcher.CreateClient(clientID)
	}

	req := newMockRequest("benchmark")
	call, err := endpoint.CreateCall(req)
	if err != nil {
		b.Fatalf("failed to create call: %v", err)
	}
	data, err := call.MarshalJSON()
	if err != nil {
		b.Fatalf("failed to marshal call: %v", err)
	}
	bundle := ocppj.RequestBundle{Call: call, Data: data}

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID string) {
			defer wg.Done()
			for j := 0; j < b.N/numClients; j++ {
				_ = dispatcher.SendRequest(clientID, bundle)
			}
		}(clientIDs[i])
	}
	wg.Wait()
}
