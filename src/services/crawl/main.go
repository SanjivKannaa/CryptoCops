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
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sanjivkannaa/CryptoCops/crawl/config"
	"github.com/sanjivkannaa/CryptoCops/crawl/helpers"
	"github.com/sanjivkannaa/CryptoCops/crawl/schema"
)

type MessageData struct {
	Message string `json:"message"`
}

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
		"group.id":                 "my-consumer-group1",
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": true,
	})
	if err != nil {
		log.Fatal("Failed to create kafka consumer: ", err)
	}

	log.Println(httpClient)
	log.Println(producer)
	log.Println(consumer)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// to crawl seed url
	go func() {
		defer wg.Done()

		// Get the seed urls
		urls, err := helpers.Read("url_queue.json")
		if err != nil {
			log.Println(err)
		}
		for len(urls) > 0 {
			cur := urls[0]
			// Remove the current url from the queue
			urls = urls[1:]
			data := map[string]string{
				"message": cur.Url,
			}
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				log.Println("Failed to create json message: ", err)
			}
			topic := os.Getenv("KAFKA_TOPIC1")
			// for testing API
			// time.Sleep(1 * time.Second)
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          jsonBytes,
			}, nil)
			if err != nil {
				log.Fatal("Failed to produce kafka message: ", err)
			}
		}
	}()
	go func(httpClient *http.Client, producer *kafka.Producer, consumer *kafka.Consumer) {
		defer wg.Done()
		msg_count := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				//msg, err := consumer.ReadMessage(-1)
				if err != nil {
					log.Println("Error reading message:", err)
					break
				}

				Crawl(httpClient, producer, consumer, msg_count)
			}
		}
	}(httpClient, producer, consumer)
	waitForTerminationSignal(cancel, wg)
	//Crawl(httpClient, producer, consumer)
}

func Crawl(httpClient *http.Client, producer *kafka.Producer, consumer *kafka.Consumer, msg_count int) {
	if a := recover(); a != nil {
		log.Println("Error in crawling function", a)
		return
	}

	// check if a visited url file exists
	// _, err := os.Stat("visited.json")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		os.Create("visited.json")
	// 	}
	// }

	// // Get the seed urls
	// urls, err := helpers.Read("url_queue.json")
	// if err != nil {
	// 	log.Println(err)
	// }

	// // adding seed urls(from url_queue.json) to KAFKA_TOPIC1
	// for len(urls) > 0 {
	// 	cur := urls[0]
	// 	// Remove the current url from the queue
	// 	urls = urls[1:]
	// 	data := map[string]string{
	// 		"message": cur.Url,
	// 	}
	// 	jsonBytes, err := json.Marshal(data)
	// 	if err != nil {
	// 		log.Println("Failed to create json message: ", err)
	// 	}
	// 	topic := os.Getenv("KAFKA_TOPIC1")
	// 	err = producer.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 		Value:          jsonBytes,
	// 	}, nil)
	// 	if err != nil {
	// 		log.Fatal("Failed to produce kafka message: ", err)
	// 	}
	// links, err := Scrap(httpClient, cur.Url)
	// if err != nil {
	// 	log.Println("Error scraping the url => ", cur.Url, " : ", err)
	// 	continue
	// }

	// Update the visited urls
	// visited, err := helpers.Read("visited.json")
	// if err != nil {
	// 	log.Println(err)
	// }
	// visited = append(visited, cur)
	// err = helpers.Write(visited, "visited.json")
	// if err != nil {
	// 	log.Println(err)
	// }

	// Update the url queue
	// for _, i := range links {
	// 	if i != "" {
	// 		urls = append(urls, helpers.Link{Url: i})
	// 	}
	// }
	// err = helpers.Write(urls, "url_queue.json")
	// if err != nil {
	// 	log.Println(err)
	// }
	//}
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	topic1 := os.Getenv("KAFKA_TOPIC3")
	err := consumer.SubscribeTopics([]string{topic1}, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to kafka topic: ", err)
	}
	msg, err := consumer.ReadMessage(-1)
	if err != nil {
		log.Println("Error reading message:", err)
	}
	msg_count += 1
	// if msg_count % MIN_COMMIT_COUNT == 0 {
	// 	go func() {
	// 		consumer.Commit()
	// 	}()
	// }
	// fmt.Printf("%% Message on %s:\n%s\n",
	// 	e.TopicPartition, string(e.Value))
	var data MessageData
	err = json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Println("Failed to parse JSON message:", err)
		return
	}
	message := data.Message
	fmt.Println("recieved message: ", message)
	crawlUrl(message, producer, httpClient)
	// run := true
	// go func() {
	// for run == true {
	// 	select {
	// 	case sig := <-signalChan:
	// 		consumer.Close()
	// 		fmt.Println("Terminating: ", sig)
	// 		close(signalChan)
	// 		return
	// 	case event := <-consumer.Events():
	// 		switch e := event.(type) {
	// 		case *kafka.Message:
	// 			var mes map[string]string
	// 			err = json.Unmarshal(e.Value, &mes)
	// 			fmt.Println("The message is: ", mes["message"])
	// 			topic := os.Getenv("KAFKA_TOPIC1")
	// 			err = producer.Produce(&kafka.Message{
	// 				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 				Value:          e.Value,
	// 			}, nil)
	// 			if err != nil {
	// 				log.Fatal("Error in producing messages: ", err)
	// 			}
	// 		default:
	// 			// fmt.Println("Error in consuming messages: ", e)
	// 			run = true
	// 		}
	// 	default:
	// 		// fmt.Println("gonna terminate")
	// 		run = true

	// 	}
	// }
	// }()

}

func waitForTerminationSignal(cancel context.CancelFunc, wg *sync.WaitGroup) {
	// err := redisClient.Close()
	// if err != nil {
	// log.Fatalf("Failed to close Redis client: %v", err)
	// }
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)
	<-signalChan
	log.Println("Termination signal received. Shutting down...")
	cancel()
	wg.Wait()
	os.Exit(0)
}

func crawlUrl(url string, producer *kafka.Producer, httpClient *http.Client) {
	topicMiddleware := os.Getenv("KAFKA_TOPIC1")
	resp, err := httpClient.Get(url)
	//log.Println(url)
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
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		link = strings.Trim(link, " ")
		mes := map[string]string{
			"message": url,
		}
		jsonBytes, err := json.Marshal(mes)
		if err != nil {
			log.Println("Failed to marshal data: ", err)
		}
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicMiddleware, Partition: kafka.PartitionAny},
			Value:          jsonBytes,
		}, nil)
		if err != nil {
			fmt.Println("Failed to produce messages: ", err)
		}
		fmt.Println(link)
	})
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
