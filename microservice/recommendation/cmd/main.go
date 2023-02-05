package main

import (
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"slowteetoe.com/recommentations/recommendation/internal/recommendation"
	"slowteetoe.com/recommentations/recommendation/internal/transport"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	c := retryablehttp.NewClient()
	c.RetryMax = 3

	partnerAdapter, err := recommendation.NewPartnershipAdapter(
		c.StandardClient(), "http://localhost:3031",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a partnerAdapter: ")
	}
	svc, err := recommendation.NewService(partnerAdapter)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a service: ")
	}
	handler, err := recommendation.NewHandler(*svc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a service: ")
	}
	m := transport.NewMux(*handler)
	if err := http.ListenAndServe(":4040", m); err != nil {
		log.Fatal().Err(err).Msg("server errored: ")
	}

}
