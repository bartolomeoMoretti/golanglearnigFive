package main

func main() {

    // token - для аутентификации в ТГ; получаем как значение флага при запуске программы
    // token = flags.Get(token)

    // tgClient необходим для общения с api ТГ
    // tgClient = telegram.New(token)

    // fetcher направляет в api ТГ запрос на наличие новых событий и затем их получение от ТГ
    // fetcher = fetcher.New()

    // processor после обработки будет отправлять нам новые сообщения
    // processor = processor.New()

    // consumer - получает события (при помощи fetcher и обрабатывает пр помощи processor)
	// consumer.Start(fetcher, processor)
}

func mustToken() string {
    
}