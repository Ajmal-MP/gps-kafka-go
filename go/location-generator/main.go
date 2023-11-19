package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	VehicleID int     `json:"vehicle_id"`
}

func main() {
	currentLocation := Location{
		Latitude:  37.7749,
		Longitude: -122.4194,
		VehicleID: 1,
	}

	// Encode the Location struct to JSON
	locationJSON, err := json.Marshal(currentLocation)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Set up the Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokerAddress := "139.5.189.204:9092" // Replace with your Kafka broker address
	topic := fmt.Sprintf("vehicleID_%d", currentLocation.VehicleID)

	// Create a Kafka producer
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, config)
	if err != nil {
		fmt.Println("Error creating Kafka producer:", err)
		return
	}
	defer producer.Close()

	// Produce a message to the Kafka topic
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(locationJSON),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Error writing to Kafka:", err)
		return
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
