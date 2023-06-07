package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer[T comparable] struct {
	reader *kafka.Reader
	dialer *kafka.Dialer
	topic  string
}

func (c *Consumer[T]) CreateConnection() {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     c.topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Millisecond * 10,
		Dialer:    c.dialer,
	})

	if err := c.reader.SetOffset(0); err != nil {
		log.Println(err)
	}
}

func (c *Consumer[T]) Read(model T, callback func(T, error)) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*80)
		defer cancel()
		message, err := c.reader.ReadMessage(ctx)

		if err != nil {
			callback(model, err)
			return
		}

		err = json.Unmarshal(message.Value, &model)

		if err != nil {
			callback(model, err)
			continue
		}

		callback(model, nil)
	}
}

func ConsumeCompany() {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	companyConsumer := Consumer[Company]{
		dialer: dialer,
		topic:  "producer-image-table",
	}

	companyConsumer.CreateConnection()

	companyConsumer.Read(Company{}, func(company Company, err error) {
		collectedCompany := CollectedCompany{
			ID:                company.ID,
			Name:              company.Name,
			Description:       company.Name,
			NumberOfEmployees: company.NumberOfEmployees,
			Registered:        company.Registered,
			Type:              company.Type,
			CreatedAt:         time.Now(),
		}

		fmt.Println(collectedCompany)
	})

	if err := companyConsumer.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
