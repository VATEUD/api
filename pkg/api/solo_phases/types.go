package solo_phases

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

var (
	userIDNotProvided = errors.New("user ID not provided")
	validUntilNotProvided = errors.New("expiry date not provided")
)

type soloPhaseRequest struct {
	User       int
	Position   string
	ValidUntil time.Time
	Extensions int8
}

type soloPhaseSaveResult struct {
	soloPhase models.SoloPhase
	err error
}

type validator interface {
	validate() error
}

func validate(request validator) error {
	return request.validate()
}

func newSoloPhaseRequest(r *http.Request) (*soloPhaseRequest, error) {
	if r == nil {
		return nil, errors.New("request can not be nil")
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if len(r.PostForm.Get("user_id")) < 1 {
		return nil, userIDNotProvided
	}

	if len(r.PostForm.Get("valid_until")) < 1 {
		return nil, validUntilNotProvided
	}

	userID, err := strconv.Atoi(r.PostForm.Get("user_id"))

	if err != nil {
		return nil, err
	}

	t, err := time.Parse(dateFormat, r.PostForm.Get("valid_until"))

	if err != nil {
		return nil, err
	}

	return &soloPhaseRequest{
		User:       userID,
		Position:   r.PostForm.Get("position"),
		ValidUntil: t,
	}, nil
}

func (req soloPhaseRequest) validate() error {
	if req.User < 80000 || req.User > 2000000 {
		return errors.New("invalid CID provided")
	}

	if len(req.Position) < 3 || len(req.Position) > 12 || !strings.Contains(req.Position, "_") {
		return errors.New("invalid position provided")
	}

	if !req.ValidUntil.After(time.Now().UTC()) {
		return errors.New("invalid valid date provided")
	}

	return nil
}

func (req soloPhaseRequest) save(channel chan soloPhaseSaveResult, subdivision models.Subdivision) {
	solo := models.SoloPhase{
		UserID:        uint(req.User),
		Position:      req.Position,
		ValidUntil:    req.ValidUntil,
		Extensions:    0,
		Expired:       false,
		SubdivisionID: subdivision.ID,
		Subdivision:   &subdivision,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := database.DB.Create(&solo).Error; err != nil {
		channel <- soloPhaseSaveResult{err: err}
		return
	}

	channel <- soloPhaseSaveResult{soloPhase: solo}
}
