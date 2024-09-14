package main

import (
	"context"

	database "github.com/DEVunderdog/user-management-gRPC/database/sqlc"
	"github.com/DEVunderdog/user-management-gRPC/token"
	"github.com/DEVunderdog/user-management-gRPC/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	config, err := utils.LoadConfig("app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to the database")
	}

	store := database.NewStore(connPool)

	err = token.InitializeJWTKeys(config.Passphrase, store, ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

}