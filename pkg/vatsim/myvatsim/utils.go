package myvatsim

import (
	"api/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const myVATSIMEndpoint = "v2/events/view/division/EUD"

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
