package myvatsim

// EventsFeed represents the myVATSIM datafeed structure
type EventsFeed struct {
	Data []Event `json:"data"`
	err  error
}

// Events represents the event JSON structure
type Event struct {
	ID               int         `json:"id"`
	Type             string      `json:"type"`
	VsoName          string      `json:"vso_name"`
	Name             string      `json:"name"`
	Link             string      `json:"link"`
	Organisers       []Organiser `json:"organisers"`
	Airports         []Airport   `json:"airports"`
	Routes           []Route     `json:"routes"`
	StartTime        string      `json:"start_time"`
	EndTime          string      `json:"end_time"`
	ShortDescription string      `json:"short_description"`
	Description      string      `json:"description"`
	Banner           string      `json:"banner"`
}

// Organiser represents the event organiser structure
type Organiser struct {
	Region            string `json:"region"`
	Division          string `json:"division"`
	Subdivision       string `json:"subdivision"`
	OrganisedByVatsim bool   `json:"organised_by_vatsim"`
}

// Airport represents the event airport structure
type Airport struct {
	ICAO string `json:"icao"`
}

// Route represents the event route structure
type Route struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Route     string `json:"route"`
}
