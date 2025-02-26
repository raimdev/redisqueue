package redisqueue

// Message constitutes a message that will be enqueued and dequeued from Redis.
// When enqueuing, it's recommended to leave ID empty and let Redis generate it,
// unless you know what you're doing.
type Message struct {
	ID              string
	Stream          string
	StreamMaxLength int64
	StreamMinID     string
	TrimLimit       int64
	Values          map[string]interface{}
}
