package myvatsim

import (
	"api/pkg/cache"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

const myVATSIMEndpoint = "v2/events/view/division/EUD"

func Events(w http.ResponseWriter, r *http.Request) {

	attrs := mux.Vars(r)

	val, err := cache.RedisCache.Get(cacheKey(attrs))

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

	if err := cache.RedisCache.Set(cacheKey(attrs), string(bytes), 2 * time.Minute); err != nil {
		log.Println("Error saving events to cache. Error:", err.Error())
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing the response.")
	}
}

func getEvents(eventsChannel chan EventsFeed) {
	client := &http.Client{}
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", utils.Getenv("MYVATSIM_API_URL", "https://my.vatsim.net/api"), myVATSIMEndpoint), nil)

	if err != nil {
		eventsChannel <- EventsFeed{err: err}
		return
	}

	res, err := client.Do(r)

	if err != nil {
		eventsChannel <- EventsFeed{err: err}
		return
	}

	defer func(closer io.Closer) {
		if err := closer.Close(); err != nil {
			log.Printf("Error occurred while closing the response body. Error: %s.", err.Error())
		}
	}(res.Body)

	var eventsFeed EventsFeed

	if err := json.NewDecoder(res.Body).Decode(&eventsFeed); err != nil {
		eventsChannel <- EventsFeed{err: err}
		return
	}

	eventsChannel <- eventsFeed
}

func cacheKey(attrs map[string]string) string {

	amount, ok := attrs["amount"]

	if !ok || len(amount) < 1 {
		return "EVENTS_ALL"
	}

	return fmt.Sprintf("EVENETS_%s", amount)
}
