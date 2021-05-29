package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const OtsPrefix string = "ots:"

type SecretsPayload struct {
	Message  string `json:"message,omitempty"`
	Password string `json:"password"`
}

// pingRedis verifies a successful connection with redis
func pingRedis(c redis.Conn) error {
	_, err := redis.String(c.Do("PING"))
	return err
}

// Ready handles the GET /ready API that returns the state of the
// service based on the connectivity with the REDIS database
func Ready(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request from %+v to check the readiness of the API server...", r.Host)
	conn := pool.Get()
	defer conn.Close()
	if err := pingRedis(conn); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `{"isReady": false}`)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"isReady": true}`)
	}
}

// CreateSecret handles the POST /secrets API that creates a new
// one-time secret and returns the secret UUID
func CreateSecret(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request from %+v to create a new secret...", r.Host)
	secretID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(r.Body)
	var payload SecretsPayload
	err = decoder.Decode(&payload)
	if err != nil {
		log.Fatal(err)
	}
	conn := pool.Get()
	defer conn.Close()
	ciphertext := encrypt(payload.Message, payload.Password)
	_, err = conn.Do("SET", OtsPrefix+secretID.String(), ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"uuid": `+secretID.String()+`}`)
}

// ReadSecret handles the GET /secrets/{secretId} API handler. It accepts
// the secret ID and the password associated with it and returns the secret
func ReadSecret(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request from %+v to read a secret...", r.Host)
	params := mux.Vars(r)
	secretID := params["secretID"]

	decoder := json.NewDecoder(r.Body)
	var payload SecretsPayload
	err := decoder.Decode(&payload)
	if err != nil {
		log.Fatal(err)
	}

	conn := pool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", OtsPrefix+secretID)
	if err != nil || reply == nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := `{"error": "Could not fetch the secret ` + secretID + `"}`
		fmt.Fprint(w, errMsg)
		log.Println(errMsg)
		return
	}
	plaintext, err := decrypt(string(reply.([]byte)), payload.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := `{"error": "Could not fetch the secret ` + secretID + `"}`
		fmt.Fprint(w, errMsg)
		log.Printf("SecretID: %s. Error: %+v", secretID, err)
		return
	}
	_, err = conn.Do("DEL", OtsPrefix+secretID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := `{"error": "Could not fetch the secret ` + secretID + `"}`
		fmt.Fprint(w, errMsg)
		log.Printf("SecretID: %s. Error: %+v", secretID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message": "`+plaintext+`"}`)
}
