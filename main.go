package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/turnage/graw/reddit"
)

var (
	Log      = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	bot      reddit.Bot
	posts    []*reddit.Post
	threads  []*reddit.Post
	comments []*reddit.Comment
	after    string
	keys     = make(map[string]bool)
	cfg      = reddit.BotConfig{
		Agent: "graw:saved_file_bot:0.0.1 by /u/DiscombobulatedAd988",
	}
	dbURL  string
	dbpool *pgxpool.Pool
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// If environmental variable is blank, read from .env file
	if os.Getenv(`BROKER_TASTYWORKS_USERNAME`) == `` {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			if err, ok := err.(viper.ConfigFileNotFoundError); ok {
				Log.Fatal().Err(err).Msgf(`init().viper.ReadInConfig()`)
			} else {
				Log.Fatal().Err(err).Msg(`init().viper.ReadInConfig(): Config file found but had an error`)
			}
		}

		Log.Info().Msg(`init().viper.ReadInConfig(): Config loaded successfully`)
	} else {
		// For Docker, read from environmental variables
		viper.AutomaticEnv()
	}

	dbURL = fmt.Sprintf("postgres://%s:%s@%s/%s",
		viper.GetString(`DB_USER`),
		viper.GetString(`DB_PASS`),
		viper.GetString(`DB_HOST`),
		viper.GetString(`DB_NAME`),
	)

	cfg.App = reddit.App{
		ID:       viper.GetString(`REDDIT_ID`),
		Secret:   viper.GetString(`REDDIT_SECRET`),
		Username: viper.GetString(`REDDIT_USERNAME`),
		Password: viper.GetString(`REDDIT_PASSWORD`),
	}

	var err error
	bot, err = reddit.NewBot(cfg)
	if err != nil {
		Log.Error().Err(err).Send()
	}

	// Setup database connection pool
	dbpool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		Log.Error().Err(err).Send()
		os.Exit(1)
	}
}

func main() {
	readFile()
	Start()
	defer dbpool.Close()
}
