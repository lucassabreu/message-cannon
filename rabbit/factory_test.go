package rabbit

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"go.uber.org/zap"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_factory(t *testing.T) {
	c := getConfig(t, "valid_queue_and_exchange_config.yml")
	factory, err := NewFactory(c, zap.NewNop())
	failIfErr(t, err, "Failed to create the factory")
	assert.Len(t, factory.conns, 2)
	t.Run("When call CreateConsumers and we got all the consumers from config", func(t *testing.T) {
		var consumers []*consumer
		consumers, err = factory.CreateConsumers()
		failIfErr(t, err, "Failed to create all the consumers")
		assert.Len(t, consumers, 2)
		for _, consumer := range consumers {
			assert.True(t, consumer.Alive(), "The consumer ", consumer.Name(), "is not alive")
		}
	})
	t.Run("When call CreateConsumer and we got a specific consumer", func(t *testing.T) {
		var consumer *consumer
		consumer, err = factory.CreateConsumer("test1")
		failIfErr(t, err, "Failed to create all the consumers")
		assert.NotNil(t, consumer)
		assert.True(t, consumer.Alive(), "The consumer ", consumer.Name(), "is not alive")
	})
	ch, err := factory.conns["default"].Channel()
	failIfErr(t, err, "Error opening a channel")
	for _, cfg := range factory.config.Consumers {
		_, err := ch.QueueDelete(cfg.Queue.Name, false, false, false)
		failIfErr(t, err)
	}
	for name := range factory.config.Exchanges {
		err := ch.ExchangeDelete(name, false, false)
		failIfErr(t, err)
	}
}

func Test_factory_should_return_error(t *testing.T) {
	c := getConfig(t, "valid_queue_and_exchange_config.yml")
	t.Run("On factory, when we pass a invalid dsn we get an error", func(t *testing.T) {
		conn := c.Connections["test1"]
		conn.DSN = "amqp://guest:guest@localhost:5672/foo"
		c.Connections["test1"] = conn
		_, err := NewFactory(c, zap.NewNop())
		assert.EqualError(t, err, "error opening the connection \"test1\": Exception (403) Reason: \"no access to this vhost\"")
	})
}

func getConfig(t *testing.T, configFile string) Config {
	c := Config{}
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("test")
	yaml, err := ioutil.ReadFile(filepath.Join("testdata", configFile))
	failIfErr(t, err, "Failed to read the config file: ")
	err = viper.ReadConfig(bytes.NewBuffer(yaml))
	failIfErr(t, err)
	viper.AutomaticEnv()

	err = viper.UnmarshalKey("rabbitmq", &c)
	failIfErr(t, err, "Failed to marshal the config struct: ")
	return c
}

func failIfErr(t *testing.T, err error, msg ...interface{}) {
	if err != nil {
		t.Fatal(err, msg)
	}
}