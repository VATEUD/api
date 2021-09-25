package myvatsim

import (
	"api/pkg/cache"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AllEvents(w http.ResponseWriter, r *http.Request) {

	val, err := cache.RedisCache.Get("EVENTS_ALL")

	if err != nil && err != redis.Nil {
		log.Printf("Error occurred while fetching events from cache. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if len(val) > 0 {
		if _, err := w.Write([]byte(val)); err != nil {
			log.Println("Error writing the response.")
		}
		return
	}

	eventsChannel := make(chan EventsFeed)
	go getEvents(eventsChannel)
	events := <-eventsChannel

	if events.err != nil {
		log.Printf("Error occurred while fetching events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(events.Data)

	if err != nil {
		log.Printf("Error occurred while marshalling events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := cache.RedisCache.Set("EVENTS_ALL", string(bytes), 2 * time.Minute); err != nil {
		log.Println("Error saving events to cache. Error:", err.Error())
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing the response.")
	}
}

func EventsByAmount(w http.ResponseWriter, r *http.Request) {

	attrs := mux.Vars(r)

	val, err := cache.RedisCache.Get(fmt.Sprintf("EVENTS_%s", attrs["amount"]))

	if err != nil && err != redis.Nil {
		log.Printf("Error occurred while fetching events from cache. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if len(val) > 0 {
		if _, err := w.Write([]byte(val)); err != nil {
			log.Println("Error writing the response.")
		}
		return
	}

	eventsChannel := make(chan EventsFeed)
	go getEvents(eventsChannel)
	events := <-eventsChannel

	if events.err != nil {
		log.Printf("Error occurred while fetching events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	amount, err := strconv.Atoi(attrs["amount"])

	if err != nil {
		log.Printf("Error occurred while converting the amont. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(events.Data[:amount])

	if err != nil {
		log.Printf("Error occurred while marshalling events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := cache.RedisCache.Set(fmt.Sprintf("EVENTS_%s", attrs["amount"]), string(bytes), 2 * time.Minute); err != nil {
		log.Println("Error saving events to cache. Error:", err.Error())
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing the response.")
	}
}
