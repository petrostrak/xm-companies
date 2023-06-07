package producer

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/petrostrak/xm-companies/internal/core/domain"
)

func GetNewProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "socket.gethostname()",
		"acks":              "all",
	})

	return producer, err
}

func Produce(key []byte, value []byte, method string, topic string, producer *kafka.Producer) {
	deliveryChan := make(chan kafka.Event, 1)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
		Headers: []kafka.Header{
			{
				Key:   "Method",
				Value: []byte(method),
			},
		},
	}, deliveryChan)

	if err != nil {
		panic(err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("delivery failed %v \n", m.TopicPartition.Error)
	} else {
		fmt.Printf("message delivered topic: %s | key: %s\n", topic, string(key))
	}

	close(deliveryChan)
}

func ProduceCompany(company *domain.Company, method string) error {
	prod, err := GetNewProducer()
	if err != nil {
		return err
	}

	companyTopic := "producer.company"

	data, err := json.Marshal(&company)
	if err != nil {
		return err
	}
	Produce([]byte(company.ID.String()), data, method, companyTopic, prod)

	return nil
}
