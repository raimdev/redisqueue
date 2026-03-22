# Project Summary

`redisqueue` is a lightweight Go library that implements a queue on top of Redis Streams.
It provides a `Producer` API for enqueueing messages and a `Consumer` API for concurrent
processing with at-least-once delivery. The consumer uses Redis consumer groups,
visibility timeouts, and reclaim logic to reassign messages when workers fail.

## Key Concepts

- Streams: each queue is a Redis stream.
- Consumer groups: consumers coordinate via Redis groups and acknowledgments.
- Visibility timeout: idle pending messages are reclaimed after a timeout.
- Concurrency: multiple goroutines process messages with a shared in-memory buffer.

## Main Components

- `producer.go`: enqueue messages with `XADD` and configurable stream trimming.
- `consumer.go`: read/claim messages with `XREADGROUP` and `XPENDING`/`XCLAIM`.
- `redis.go`: Redis client setup and server version checks.
- `message.go`: message model shared by producer and consumer.

## Delivery Semantics

- At-least-once delivery. Messages are acknowledged after successful processing.
- Failed or timed-out messages are reclaimed and reprocessed.

## Recommended Usage

- Use custom `ProducerOptions` and `ConsumerOptions` for production defaults.
- Keep `BlockingTimeout` non-zero to allow graceful shutdown.
- Set `StreamMaxLength` high enough to avoid dropping unprocessed messages.
