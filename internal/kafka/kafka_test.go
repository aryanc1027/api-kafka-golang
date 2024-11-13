package kafka

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/aryanc1027/api-kafka-golang/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	cfg := &config.Config{
		KafkaHosts: []string{"localhost:9092"},
	}

	producer, err := NewProducer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	err = producer.Close()
	assert.NoError(t, err)
}

func TestProducer_SendMessage(t *testing.T) {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true


	mockProducer := mocks.NewSyncProducer(t, cfg)
	mockProducer.ExpectSendMessageAndSucceed()


	producer := &SaramaProducer{
		producer: mockProducer,
	}


	err := producer.SendMessage("test-topic", []byte("test-message"))
	assert.NoError(t, err)


	err = mockProducer.Close()
	assert.NoError(t, err)
}

func TestNewConsumer(t *testing.T) {
	cfg := &config.Config{
		KafkaHosts: []string{"localhost:9092"},
	}

	consumer, err := NewConsumer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	err = consumer.Close()
	assert.NoError(t, err)
}

func TestConsumer_ConsumeMessages(t *testing.T) {

	cfg := sarama.NewConfig()


	mockConsumer := mocks.NewConsumer(t, cfg)


	consumer := &SaramaConsumer{
		consumer: mockConsumer,
	}


	mockPartitionConsumer := mockConsumer.ExpectConsumePartition("test-topic", 0, sarama.OffsetNewest)
	mockPartitionConsumer.ExpectMessagesDrainedOnClose()

	message := &sarama.ConsumerMessage{
		Value: []byte("test-message"),
	}
	mockPartitionConsumer.YieldMessage(message)


	messages := make(chan []byte)
	errors := make(chan error)


	go consumer.ConsumeMessages("test-topic", messages, errors)


	select {
	case msg := <-messages:
		assert.Equal(t, []byte("test-message"), msg)
	case err := <-errors:
		t.Fatalf("Unexpected error: %v", err)
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for message")
	}


	mockPartitionConsumer.Close()


	err := mockConsumer.Close()
	assert.NoError(t, err)
}