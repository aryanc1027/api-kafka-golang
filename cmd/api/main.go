package main

import (
	"log"

	"github.com/aryanc1027/api-kafka-golang/internal/api"
	"github.com/aryanc1027/api-kafka-golang/internal/config"
	"github.com/aryanc1027/api-kafka-golang/internal/kafka"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}


	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()


	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()


	router := api.SetupRouter(cfg, producer, consumer)


	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}