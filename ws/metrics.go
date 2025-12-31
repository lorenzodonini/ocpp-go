package ws

import (
	"context"
	"fmt"
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

type serverMetrics struct {
	connectedChargePoints atomic.Int64

	chargePointsConnectedMetric metric.Int64ObservableGauge
	pingPongDurationMetric      metric.Float64Histogram
	messageRate                 metric.Int64Histogram
}

func newServerMetrics(meterProvider metric.MeterProvider) (*serverMetrics, error) {
	meter := meterProvider.Meter("websocket")
	m := &serverMetrics{}

	chargePointsConnected, err := meter.Int64ObservableGauge(
		chargePointsConnectedMetric,
		metric.WithDescription("Number of currently connected charge points"),
		metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
			io.Observe(m.connectedChargePoints.Load())
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

func (m *serverMetrics) IncrementChargePoints() {
	m.connectedChargePoints.Add(1)
}

func (m *serverMetrics) DecrementChargePoints() {
	// Only positive values are allowed
	if m.connectedChargePoints.Load() == 0 {
		return
	}

	m.connectedChargePoints.Add(1)
}

func (m *serverMetrics) RecordMessageRate(ctx context.Context, chargePointId string, direction string) {
	attributes := metric.WithAttributes(
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeDirection, direction),
	)
	m.messageRate.Record(ctx, 1, attributes)
}

func (m *serverMetrics) RecordPingPongDuration(ctx context.Context, duration time.Duration, chargePointId string) {
	attributes := metric.WithAttributes(
		attribute.String(attributeChargePointId, chargePointId),
	)
	m.pingPongDurationMetric.Record(ctx, duration.Seconds(), attributes)
}
