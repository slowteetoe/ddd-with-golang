package recommendation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Rhymond/go-money"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) (*Handler, error) {
	if svc == (Service{}) {
		return nil, errors.New("service cannot be empty")
	}
	return &Handler{svc: svc}, nil
}

type GetRecommendationResponse struct {
	HotelName string `json:"hotelName"`
	TotalCost struct {
		Cost     int64  `json:"cost"`
		Currency string `json:"currency"`
	} `json:"totalCost"`
}

func (h Handler) GetRecommendation(w http.ResponseWriter, req *http.Request) {
	// caution: Get returns the first value only!
	location := req.URL.Query().Get("location")
	if location == "" {
		log.Info().Msg("missing 'location' param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	from := req.URL.Query().Get("from")
	if from == "" {
		log.Info().Msg("missing 'from' param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	to := req.URL.Query().Get("to")
	if to == "" {
		log.Info().Msg("missing 'to' param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	budget := req.URL.Query().Get("budget")
	if budget == "" {
		log.Info().Msg("missing 'budget' param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	const expectedFormat = "2006-01-02"

	formattedStart, err := time.Parse(expectedFormat, from)
	if err != nil {
		log.Info().Msgf("invalid 'from' date: %s %w", from, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	formattedEnd, err := time.Parse(expectedFormat, to)
	if err != nil {
		log.Info().Msgf("invalid 'to' date: %s %w", to, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := strconv.ParseInt(budget, 10, 64)
	if err != nil {
		log.Info().Msgf("invalid 'budget' param: %s %w", budget, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	budgetMon := money.New(b, "USD")

	rec, err := h.svc.Get(req.Context(), formattedStart, formattedEnd, location, *budgetMon)
	if err != nil {
		log.Info().Msgf("unable to pull recommendations: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(GetRecommendationResponse{
		HotelName: rec.HotelName,
		TotalCost: struct {
			Cost     int64  "json:\"cost\""
			Currency string "json:\"currency\""
		}{
			Cost:     rec.TripPrice.Amount(),
			Currency: "USD",
		},
	})
	if err != nil {
		log.Info().Err(err).Msg("something went wrong: ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
