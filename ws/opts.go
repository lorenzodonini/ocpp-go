package ws

import "go.opentelemetry.io/otel/metric"

type ServerOpt func(server *ServerOpts)

type ServerOpts struct {
	meterProvider metric.MeterProvider
}

func WithMeterProvider(meterProvider metric.MeterProvider) ServerOpt {
	return func(opts *ServerOpts) {
		opts.meterProvider = meterProvider
	}
}
