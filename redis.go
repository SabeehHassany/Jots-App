package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

// Define the Redis channel for new jots notifications
const newJotsChannel = "new_jots_channel"

func init() {
	redisURL := os.Getenv("redis://default:qxTNvTjxSCBiyIyYcedBBoRCslvdulvl@redis.railway.internal:6379")
	if redisURL == "" {
		log.Fatalf("REDIS_URL not set")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL, // Use the Redis URL from the environment variable
		Password: "",       // No password set by default, adjust if needed
		DB:       0,        // Use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}

/*
package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

// Define the Redis channel for new jots notifications
const newJotsChannel = "new_jots_channel"

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}
*/
