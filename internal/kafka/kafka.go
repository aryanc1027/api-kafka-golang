package kafka 

import (
    "github.com/IBM/sarama"
    "github.com/aryanc1027/api-kafka-golang/internal/config"
)


type Producer interface {
    SendMessage(topic string, message []byte) error
    Close() error
}


type Consumer interface {
    ConsumeMessages(topic string, messages chan<- []byte, errors chan<- error)
    Close() error
}

type SaramaProducer struct {
    producer sarama.SyncProducer
}

// SaramaConsumer implements the Consumer interface
type SaramaConsumer struct {
    consumer sarama.Consumer
}

func NewProducer(cfg *config.Config) (Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5

    producer, err := sarama.NewSyncProducer(cfg.KafkaHosts, config)
    if err != nil {
        return nil, err
    }
    return &SaramaProducer{producer: producer}, nil
}

func (p *SaramaProducer) SendMessage(topic string, message []byte) error {
    _, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.ByteEncoder(message),
    })
    return err
}

func (p *SaramaProducer) Close() error {
    return p.producer.Close()
}

func NewConsumer(cfg *config.Config) (Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true

    consumer, err := sarama.NewConsumer(cfg.KafkaHosts, config)
    if err != nil {
        return nil, err
    }

    return &SaramaConsumer{consumer: consumer}, nil
}

func (c *SaramaConsumer) ConsumeMessages(topic string, messages chan<- []byte, errors chan<- error) {
    partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        errors <- err
        return
    }
    defer partitionConsumer.Close()

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            messages <- msg.Value
        case err := <-partitionConsumer.Errors():
            errors <- err
        }
    }
}

func (c *SaramaConsumer) Close() error {
    return c.consumer.Close()
}