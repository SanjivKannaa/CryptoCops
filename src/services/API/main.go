package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	port := ":" + os.Getenv(("API_PORT"))
	http.HandleFunc("/api", handleAPIRequest)
	http.HandleFunc("/api/pause", API_pause)
	http.HandleFunc("/api/resume", API_resume)
	fmt.Println("Server is running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "its working!!",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func API_pause(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Paused!",
	}
	Pause()
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func API_resume(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "its working!!",
	}
	Resume()
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func Pause() {
	kafka_url := os.Getenv("KAFKA_URL")
	config := kafka.ConfigMap{
		"bootstrap.servers": kafka_url,
		"group.id":          "my-consumer-group2",
		"auto.offset.reset": "earliest",
		// "auto.commit.interval.ms": 8000,
		//		"enable.auto.commit": "false",
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}

	topic := os.Getenv("KAFKA_TOPIC1")
	// topic := "to_middleware"
	fmt.Println("listening to topic:", topic)
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic:", err)
	}
	consumer.Pause([]kafka.TopicPartition{})
}

func Resume() {
	kafka_url := os.Getenv("KAFKA_URL")
	config := kafka.ConfigMap{
		"bootstrap.servers": kafka_url,
		"group.id":          "my-consumer-group2",
		"auto.offset.reset": "earliest",
		// "auto.commit.interval.ms": 8000,
		//		"enable.auto.commit": "false",
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}

	topic := os.Getenv("KAFKA_TOPIC1")
	// topic := "to_middleware"
	fmt.Println("listening to topic:", topic)
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic:", err)
	}
	consumer.Resume([]kafka.TopicPartition{})
}
