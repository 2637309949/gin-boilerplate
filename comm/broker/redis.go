package broker

import (
	"context"
	"errors"
	"gin-boilerplate/comm/codec"
	"gin-boilerplate/comm/codec/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	DefaultMaxActive      = 0
	DefaultMaxIdle        = 5
	DefaultIdleTimeout    = 2 * time.Minute
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = 5 * time.Second
)

// publication is an internal publication for the Redis
type publication struct {
	topic   string
	message *Message
	err     error
}

// Topic returns the topic this publication applies to.
func (p *publication) Topic() string {
	return p.topic
}

// Message returns the broker message of the publication.
func (p *publication) Message() *Message {
	return p.message
}

// Ack sends an acknowledgement to the  However this is not supported
// is Redis and therefore this is a no-op.
func (p *publication) Ack() error {
	return nil
}

func (p *publication) Error() error {
	return p.err
}

// subscriber proxies and handles Redis messages as broker publications.
type subscriber struct {
	codec  codec.Marshaler
	conn   *redis.PubSub
	topic  string
	handle Handler
	opts   SubscribeOptions
}

// recv loops to receive new messages from Redis and handle them
// as publications.
func (s *subscriber) recv() {
	// Close the connection once the subscriber stops receiving.
	defer s.conn.Close()

	for {
		m, err := s.conn.Receive(s.opts.Context)
		if err != nil {
			return
		}
		switch x := m.(type) {
		case redis.Message:
			var m Message

			// Handle error? Only a log would be necessary since this type
			// of issue cannot be fixed.
			if err := s.codec.Unmarshal([]byte(x.Payload), &m); err != nil {
				break
			}

			p := publication{
				topic:   x.Channel,
				message: &m,
			}

			// Handle error? Retry?
			if p.err = s.handle(&m); p.err != nil {
				break
			}

			// Added for posterity, however Ack is a no-op.
			if s.opts.AutoAck {
				if err := p.Ack(); err != nil {
					break
				}
			}

		case redis.Subscription:
			if x.Count == 0 {
				return
			}

		case error:
			return
		}
	}
}

// Options returns the subscriber options.
func (s *subscriber) Options() SubscribeOptions {
	return s.opts
}

// Topic returns the topic of the subscriber.
func (s *subscriber) Topic() string {
	return s.topic
}

// Unsubscribe unsubscribes the subscriber and frees the connection.
func (s *subscriber) Unsubscribe() error {
	return s.conn.Unsubscribe(s.opts.Context)
}

// broker implementation for Redis.
type redisBroker struct {
	addr string
	pool *redis.Client
	opts Options
}

// String returns the name of the broker implementation.
func (b *redisBroker) String() string {
	return "redis"
}

// Options returns the options defined for the
func (b *redisBroker) Options() Options {
	return b.opts
}

// Address returns the address the broker will use to create new connections.
// This will be set only after Connect is called.
func (b *redisBroker) Address() string {
	return b.addr
}

// Init sets or overrides broker options.
func (b *redisBroker) Init(opts ...Option) error {
	if b.pool != nil {
		return errors.New("redis: cannot init while connected")
	}

	for _, o := range opts {
		o(&b.opts)
	}

	return nil
}

// Connect establishes a connection to Redis which provides the
// pub/sub implementation.
func (b *redisBroker) Connect() error {
	if b.pool != nil {
		return nil
	}

	var addr string
	if len(b.opts.Addrs) == 0 || b.opts.Addrs[0] == "" {
		addr = "redis://127.0.0.1:6379"
	} else {
		addr = b.opts.Addrs[0]
		if !strings.HasPrefix("redis://", addr) {
			addr = "redis://" + addr
		}
	}
	b.pool = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return nil
}

// Disconnect closes the connection pool.
func (b *redisBroker) Disconnect() error {
	err := b.pool.Close()
	b.pool = nil
	b.addr = ""
	return err
}

// Publish publishes a message.
func (b *redisBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	v, err := b.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}

	err = b.pool.Publish(b.opts.Context, topic, v).Err()
	return err
}

// Subscribe returns a subscriber for the topic and handler.
func (b *redisBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	var options SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	s := subscriber{
		codec:  b.opts.Codec,
		conn:   b.pool.Subscribe(b.opts.Context, topic),
		topic:  topic,
		handle: handler,
		opts:   options,
	}

	// Run the receiver routine.
	go s.recv()

	return &s, nil
}

// NewBroker returns a new broker implemented using the Redis pub/sub
// protocol. The connection address may be a fully qualified IANA address such
// as: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
func NewBroker(opts ...Option) Broker {

	// Initialize with empty broker options.
	options := Options{
		Codec:   json.Marshaler{},
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &redisBroker{
		opts: options,
	}
}
