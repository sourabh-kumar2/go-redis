package main

import (
	"flag"
	"log"

	"github.com/sourabh-kumar2/go-redis/config"
	"github.com/sourabh-kumar2/go-redis/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", config.DEFAULT_HOST, "host for the server")
	flag.IntVar(&config.Port, "port", config.DEFAULT_PORT, "port for the server")
	flag.Parse()
}

func main() {
	setupFlags()

	log.Println("rolling the server")

	if err := server.RunTCPAsyncServer(); err != nil {
		log.Fatal(err)
	}
}
