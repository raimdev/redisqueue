# Roadmap

These items move the project toward a robust alternative to RabbitMQ-class systems.

1. Delivery controls: explicit Ack/Nack API, configurable retry policy with backoff and jitter.
2. Dead-letter queue (DLQ) stream with failure metadata and original payload.
3. Observability: Prometheus metrics, structured logs, and optional OpenTelemetry spans.
4. Operational guidance: persistence requirements (AOF), HA topology recommendations, and cluster/Sentinel examples.
5. Performance: use `XAUTOCLAIM` for reclaiming, and ensure fair scheduling across streams.
6. Ergonomics: add `Close()` to release Redis connections and option to disable signal handling.
7. Testing: integration tests for shutdown draining, max-delivery, DLQ routing, and ctx cancellation under load.
