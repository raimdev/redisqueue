package redisqueue

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

// ProducerOptions provide options to configure the Producer.
type ProducerOptions struct {
	// StreamMaxLength sets the MAXLEN option when calling XADD. This creates a
	// capped stream to prevent the stream from taking up memory indefinitely.
	// It's important to note though that this isn't the maximum number of
	// _completed_ messages, but the maximum number of _total_ messages. This
	// means that if all consumers are down, but producers are still enqueuing,
	// and the maximum is reached, unprocessed message will start to be dropped.
	// So ideally, you'll set this number to be as high as you can makee it.
	// More info here: https://redis.io/commands/xadd#capped-streams.
	StreamMaxLength int64

	// StreamMinID sets the minimum ID that will be used when calling XADD. This
	// is useful when you want to ensure that the stream is trimmed to a certain
	// point. More info here: https://redis.io/commands/xadd#capped-streams.
	StreamMinID string

	// UseApproximate determines whether to use the ~ with the MAXLEN and MINID
	// option. This allows the stream trimming to done in a more efficient
	// manner. More info here: https://redis.io/commands/xadd#capped-streams.
	UseApproximate bool

	// TrimLimit sets LIMIT for XADD
	TrimLimit int64

	// RedisClient supersedes the RedisOptions field, and allows you to inject
	// an already-made Redis Client for use in the consumer. This may be either
	// the standard client or a cluster client.
	RedisClient redis.UniversalClient
	// RedisOptions allows you to configure the underlying Redis connection.
	// More info here:
	// https://pkg.go.dev/github.com/redis/go-redis/v9#Options.
	//
	// This field is used if RedisClient field is nil.
	RedisOptions *RedisOptions
}

// Producer adds a convenient wrapper around enqueuing messages that will be
// processed later by a Consumer.
type Producer struct {
	options *ProducerOptions
	redis   redis.UniversalClient
}

var defaultProducerOptions = &ProducerOptions{
	StreamMaxLength: 1000,
	UseApproximate:  true,
}

// NewProducer uses a default set of options to create a Producer. It sets
// StreamMaxLength to 1000 and UseApproximate to true. In most production
// environments, you'll want to use NewProducerWithOptions.
func NewProducer() (*Producer, error) {
	return NewProducerWithOptions(context.Background(), defaultProducerOptions)
}

// NewProducerWithOptions creates a Producer using custom ProducerOptions.
func NewProducerWithOptions(ctx context.Context, options *ProducerOptions) (*Producer, error) {
	var r redis.UniversalClient

	if options.RedisClient != nil {
		r = options.RedisClient
	} else {
		r = newRedisClient(options.RedisOptions)
	}

	if err := redisPreflightChecks(ctx, r); err != nil {
		return nil, err
	}

	return &Producer{
		options: options,
		redis:   r,
	}, nil
}

// Enqueue takes in a pointer to Message and enqueues it into the stream set at
// msg.Stream. While you can set msg.ID, unless you know what you're doing, you
// should let Redis auto-generate the ID. If an ID is auto-generated, it will be
// set on msg.ID for your reference. msg.Values is also required.
func (p *Producer) Enqueue(ctx context.Context, msg *Message) error {
	maxLen := cmp.Or(msg.StreamMaxLength, p.options.StreamMaxLength)
	minID := cmp.Or(msg.StreamMinID, p.options.StreamMinID)
	if maxLen > 0 {
		minID = ""
	}

	args := &redis.XAddArgs{
		ID:     msg.ID,
		Stream: msg.Stream,
		Values: msg.Values,
		MaxLen: maxLen,
		MinID:  minID,
		Approx: p.options.UseApproximate,
		Limit:  cmp.Or(msg.TrimLimit, p.options.TrimLimit),
	}
	id, err := p.redis.XAdd(ctx, args).Result()
	if err != nil {
		return err
	}
	msg.ID = id
	return nil
}
