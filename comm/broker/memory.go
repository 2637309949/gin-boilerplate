package broker

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type memoryBroker struct {
	opts Options

	addr string
	sync.RWMutex
	connected   bool
	Subscribers map[string][]*memorySubscriber
}

type memorySubscriber struct {
	id      string
	topic   string
	exit    chan bool
	handler Handler
	opts    SubscribeOptions
}

func (m *memoryBroker) Options() Options {
	return m.opts
}

func (m *memoryBroker) Address() string {
	return m.addr
}

func (m *memoryBroker) Connect() error {
	m.Lock()
	defer m.Unlock()

	if m.connected {
		return nil
	}

	// use 127.0.0.1 to avoid scan of all network interfaces
	addr, err := Extract("127.0.0.1")
	if err != nil {
		return err
	}
	i := rand.Intn(20000)
	// set addr with port
	addr = HostPort(addr, 10000+i)

	m.addr = addr
	m.connected = true

	return nil
}

func (m *memoryBroker) Disconnect() error {
	m.Lock()
	defer m.Unlock()

	if !m.connected {
		return nil
	}

	m.connected = false

	return nil
}

func (m *memoryBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *memoryBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return errors.New("not connected")
	}

	subs, ok := m.Subscribers[topic]
	m.RUnlock()
	if !ok {
		return nil
	}

	for _, sub := range subs {
		if err := sub.handler(msg); err != nil {
			if eh := sub.opts.ErrorHandler; eh != nil {
				eh(msg, err)
			}
			continue
		}
	}

	return nil
}

func (m *memoryBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, errors.New("not connected")
	}
	m.RUnlock()

	var options SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	sub := &memorySubscriber{
		exit:    make(chan bool, 1),
		id:      uuid.New().String(),
		topic:   topic,
		handler: handler,
		opts:    options,
	}

	m.Lock()
	m.Subscribers[topic] = append(m.Subscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		var newSubscribers []*memorySubscriber
		for _, sb := range m.Subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.Subscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil
}

func (m *memoryBroker) String() string {
	return "memory"
}

func (m *memorySubscriber) Options() SubscribeOptions {
	return m.opts
}

func (m *memorySubscriber) Topic() string {
	return m.topic
}

func (m *memorySubscriber) Unsubscribe() error {
	m.exit <- true
	return nil
}

func NewMemoryBroker(opts ...Option) Broker {
	options := Options{
		Context: context.Background(),
	}

	rand.Seed(time.Now().UnixNano())
	for _, o := range opts {
		o(&options)
	}

	return &memoryBroker{
		opts:        options,
		Subscribers: make(map[string][]*memorySubscriber),
	}
}
