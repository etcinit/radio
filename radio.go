package radio

import (
	"errors"
	"runtime"
)

// A Broadcaster should be able to take messages an broadcast them into multiple
// channels.
type Broadcaster interface {
	Broadcast(message interface{}) error
	Listen() (<-chan interface{}, uint)
	Stop(listener uint) error
	Call(listener uint, message interface{}) error
}

// A Radio is a simple implementation of a Broadcaster. It is capable of
// broadcasting messages to multiple channels and sending messages to individual
// channels.
type Radio struct {
	listeners map[uint]chan interface{}
	antenna   chan interface{}
	lastID    uint
	live      bool
}

// Noise is a simple message that can be sent through a Radio
type Noise struct{}

// NewRadio creates a new instance of a Radio.
func NewRadio() *Radio {
	instance := &Radio{
		antenna:   make(chan interface{}),
		listeners: make(map[uint]chan interface{}),
	}

	instance.powerup()

	runtime.SetFinalizer(instance, func(radio *Radio) {
		instance.antenna <- Noise{}
	})

	return instance
}

func repeat(listener chan interface{}, message interface{}) {
	listener <- message
}

func (r *Radio) powerup() {
	r.live = true

	go func() {
		defer close(r.antenna)

		for message := range r.antenna {
			if _, ok := message.(Noise); ok {
				r.live = false
				break
			}

			for _, listener := range r.listeners {
				go repeat(listener, message)
			}
		}
	}()
}

// Broadcast repeats the same message to all the Radio listeners.
func (r *Radio) Broadcast(message interface{}) error {
	if !r.live {
		return errors.New("This radio channel is offline")
	}

	r.antenna <- message

	return nil
}

// Listen creates a new read-only listener channel that will receive messages
// broadcast on this Radio. This function also returns an identifier that can be
// used to reference to this listerner on other Radio functions.
func (r *Radio) Listen() (<-chan interface{}, uint) {
	listener := make(chan interface{})

	r.lastID = r.lastID + 1
	r.listeners[r.lastID] = listener

	return listener, r.lastID
}

// Call sends a message only to the specified listener.
func (r *Radio) Call(listener uint, message interface{}) error {
	if ch, ok := r.listeners[listener]; ok {
		go repeat(ch, message)

		return nil
	}

	return errors.New("Listener does not exist")
}

// Stop removes a listener and closes their channel.
func (r *Radio) Stop(listener uint) error {
	if ch, ok := r.listeners[listener]; ok {
		delete(r.listeners, listener)

		close(ch)

		return nil
	}

	return errors.New("Listener does not exist")
}
