//
//  Practicing Kafka
//
//  Copyright © 2016. All rights reserved.
//

package config

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// InitKafkaProducer will create a variable that represent the sarama.AsyncProducer
func InitKafkaProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{Configuration.Kafka.Addr}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping connection to kafka: %s", err.Error())
	}

	return producer, nil
}
