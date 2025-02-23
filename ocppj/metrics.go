package ocppj

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	requestsInboundMetric        = "ocpp_requests_inbound"
	requestsInboundFailedMetric  = "ocpp_requests_inbound_failed"
	requestsOutboundMetric       = "ocpp_requests_outbound"
	requestsOutboundFailedMetric = "ocpp_requests_outbound_failed"
)

const (
	attributeChargePointId = "charge_point_id"
	attributeOcppVersion   = "ocpp_version"
	attributeFeature       = "feature"
	attributeError         = "error"
)

type OcppMetricsError string

const (
	ChargePointError = OcppMetricsError("charge_point_error")
	PayloadError     = OcppMetricsError("payload_error")
	ValidationError  = OcppMetricsError("validation_error")
)

type Opt func(metricContext *MetricContext)

func WithAdditionalAttributes(attrs []attribute.KeyValue) Opt {
	return func(metricContext *MetricContext) {
		metricContext.additionalAttributes = attrs
	}
}

type MetricContext struct {
	ChargePointId string
	Message       string
	// Error if the request failed
	Error                *OcppMetricsError
	additionalAttributes []attribute.KeyValue
}

func (ctx *MetricContext) toAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, ctx.ChargePointId),
		attribute.String(attributeFeature, ctx.Message),
	}

	if ctx.Error != nil {
		attrs = append(attrs, attribute.String(attributeError, string(*ctx.Error)))
	}

	if ctx.additionalAttributes != nil {
		attrs = append(attrs, ctx.additionalAttributes...)
	}

	return attrs
}

type ocppMetrics struct {
	requestsIn        metric.Int64Histogram
	requestsOut       metric.Int64Histogram
	requestsInFailed  metric.Int64Histogram
	requestsOutFailed metric.Int64Histogram
}

// newOcppMetrics Creates a new metrics instance
func newOcppMetrics(meterProvider metric.MeterProvider, ocppVersion string) (*ocppMetrics, error) {
	if meterProvider == nil {
		return nil, errors.New("meterProvider is required")
	}

	meter := meterProvider.Meter(
		"ocpp",
		metric.WithInstrumentationAttributes(attribute.String(attributeOcppVersion, ocppVersion)),
	)

	requestsIn, err := meter.Int64Histogram(
		requestsInboundMetric,
		metric.WithDescription("Number of inbound requests"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsInboundMetric))
	}

	requestsInFailed, err := meter.Int64Histogram(
		requestsInboundFailedMetric,
		metric.WithDescription("Number of inbound requests failed"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsInboundFailedMetric))
	}

	requestsOut, err := meter.Int64Histogram(
		requestsOutboundMetric,
		metric.WithDescription("Number of outbound requests"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsOutboundMetric))
	}

	requestsOutFailed, err := meter.Int64Histogram(
		requestsOutboundFailedMetric,
		metric.WithDescription("Number of outbound requests failed"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsOutboundFailedMetric))
	}

	metrics := &ocppMetrics{
		requestsIn:        requestsIn,
		requestsInFailed:  requestsInFailed,
		requestsOut:       requestsOut,
		requestsOutFailed: requestsOutFailed,
	}
	return metrics, nil
}

func (m *ocppMetrics) IncrementInboundRequests(ctx context.Context, chargePointId, requestName string, error *OcppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeFeature, requestName),
	}

	if error != nil {
		attrs = append(attrs, attribute.String(attributeError, string(*error)))
	}

	metricAttrs := metric.WithAttributes(attrs...)
	if error != nil {
		m.requestsInFailed.Record(ctx, 1, metricAttrs)
	}

	m.requestsIn.Record(ctx, 1, metricAttrs)
}

func (m *ocppMetrics) IncrementInboundFailedRequests(ctx context.Context, chargePointId, requestName string, error OcppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeFeature, requestName),
		attribute.String(attributeError, string(error)),
	}
	m.requestsInFailed.Record(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *ocppMetrics) IncrementOutboundRequests(ctx context.Context, chargePointId, requestName string, error *OcppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeFeature, requestName),
	}

	if error != nil {
		attrs = append(attrs, attribute.String(attributeError, string(*error)))
	}

	metricAttrs := metric.WithAttributes(attrs...)
	if error != nil {
		m.requestsOutFailed.Record(ctx, 1, metricAttrs)
	}

	m.requestsOut.Record(ctx, 1, metricAttrs)
}

func (m *ocppMetrics) IncrementOutboundFailedRequests(ctx context.Context, chargePointId, requestName string, error OcppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
		attribute.String(attributeFeature, requestName),
		attribute.String(attributeError, string(error)),
	}
	m.requestsOutFailed.Record(ctx, 1, metric.WithAttributes(attrs...))
}
