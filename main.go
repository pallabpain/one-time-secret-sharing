package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var pool *redis.Pool

func main() {
	redisUrl := getEnv("REDIS_URL", "localhost:6379")

	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisUrl)
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ready", Ready)
	router.HandleFunc("/secrets", CreateSecret).Methods("POST")
	router.HandleFunc("/secrets/{secretID}", ReadSecret).Methods("POST")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatalf("%+v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
