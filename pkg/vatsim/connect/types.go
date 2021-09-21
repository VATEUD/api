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
	Personal Personal `json:"personal"`
	Vatsim   Vatsim   `json:"vatsim"`
}

type Personal struct {
	NameFirst string  `json:"name_first"`
	NameLast  string  `json:"name_last"`
	NameFull  string  `json:"name_full"`
	Email     string  `json:"email"`
	Country   Country `json:"country"`
}

type Country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Vatsim struct {
	Rating      Rating      `json:"rating"`
	PilotRating PilotRating `json:"pilotrating"`
	Region      Region      `json:"region"`
	Division    Division    `json:"division"`
	Subdivision Subdivision `json:"subdivision"`
}

type Rating struct {
	ID    int    `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
}

type PilotRating struct {
	ID    int    `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
}

type Division struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Subdivision struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
