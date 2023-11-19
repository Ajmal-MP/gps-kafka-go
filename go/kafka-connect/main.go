package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

const (
	broker = "139.5.189.204:9092" // Replace with your Kafka broker address
	topic  = "my_topic"           // Replace with your Kafka topic
)

func produceMessages() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		log.Fatal("Error creating Kafka producer:", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal("Error closing Kafka producer:", err)
		}
	}()

	for i := 0; i < 10; i++ {
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("Message %d", i)),
		}

		_, _, err := producer.SendMessage(message)
		if err != nil {
			log.Println("Failed to produce message:", err)
		} else {
			fmt.Println("Produced message:", message.Value)
		}

		time.Sleep(time.Second)
	}
}

func consumeMessages() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatal("Error creating Kafka consumer:", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal("Error closing Kafka consumer:", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error creating partition consumer:", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Println("Consumed message:", string(msg.Value))
			case err := <-partitionConsumer.Errors():
				fmt.Println("Error:", err)
			case <-signals:
				return
			}
		}
	}()

	wg.Wait()
}

func main() {
	go produceMessages()
	go consumeMessages()

	// Let the program run for some time to produce and consume messages
	time.Sleep(20 * time.Second)
}
