package main

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	"github.com/Shopify/sarama"
	"log"
	"strings"
)

type User struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
}

const (
	PRODUCER_URL string = "localhost:9092"
	KAFKA_TOPIC string = "simple-avro-kafka-golang"
)

func main(){

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

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Consumer.Return.Errors = true

	producerUrl := strings.Split(PRODUCER_URL, ",")

	consumer, err := sarama.NewConsumer(producerUrl, config)
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

				log.Print("raw data : ",data)
				log.Print("id : ",user.ID)
				log.Print("name : ",user.Name)

			}

		}

	}

}