//
//  Practicing Kafka
//
//  Copyright © 2016. All rights reserved.
//

package config

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
)

// InitKafkaConsumer will create a variable that represent the sarama.Consumer
func InitKafkaConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	producerUrl := strings.Split(Configuration.Kafka.Addr, ",")

	consumer, err := sarama.NewConsumer(producerUrl, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping connection to kafka: %s", err.Error())
	}

	return consumer, nil
}
