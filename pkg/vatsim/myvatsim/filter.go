package myvatsim

import (
	"api/internal/pkg/logger"
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

func EventsFilterDays(w http.ResponseWriter, r *http.Request) {

	utils.Allow(w, "*")
	attrs := mux.Vars(r)

	val, err := cache.RedisCache.Get(fmt.Sprintf("EVENTS_DAYS_%s", attrs["days"]))

	if err != nil && err != redis.Nil {
		logger.Log.Errorf("Error occurred while fetching events from cache. Error: %s.\n", err.Error())
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
		logger.Log.Errorf("Error occurred while fetching events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	days, err := strconv.Atoi(attrs["days"])

	if err != nil {
		logger.Log.Errorf("Error occurred while converting days. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(events.FilterDays(uint(days)))

	if err != nil {
		logger.Log.Errorf("Error occurred while marshalling events. Error: %s.\n", events.err.Error())
		res := response.New(w, r, "Internal error occurred while fetching events.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := cache.RedisCache.Set(fmt.Sprintf("EVENTS_DAYS_%s", attrs["days"]), string(bytes), 2*time.Minute); err != nil {
		logger.Log.Errorln("Error saving events to cache. Error:", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		logger.Log.Errorln("Error writing the response.")
	}
}
