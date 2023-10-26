package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"syscall"

	"github.com/PuerkitoBio/goquery"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sanjivkannaa/CryptoCops/scrape/config"
	"github.com/sanjivkannaa/CryptoCops/scrape/schema"
)

func main() {
	// err := godotenv.Load("../../env/.env")
	// if err != nil {
	// 	log.Fatal("Error in loading the .env file: ", err)
	// }
	proxyPort := os.Getenv("TOR_PROXY_PORT")
	println("Tor proxy port: ", proxyPort)
	// Set the tor proxy
	proxy, err := url.Parse("http://host.docker.internal:" + proxyPort)
	if err != nil {
		log.Fatal("Unable to parse proxy url", err)
	}
	torTransport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	// No timeout
	httpClient := &http.Client{Transport: torTransport}
	log.Println("Tor proxy setup done!!!")
	config.InitDB()
	kafkaUrl := os.Getenv("KAFKA_URL")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaUrl,
	})
	if err != nil {
		log.Fatal("Failed to create kafka producer: ", err)
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        kafkaUrl,
		"group.id":                 "my-consumer-group3",
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": true,
	})
	if err != nil {
		log.Fatal("Failed to create kafka consumer: ", err)
	}
	Scrap(httpClient, producer, consumer)
}

func Scrap(httpClient *http.Client, producer *kafka.Producer, consumer *kafka.Consumer) {
	if a := recover(); a != nil {
		log.Println("Error in Scrap function", a)
	}
	topic := os.Getenv("KAFKA_TOPIC2")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal("Error in subscribing to the kafka topic: ", err)
	}
	run := true
	for run == true {
		select {
		case sig := <-signalChan:
			consumer.Close()
			fmt.Println("Terminating: ", sig)
			close(signalChan)
			return
		case event := <-consumer.Events():
			switch e := event.(type) {
			case *kafka.Message:
				var mes map[string]string
				err = json.Unmarshal(e.Value, &mes)
				if err != nil {
					log.Println("Error in unmarshalling data: ", err)
				}
				ScrapeUrl(mes["message"], producer, httpClient)
				continue
			default:
				run = true
			}
		default:
			run = true

		}
	}
}

func ScrapeUrl(url string, producer *kafka.Producer, httpClient *http.Client) {
	//topicCrawler := os.Getenv("KAFKA_TOPIC3")
	topicMlmodel := os.Getenv("KAFKA_TOPIC4")
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("Error getting data from the url ", url, err)
		ErroneousUrl(url)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Error in fetching the page", url, resp.StatusCode)
		ErroneousUrl(url)
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error parsing data from the url ", url, err)
		ErroneousUrl(url)
		return
	}
	// doc.Find("a").Each(func(i int, s *goquery.Selection) {
	// 	link, _ := s.Attr("href")
	// 	link = strings.Trim(link, " ")
	// 	mes := map[string]string{
	// 		"message": url,
	// 	}
	// 	jsonBytes, err := json.Marshal(mes)
	// 	if err != nil {
	// 		log.Println("Failed to marshal data: ", err)
	// 	}
	// 	err = producer.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &topicCrawler, Partition: kafka.PartitionAny},
	// 		Value:          jsonBytes,
	// 	}, nil)
	// 	if err != nil {
	// 		fmt.Println("Failed to produce messages: ", err)
	// 	}
	// 	fmt.Println(link)
	// })
	document := map[string](string){
		"url":     url,
		"message": doc.Text(),
	}
	jsonBytes, err := json.Marshal(document)
	if err != nil {
		log.Fatal("Failed to marshal data: ", err)
	}
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicMlmodel, Partition: kafka.PartitionAny},
		Value:          jsonBytes,
	}, nil)
	if err != nil {
		log.Fatal("Failed to produce message: ", err)
	}

}

func ErroneousUrl(url string) {
	db := config.Getdb()
	unscrapedUrl := db.Collection("UnscrapedUrl")
	data := schema.UnscrapedURL{
		URL: url,
	}
	// ctx := context.Background()
	res, err := unscrapedUrl.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal("Failed to insert data into database: ", err)
	} else {
		log.Println("Inserted data into database successfully: ", res)
	}
}
