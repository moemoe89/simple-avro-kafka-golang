//
//  Practicing Avro Kafka
//
//  Copyright © 2016. All rights reserved.
//

package main

import (
	conf "github.com/moemoe89/simple-avro-kafka-golang/consumer/config"

	"bytes"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/linkedin/goavro"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	KAFKA_TOPIC = "simple-avro-kafka-golang"
)

func main() {

	recordSchemaJSON := `
	{
	  "type": "record",
	  "name": "users",
	  "doc:": "A basic schema for storing data of users",
	  "namespace": "com.avro.kafka.golang",
	  "fields": [
		{
		  "doc": "Name of role",
		  "type": "string",
		  "name": "role"
		},
		{
		  "doc": "The content of the user's data",
		  "type": "string",
		  "name": "data"
		},
		{
		  "doc": "Unix epoch time in milliseconds",
		  "type": "long",
		  "name": "timestamp"
		}
	  ]
	}
	`

	consumer, err := conf.InitKafkaConsumer()
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(KAFKA_TOPIC, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	log.Print("Connected to kafka broker")

	for m := range partitionConsumer.Messages() {

		codec, err := goavro.NewCodec(recordSchemaJSON)
		if err != nil {
			panic(err)
		}

		encoded := []byte(m.Value)
		bb := bytes.NewBuffer(encoded)
		decoded, err := codec.Decode(bb)

		record := decoded.(*goavro.Record)
		log.Print("Record Name:", record.Name)
		log.Print("Record Fields:")
		for i, field := range record.Fields {
			log.Print("field", i, field.Name, ":", field.Datum)

			if field.Name == "com.avro.kafka.golang.data" {

				data, _ := field.Datum.(string)
				bytesText := []byte(data)

				var user User
				json.Unmarshal(bytesText, &user)

				log.Print("raw data : ", data)
				log.Print("id : ", user.ID)
				log.Print("name : ", user.Name)

			}

		}

	}

}
