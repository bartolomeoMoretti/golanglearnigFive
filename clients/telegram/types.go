package telegram

// UpdatesResponse в случае, если получено ok=true, то в теле ответа приходит result как слайс из апдейтов
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update тут все типы, с которым работает клиент
type Update struct {
	// ID - это update_id в апи ТГ: Getting updates -> Update
	// Message - это message в апи ТГ: Getting updates -> Update
	// определённые ниже json-теги нужны, чтобы парсить json в эту структуру
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text     string `json:"text"`
	From     From   `json:"from"`
	Chat     Chat   `json:"chat"`
	DateSent int    `json:"date"`
}

type From struct {
	Username  string `json:"username"`
	UserID    int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Chat struct {
	ID int `json:"id"`
}
