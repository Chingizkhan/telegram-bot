package main

import (
	"context"
	"flag"
	"log"
	tgClient "telegram-bot/clients/telegram"
	"telegram-bot/consumer/event_consumer"
	"telegram-bot/events/telegram"
	"telegram-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	fileStoragePath   = "files_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	log.Println("bot starting")

	//s := files.New(fileStoragePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalln("can't connect to storage: ", err)
	}

	if err = s.Init(context.TODO()); err != nil {
		log.Fatalln("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalln("service is stopped", err)
	}
}

func mustToken() string {
	t := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *t == "" {
		log.Fatalln("token is not specified")
	}

	return *t
}
