package ocppj

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	dispatcherQueueSize      = "dispatcher_queue_size"
	dispatcherPendingRequest = "dispatcher_pending_requests"
)

type dispatcherMetrics struct {
	meter metric.Meter

	requestQueue    metric.Int64ObservableGauge
	pendingRequests metric.Int64ObservableUpDownCounter
}

func newDispatcherMetrics(meterProvider metric.MeterProvider) (*dispatcherMetrics, error) {
	if meterProvider == nil {
		return nil, errors.New("meterProvider is nil")
	}

	meter := meterProvider.Meter("server_dispatcher")

	clientQueue, err := meter.Int64ObservableGauge(
		dispatcherQueueSize,
		metric.WithDescription("Number of messages in the dispatcher's queue"),
	)
	if err != nil {
		return nil, err
	}

	clientPendingRequest, err := meter.Int64ObservableUpDownCounter(
		dispatcherPendingRequest,
		metric.WithDescription("Number of pending requests in the dispatcher"),
	)
	if err != nil {
		return nil, err
	}

	dispatcher := &dispatcherMetrics{
		meter:           meter,
		requestQueue:    clientQueue,
		pendingRequests: clientPendingRequest,
	}

	return dispatcher, nil
}

func (d *dispatcherMetrics) ObserveInFlightRequests(state *serverState) {
	if d.meter == nil {
		return
	}

	_, err := d.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			state.mutex.RLock()
			currentState := state.pendingRequestState
			state.mutex.RUnlock()

			for clientID, clientState := range currentState {
				inFlightRequest := int64(0)
				if clientState.HasPendingRequest() {
					inFlightRequest = 1
				}
				obs.ObserveInt64(
					d.pendingRequests,
					inFlightRequest,
					metric.WithAttributes(attribute.String("client_id", clientID)),
				)
			}
			return nil
		},
		d.pendingRequests)
	if err != nil {
		log.Errorf("failed to register callback for inflight queue size: %v", err)
	}
}

func (d *dispatcherMetrics) ObserveQueues(queue ServerQueueMap) {
	if d.meter == nil {
		return
	}

	_, err := d.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			for clientID, queueSize := range queue.SizePerClient() {
				obs.ObserveInt64(
					d.requestQueue,
					int64(queueSize),
					metric.WithAttributes(attribute.String("client_id", clientID)),
				)
			}
			return nil
		},
		d.requestQueue)
	if err != nil {
		log.Errorf("failed to register callback for dispatcher queue size: %v", err)
	}
}
