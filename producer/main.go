package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	searchDocumentEvented = "searchDocumentEvented"
	foundDocumentEvented  = "foundDocumentEvented"
)

var (
	kafkaURL = ""
	topic    = ""
	groupID  = ""
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

func searchDocumentEventHandler(kafkaWriter *kafka.Writer) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln(err)
		}
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
			Value: body,
		}
		err = kafkaWriter.WriteMessages(req.Context(), msg)

		if err != nil {
			wrt.Write([]byte(err.Error()))
			log.Fatalln(err)
		}
	})
}

func main() {
	setupConfig()
	kafkaURL := ""
	topic := ""
	mode, ok := viper.Get("LOCAL").(string)
	if ok && mode == "1" {
		kafkaURL = viper.Get("KAFKA.URL").(string)
		topic = viper.Get("KAFKA.TOPIC").(string)
	} else {
		// get kafka writer using environment variables.
		kafkaURL = os.Getenv("kafkaURL")
		topic = os.Getenv("topic")
	}

	kafkaWriter := getKafkaWriter(kafkaURL, topic)

	defer kafkaWriter.Close()
	go consumeSearchDocumentEvented()
	go consumeFoundDocumentEvented()

	// Add handle func for producer.
	http.HandleFunc("/api/v1/account/searchDocumentEvent", searchDocumentEventHandler(kafkaWriter))

	// Run the web server.
	fmt.Println("start producer-api ... !!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func consumeSearchDocumentEvented() {
	kafkaURL := ""
	topic := ""
	groupID := ""
	mode, ok := viper.Get("LOCAL").(string)
	if ok && mode == "1" {
		kafkaURL = viper.Get("KAFKA.URL").(string)
		topic = viper.Get("KAFKA.TOPIC").(string)
		groupID = viper.Get("KAFKA.GROUPID").(string)
	} else {
		// get kafka writer using environment variables.
		kafkaURL = os.Getenv("kafkaURL")
		topic = os.Getenv("topic")
		groupID = os.Getenv("groupID")
	}

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	fmt.Sprintf("start consuming - %s", searchDocumentEvented)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v = %s\n", m.Topic, string(m.Value))
		publishFoundDocumentEvented()
	}
}

func publishFoundDocumentEvented() {
	kafkaURL := ""
	topic := foundDocumentEvented
	mode, ok := viper.Get("LOCAL").(string)
	if ok && mode == "1" {
		kafkaURL = viper.Get("KAFKA.URL").(string)
	} else {
		// get kafka writer using environment variables.
		kafkaURL = os.Getenv("kafkaURL")
	}
	writer := getKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	fmt.Sprintf("start producing - %s", foundDocumentEvented)

	key := fmt.Sprintf("Key-%d", uuid.New())
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte("Es Cliente"),
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("produced at topic:%s = %s\n", topic, string(msg.Value))
	}
}

func consumeFoundDocumentEvented() {
	kafkaURL := ""
	topic := foundDocumentEvented
	groupID := ""
	mode, ok := viper.Get("LOCAL").(string)
	if ok && mode == "1" {
		kafkaURL = viper.Get("KAFKA.URL").(string)
		groupID = viper.Get("KAFKA.GROUPID").(string)
	} else {
		// get kafka writer using environment variables.
		kafkaURL = os.Getenv("kafkaURL")
		groupID = os.Getenv("groupID")
	}

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	fmt.Sprintf("start consuming - %s", foundDocumentEvented)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("consumed at topic:%v = %s\n", m.Topic, string(m.Value))
	}
}
