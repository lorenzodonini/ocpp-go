package ws

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	chargePointsConnectedMetric = "websocket_charge_points_connected"
	messageRateMetric           = "websocket_message_rate"
	pingPongDurationMetric      = "websocket_ping_pong_duration"
	attributeChargePointId      = "charge_point_id"
	attributeDirection          = "direction"
)

const (
	directionInbound  = "inbound"
	directionOutbound = "outbound"
)

type metrics struct {
	connectedChargePoints int64
	mu                    sync.Mutex

	chargePointsConnectedMetric metric.Int64ObservableGauge
	pingPongDurationMetric      metric.Float64Histogram
	messageRate                 metric.Int64Histogram
}

func newMetrics(meterProvider metric.MeterProvider) (*metrics, error) {
	meter := meterProvider.Meter("websocket")
	m := &metrics{}

	chargePointsConnected, err := meter.Int64ObservableGauge(
		chargePointsConnectedMetric,
		metric.WithDescription("Number of currently connected charge points"),
		metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
			io.Observe(m.connectedChargePoints)
			return nil
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", chargePointsConnectedMetric))
	}

	messageRate, err := meter.Int64Histogram(
		messageRateMetric,
		metric.WithDescription("Message rate"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", messageRateMetric))
	}

	m.pingPongDurationMetric, err = meter.Float64Histogram(
		pingPongDurationMetric,
		metric.WithDescription("Duration of ping-pong messages"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", pingPongDurationMetric))
	}

	m.chargePointsConnectedMetric = chargePointsConnected
	m.messageRate = messageRate

	return m, nil
}

func (m *metrics) IncrementChargePoints() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.AddInt64(&m.connectedChargePoints, 1)
}

func (m *metrics) DecrementChargePoints() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Only positive values are allowed
	if m.connectedChargePoints == 0 {
		return
	}

	atomic.AddInt64(&m.connectedChargePoints, -1)
}

func (m *metrics) RecordMessageRate(ctx context.Context, chargePointId string, direction string) {
	attributes := metric.WithAttributes(
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeDirection, direction),
	)
	m.messageRate.Record(ctx, 1, attributes)
}

func (m *metrics) RecordPingPongDuration(ctx context.Context, duration time.Duration, chargePointId string) {
	attributes := metric.WithAttributes(
		attribute.String(attributeChargePointId, chargePointId),
	)
	m.pingPongDurationMetric.Record(ctx, duration.Seconds(), attributes)
}
