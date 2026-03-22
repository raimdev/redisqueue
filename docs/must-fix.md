# Must-fix

These items are required to improve reliability and correctness.

- [x] Fix shutdown semantics to avoid dropping buffered messages.
- [x] Avoid deadlocks by making error reporting non-blocking (buffered errors channel or handler).
- [x] Enforce backpressure when the in-memory queue is full (do not issue reads/claims with Count 0).
- [x] Honor context cancellation by initiating shutdown when ctx is done.
- [x] Validate options: require `Concurrency > 0` and `BufferSize > 0`, or explicitly handle zero values.
