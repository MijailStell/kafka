package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
		Person{Document: "45531746"},
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

func searchDocumentEventHandler() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalln(err)
		}

		kafkaWriter := getKafkaWriter(kafkaURL, searchDocumentEvented)

		defer kafkaWriter.Close()
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
			Value: body,
		}
		err = kafkaWriter.WriteMessages(req.Context(), msg)

		if err != nil {
			wrt.Write([]byte(err.Error()))
			log.Fatalln(err)
		}

		consumeFoundDocumentEvented()

		response := <-channel // receive from channel
		fmt.Println(response)
		wrt.Write([]byte(response))
	})
}

func main() {
	setupConfig()
	// kafkaWriter := getKafkaWriter(kafkaURL, searchDocumentEvented)

	// defer kafkaWriter.Close()
	// go consumeSearchDocumentEvented()
	// go consumeFoundDocumentEvented()

	// Add handle func for producer.
	http.HandleFunc("/api/v1/account/searchDocumentEvent", searchDocumentEventHandler())

	// Run the web server.
	fmt.Println("start producer-api ... !!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func consumeFoundDocumentEvented() {
	reader := getKafkaReader(kafkaURL, foundDocumentEvented, groupID)

	defer reader.Close()

	fmt.Sprintf("start consuming - %s", foundDocumentEvented)
	// for {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("consumed at topic:%s = %s\n", m.Topic, string(m.Value))
	channel <- string(m.Value) // send to c
	// }
}
