package main

import (
    "golanglearningFive/cfg/s"
    tgClient "golanglearningFive/clients/telegram"
    eventConsumer "golanglearningFive/consumer/event-consumer"
    "golanglearningFive/events/telegram"
    "golanglearningFive/storage/files"
    "log"
)

const (
    tgBotHost   = "api.telegram.org"
    storagePath = "storage/users-files"
    batchSize = 100
)

func main() {

    // token - для аутентификации в ТГ; получаем как значение флага при запуске программы
    // token = flags.Get(token)
    t := mustToken()

    // tgClient необходим для общения с api ТГ
    // tgClient = telegram.New(token)

    eventsProcessor := telegram.New(
        tgClient.New(tgBotHost, t),
        files.New(storagePath),
    )

    log.Print("server has been started")

    consumer := eventConsumer.New(eventsProcessor,eventsProcessor,batchSize)

    if err := consumer.Start();err!=nil{
        log.Fatal("server was stopped", err)
    }

    // fetcher направляет в api ТГ запрос на наличие новых событий и затем отвечает за их получение от ТГ
    // fetcher = fetcher.New()

    // processor после обработки будет отправлять нам новые сообщения
    // processor = processor.New()

    // consumer - получает события (при помощи fetcher и обрабатывает пр помощи processor)
	// consumer.Start(fetcher, processor)
}

func mustToken() string {
    token := s.T /*flag.String(
		"tg-token",
		"",
		"token for access",
	)
    
    flag.Parse()*/

    /*if *token == "" {
        log.Fatal("token is empty")
    }*/

    return token
}