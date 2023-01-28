package telegram

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "path"
    "strconv"
)

//клиент для api-бота: читаем сообщения для бота, отправляем сообщения от бота

type Client struct {
    // host - хост api-сервиса ТГ
    // basePath - префикс (идущий после хоста), с которого начинаются все запросы
    // client - чтобы не создавать его отдельно для каждого запроса
    host     string
    basePath string
    client   http.Client
}

func New(host string, token string) *Client {
    //функция New тут создаёт клиента
    return &Client{
        host:     host,
        basePath: newBasePath(token),
        client:   http.Client{},
    }
}

func newBasePath(token string) string {
    return "bot" + token
}

// получение апдейтов
// offset, limit - параметры для getUpdates
// limit - кол-во апдейтов, которые получаем за 1 запрос
/*offset - смещение: с какого апдейта забирать очередную пачку апдейтов;
нпр., ранее получили 15 апдейтов (limit 15), хотим дальше получить поочереди и тогда подаём offset 16*/

func (c *Client) Updates(offset, limit int) ([]Update, error) {
    // функция возвращает структуру, содержащую всё, что нужно знать об апдейте
    // посколько нужен ни один тип, вынесено в отдельный файл в: telegram/types.go
    //формируем параметры запроса
    q := url.Values{}
    // url.Values.Add - добавляет указанный параметр к запросу
    q.Add("offset", strconv.Itoa(offset))
    q.Add("limit", strconv.Itoa(limit))
    //отправляем запрос
    data, err := c.doRequest("getUpdates", q)

    if err != nil {
        return nil, err
    }
    // теперь мы получили данные из ответа. Что нам с ними делать?
    // понимаем, что там json, а поэтому распарсим его

    var res UpdatesResponse
    // распарсинг json
    // в 1-м аргументе указываем что именно парсим , а во 2-м показываем куда
    if err := json.Unmarshal(data, &res); err!=nil{
        return nil,err
    }

    return res.Result, nil
}

// SendMessage отправляем сообщения
func (c *Client) SendMessage(chatID int, text string) error {
    q := url.Values{}
    q.Add("chat_id", strconv.Itoa(chatID))
    q.Add("text", text)

    _, err := c.doRequest("sendMessage", q)
    if err != nil {
        return fmt.Errorf("can't do request: %w", err)
    }

    return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
    // на вход подаём методв в виде аргумента + подаем запрос, а в ответ возвращает байты (которые получили в ответ) и ошибку
    // формируем url, на который будет отправляеться запрос
    u := url.URL{
        // указываем протокол
        Scheme: "https",
        // хост
        Host: c.host,
        // Путь, который состоит из двух частей = базовая часть + нужный метод
        // path.Join - аналог: c.basePath + "/" + method
        Path: path.Join(c.basePath, method),
    }
    // формируем объект запроса - пока не отправляем его, а подготавливаем
    // подготовка: передаём название метода передаём в NewRequest
    /*+ далее передаём url в текстовом виде (тип url реализует интерфейс Стрингер, что говорит о том,
      что у него есть обязательный метод Стринг, который будет возвращать то, что нужно)
      + далее последним аргументом указывается тело запроса: при Get обычно nil*/
    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("can't do request: %w", err)
        // почитать про errors.Is() и errors.As()
    }
    // передаем в объект req параметры запроса, которые получены в аргументе query url.Values
    // метод Encode() приведёт параметры к тому виду, который подходит для отправки на сервер
    req.URL.RawQuery = query.Encode()
    // отправляем получившийся запрос
    // для отправки используем тот клиент, который заранее приготовили
    // используем метод Do(), в который мы передаём объект запроса
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("can't do request: %w", err)
    }
    // получаем содержимое
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("can't do request: %w", err)
    }

    return body, nil
    // возвращаемся в метод Updates
}
