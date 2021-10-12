package connect

type Token struct {
	Access    string   `json:"access_token"`
	Type      string   `json:"token_type"`
	ExpiresIn int      `json:"expires_in"`
	Refresh   string   `json:"refresh_token"`
	Scopes    []string `json:"scopes"`
	err       error
}

type UserData struct {
	Data Data `json:"data"`
	err  error
}

type Data struct {
	CID      string   `json:"cid"`
	Personal Personal `json:"personal,omitempty"`
	Vatsim   Vatsim   `json:"vatsim,omitempty"`
}

type Personal struct {
	NameFirst string  `json:"name_first,omitempty"`
	NameLast  string  `json:"name_last,omitempty"`
	NameFull  string  `json:"name_full,omitempty"`
	Email     string  `json:"email,omitempty"`
	Country   Country `json:"country,omitempty"`
}

type Country struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Vatsim struct {
	Rating      Rating      `json:"rating,omitempty"`
	PilotRating PilotRating `json:"pilotrating,omitempty"`
	Region      Region      `json:"region,omitempty"`
	Division    Division    `json:"division,omitempty"`
	Subdivision Subdivision `json:"subdivision,omitempty"`
}

type Rating struct {
	ID    int    `json:"id,omitempty"`
	Long  string `json:"long,omitempty"`
	Short string `json:"short,omitempty"`
}

type PilotRating struct {
	ID    int    `json:"id,omitempty"`
	Long  string `json:"long,omitempty"`
	Short string `json:"short,omitempty"`
}

type Division struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Subdivision struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Region struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
