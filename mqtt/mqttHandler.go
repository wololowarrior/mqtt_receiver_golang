package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var redisClient *redis.Client

const redisSpeedKey = "speed_new"

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
		DB:   0,
	})
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	var data map[string]interface{}
	if err := json.Unmarshal(message.Payload(), &data); err != nil {
		log.Printf("Error decoding MQTT message: %v", err)
		return
	}

	// Extract speed value
	speed, ok := data["speed"].(float64)
	if !ok {
		log.Println("Invalid speed value in MQTT message")
		return
	}

	// Update Redis with the latest speed
	if err := redisClient.Set(context.Background(), redisSpeedKey, speed, 0).Err(); err != nil {
		log.Printf("Error updating Redis: %v", err)
	}
}

func main() {
	brokerHost := os.Getenv("MQTT_BROKER_HOST")
	topic := os.Getenv("MQTT_TOPIC")
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", brokerHost))
	opts.SetClientID("redis-client")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}

	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to MQTT topic: %v", token.Error())
	}

	fmt.Println("MQTT to Redis is running.")
	// enter in a loop and not terminate. To allow the program to receive as many requests.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan

	client.Disconnect(250)
	fmt.Println("MQTT to Redis bridge has been terminated.")
}
