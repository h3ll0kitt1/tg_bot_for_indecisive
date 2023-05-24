package config

import (
	"flag"
	"log"
)

type Config struct {
	AccessToken string
}

func Load() (Config, error) {
	Token := flag.String(
		"t",
		"",
		"token to access telegram bot",
	)

	flag.Parse()

	if *Token == "" {
		log.Fatal("token is not specified")
	}

	return Config{
		AccessToken: *Token,
	}, nil
}
