package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type Person struct {
	Document string `json:"document"`
}

const (
	searchDocumentEvented = "searchDocumentEvented"
	foundDocumentEvented  = "foundDocumentEvented"
)

var (
	kafkaURL = ""
	groupID  = ""
	Persons  []Person
	channel  = make(chan string, 1000)
)

func setupConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	mode, ok := viper.Get("LOCAL").(string)
	if ok && mode == "1" {
		kafkaURL = viper.Get("KAFKA.URL").(string)
		groupID = viper.Get("KAFKA.GROUPID").(string)
	} else {
		// get kafka writer using environment variables.
		kafkaURL = os.Getenv("kafkaURL")
		groupID = os.Getenv("groupID")
	}

	Persons = []Person{
		Person{Document: "45531745"},
		Person{Document: "18538847"},
	}
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func main() {
	setupConfig()
	kafkaWriter := getKafkaWriter(kafkaURL, searchDocumentEvented)

	defer kafkaWriter.Close()
	consumeSearchDocumentEvented()
}

func consumeSearchDocumentEvented() {
	reader := getKafkaReader(kafkaURL, searchDocumentEvented, groupID)

	defer reader.Close()

	fmt.Sprintf("start consuming - %s", searchDocumentEvented)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v = %s\n", m.Topic, string(m.Value))
		var personItem Person
		json.Unmarshal([]byte(string(m.Value)), &personItem)
		publishFoundDocumentEvented(personItem)
	}
}

func publishFoundDocumentEvented(personItem Person) {
	writer := getKafkaWriter(kafkaURL, foundDocumentEvented)
	defer writer.Close()

	fmt.Sprintf("start producing - %s", foundDocumentEvented)

	key := fmt.Sprintf("Key-%d", uuid.New())
	data := fmt.Sprintf("%s - No es Cliente", personItem.Document)

	for i := range Persons {
		if Persons[i].Document == personItem.Document {
			data = fmt.Sprintf("%s - Es Cliente", personItem.Document)
			break
		}
	}
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(data),
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("produced at topic:%s = %s\n", foundDocumentEvented, string(msg.Value))
	}
}
