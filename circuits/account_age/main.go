package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/brevis-network/brevis-sdk/sdk/prover"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	proverService, err := prover.NewService(&LPFeesBNBCakeCircuit{}, prover.ServiceConfig{
		SetupDir: "$HOME/circuitOut",
		SrsDir:   "$HOME/kzgsrs",
	})
	if err != nil {
		log.Error().Err(err).Msg("could not create prover service")
		os.Exit(1)
	}
	const port uint = 32001
	log.Info().Msgf("starting prover service on port: %d", port)
	proverService.Serve("0.0.0.0", port)
}
