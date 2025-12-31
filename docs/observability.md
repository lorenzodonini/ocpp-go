# Observability in ocpp-go

## Metrics

The library currently supports only websocket and ocpp-j server metrics, which are exported via OpenTelemetry.
To enable metrics, you need to set the metrics exporter on the server:

```go
// sets up OTLP metrics exporter
func setupMetrics(address string) error {
grpcOpts := []grpc.DialOption{
grpc.WithTransportCredentials(insecure.NewCredentials()),
}

client, err := grpc.NewClient(address, grpcOpts...)

if err != nil {
return errors.Wrap(err, "failed to create gRPC connection to collector")
}

ctx := context.Background()

exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(client))
if err != nil {
return errors.Wrap(err, "failed to create otlp metric exporter")
}

resource, err := resource.New(ctx,
resource.WithAttributes(
semconv.ServiceNameKey.String("centralSystem-demo"),
semconv.ServiceVersionKey.String("example"),
),
resource.WithFromEnv(),
resource.WithContainer(),
resource.WithOS(),
resource.WithOSType(),
resource.WithHost(),
)
if err != nil {
return errors.Wrap(err, "failed to create resource")
}

meterProvider := metricsdk.NewMeterProvider(
metricsdk.WithReader(
metricsdk.NewPeriodicReader(
exporter,
metricsdk.WithInterval(10*time.Second),
),
),
metricsdk.WithResource(resource),
)

otel.SetMeterProvider(meterProvider)
return nil
}
```

You can check out the [reference implementation](../example/1.6/cs/central_system_sim.go), and deploy it with:

```bash
make example-ocpp16-observability
```

> Note: Deploying the example requires docker and docker compose to be installed.

The deployment will start a central system with metrics enabled and a
full [observability stack](https://github.com/grafana/docker-otel-lgtm).

You can log in to Grafana at http://localhost:3000 with the credentials `admin/admin`.