package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Response struct {
	Writer  http.ResponseWriter `json:"-" xml:"-"`
	Request *http.Request       `json:"-" xml:"-"`
	Message string              `json:"message" xml:"message"`
	Status  int                 `json:"-" xml:"-"`
}

func New(w http.ResponseWriter, r *http.Request, message string, status int) Response {
	return Response{
		Writer:  w,
		Request: r,
		Message: message,
		Status:  status,
	}
}

func (response Response) Process() {
	if response.Request.Header.Get("Content-Type") == "application/xml" || response.Request.Header.Get("Content-Type") == "text/xml" {
		response.Writer.Header().Set("Content-Type", "application/xml")
		response.Writer.WriteHeader(response.Status)
		data, err := xml.Marshal(response)
		if err != nil {
			response.Writer.WriteHeader(http.StatusInternalServerError)
			response.Writer.Write([]byte(err.Error()))
			return
		}
		response.Writer.Write(data)
		return
	}

	response.Writer.Header().Set("Content-Type", "application/json")
	response.Writer.WriteHeader(response.Status)
	data, err := json.Marshal(response)
	if err != nil {
		response.Writer.WriteHeader(http.StatusInternalServerError)
		response.Writer.Write([]byte(err.Error()))
		return
	}
	response.Writer.Write(data)
}
