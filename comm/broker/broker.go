package broker

var (
	// DefaultBroker implementation
	DefaultBroker = NewMemoryBroker()
	_             = DefaultBroker.Connect()
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
type Handler func(*Message) error

type ErrorHandler func(*Message, error)

type Message struct {
	Header map[string]string
	Body   []byte
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

// Publish a message to a topic
func Publish(topic string, m *Message) error {
	return DefaultBroker.Publish(topic, m)
}

// Subscribe to a topic
func Subscribe(topic string, h Handler) (Subscriber, error) {
	return DefaultBroker.Subscribe(topic, h)
}
