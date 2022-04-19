package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
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

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
