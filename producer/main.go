package main

import (
	"bytes"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/linkedin/goavro"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	Role string `json:"role"`
	Data struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

var request Request

const (
	PRODUCER_URL string = "localhost:9092"
	KAFKA_TOPIC  string = "simple-avro-kafka-golang"
)

func message(c *gin.Context) {

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

	c.Bind(&request)
	requestMarshal, err := json.Marshal(request)

	if err != nil {
		panic(err)
	}

	requestString := string(requestMarshal)

	someRecord, err := goavro.NewRecord(goavro.RecordSchema(recordSchemaJSON))
	if err != nil {
		panic(err)
	}

	dataMarshal, err := json.Marshal(request.Data)
	if err != nil {
		panic(err)
	}

	someRecord.Set("role", request.Role)
	someRecord.Set("data", string(dataMarshal))
	someRecord.Set("com.avro.kafka.golang.timestamp", int64(1082196484))

	codec, err := goavro.NewCodec(recordSchemaJSON)
	if err != nil {
		panic(err)
	}

	bb := new(bytes.Buffer)
	if err = codec.Encode(bb, someRecord); err != nil {
		panic(err)
	}

	actual := bb.Bytes()

	dataString := string(actual)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{PRODUCER_URL}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	strTime := strconv.Itoa(int(time.Now().Unix()))

	msg := &sarama.ProducerMessage{
		Topic: KAFKA_TOPIC,
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.StringEncoder(dataString),
	}

	producer.Input() <- msg

	resp := gin.H{
		"status":  http.StatusOK,
		"message": "Message has been sent.",
		"data":    requestString,
	}

	c.IndentedJSON(http.StatusOK, resp)

}

func main() {

	router := gin.Default()
	router.POST("/", message)
	router.Run(":3000")

}
