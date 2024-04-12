package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL       string
	ClasterName string
	ClientID    string
	NatsIP      string
}

func ReadConfig() Config {
	var conf Config

	env, err := godotenv.Read("../../conf.env")
	if err != nil {
		log.Fatal("cant read config file")
	}
	conf.DbURL = env["dbURL"]
	conf.ClasterName = env["clasterName"]
	conf.ClientID = env["CLientID"]
	conf.NatsIP = env["NatsIP"]
	return conf
}
