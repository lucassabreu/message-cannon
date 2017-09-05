package supervisor

import (
	"context"
	"io"
)

// Runner encapsulate what is done with messages
type Runner func(context.Context, io.Writer, []byte) error

// Factory create consumers
type Factory interface {
	// CreateConsumers will iterate over config and create all the consumers
	CreateConsumers() ([]Consumer, error)

	// CreateConsumer create a new consumer for a specific name using the config provided.
	CreateConsumer(name string) (Consumer, error)

	// Name return the factory name
	Name() string
}

// Consumer consume messages and pass to workers who will process the messages.
type Consumer interface {
	// TODO: Create the state, we will add some metrics here
	// State returns a copy of the executor's current operation state.
	// State() State

	// Run will get the messages and pass to the runner.
	Run()

	// Kill will try to stop the internal work. Return an error in case of failure.
	Kill() error

	// Alive returns true if the tomb is not in a dying or dead state.
	Alive() bool

	// Name return the consumer name
	Name() string

	// FactoryName is the name of the factory responsible for this consumer.
	FactoryName() string
}
