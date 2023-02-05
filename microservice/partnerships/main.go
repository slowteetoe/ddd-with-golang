package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Res struct {
	AvailableHotels []struct {
		Name               string `json:"name"`
		PriceInUSDPerNight int    `json:"priceInUSDPerNight"`
	} `json:"availableHotels"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10

	sampleRes := Res{AvailableHotels: []struct {
		Name               string `json:"name"`
		PriceInUSDPerNight int    `json:"priceInUSDPerNight"`
	}{
		{
			Name:               "some hotel",
			PriceInUSDPerNight: 300,
		},
		{
			Name:               "some other hotel",
			PriceInUSDPerNight: 30,
		},
		{
			Name:               "some third hotel",
			PriceInUSDPerNight: 90,
		},
		{
			Name:               "some fourth hotel",
			PriceInUSDPerNight: 80,
		},
	}}

	b, err := json.Marshal(sampleRes)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/partnerships",
		func(writer http.ResponseWriter, request *http.Request) {
			ran := rand.Intn(max - min + 1)
			if ran > 7 {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(b)
		})

	log.Info().Msg("running...")
	if err := http.ListenAndServe(":3031", r); err != nil {
		log.Fatal().Err(err).Send()
	}
}
