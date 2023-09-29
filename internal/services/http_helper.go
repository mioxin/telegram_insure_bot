package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HttpHelper предоставляет удобный интерфейс
// для выполнения HTTP-запросов
type HttpHelper struct {
	client                   *http.Client
	rbody                    any
	uri                      string
	form                     map[string]string
	cookies, headers, params []pair
}
type pair struct {
	k, v string
}

// NewHttpHelper создает новый экземпляр HttpHelper
func NewHttpHelper() *HttpHelper {
	//req.ProxyHeader.Set("User-Agent", "Wget/1.9.1")
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 10
	transport.GetProxyConnectHeader = func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error) {
		return http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"}}, nil
	}
	//transport.DialContext =
	cl := http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: transport,
	}
	uri := ""
	return &HttpHelper{&cl, nil, uri, make(map[string]string), make([]pair, 0), make([]pair, 0), make([]pair, 0)}
}

// URL устанавливает URL, на который пойдет запрос
func (h *HttpHelper) URL(uri string) *HttpHelper {
	h.uri = uri
	return h
}

// Client устанавливает HTTP-клиента
// вместо умолчательного http.DefaultClient
func (h *HttpHelper) Client(cl *http.Client) *HttpHelper {
	h.client = cl
	return h
}

// Header устанавливает значение заголовка
func (h *HttpHelper) Header(key, value string) *HttpHelper {
	h.headers = append(h.headers, pair{key, value})
	return h
}

func (h *HttpHelper) Cookies(key, value string) *HttpHelper {
	h.cookies = append(h.cookies, pair{key, value})
	return h
}

// Param устанавливает значение URL-параметра
func (h *HttpHelper) Param(key, value string) *HttpHelper {
	h.params = append(h.params, pair{key, value})
	return h
}

// Form устанавливает данные, которые будут закодированы
// как application/x-www-form-urlencoded и отправлены в теле запроса
// с соответствующим content-type
func (h *HttpHelper) Form(form map[string]string) *HttpHelper {
	for k, v := range form {
		h.form[k] = v
	}
	h.headers = append(h.headers, pair{"Content-Type", "application/x-www-form-urlencoded"})

	return h
}

// JSON устанавливает данные, которые будут закодированы
// как application/json и отправлены в теле запроса
// с соответствующим content-type
func (h *HttpHelper) JSON(v any) *HttpHelper {
	h.rbody = v
	h.headers = append(h.headers, pair{"Content-Type", "application/json"})
	h.headers = append(h.headers, pair{"Accept", "application/json"})
	return h
}

// Get выполняет GET-запрос с настроенными ранее параметрами
func (h *HttpHelper) Get() *HttpHelperResponse {
	var err error
	req, err := http.NewRequest(http.MethodGet, h.uri, nil)
	if err != nil {
		return &HttpHelperResponse{0, "", err, nil}
	}
	prm := url.Values{}
	for _, p := range h.params {
		prm.Add(p.k, p.v)
	}
	req.URL.RawQuery = prm.Encode()
	for _, h := range h.headers {
		req.Header.Add(h.k, h.v)
	}
	resp, err := h.client.Do(req)

	if err != nil {
		if resp != nil {
			defer resp.Body.Close()
		}
		return &HttpHelperResponse{0, "", err, nil}
	} else {
		return NewHttpHelperResponse(resp)
	}
}

// Post выполняет POST-запрос с настроенными ранее параметрами
func (h *HttpHelper) Post() *HttpHelperResponse {
	var rdr io.Reader
	if len(h.form) > 0 {
		f := url.Values{}
		for k, v := range h.form {
			f.Add(k, v)
		}
		rdr = bytes.NewReader([]byte(f.Encode()))
	}

	if h.rbody != nil {
		if jsondata, err := json.Marshal(h.rbody); err != nil {
			return &HttpHelperResponse{0, "", err, nil}
		} else {
			rdr = bytes.NewReader(jsondata)
		}
	}

	req, err := http.NewRequest(http.MethodPost, h.uri, rdr)
	if err != nil {
		return &HttpHelperResponse{0, "", err, nil}
	}

	prm := url.Values{}
	for _, p := range h.params {
		prm.Add(p.k, p.v)
	}
	req.URL.RawQuery = prm.Encode()

	for _, h := range h.headers {
		req.Header.Add(h.k, h.v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		if resp != nil {
			defer resp.Body.Close()
		}
		return &HttpHelperResponse{0, "", err, nil}
	} else {
		return NewHttpHelperResponse(resp)
	}

}

// HttpHelperResponse представляет ответ на HTTP-запрос
type HttpHelperResponse struct {
	StatusCode int
	Status     string
	err        error
	response   []byte
}

func NewHttpHelperResponse(rsp *http.Response) *HttpHelperResponse {
	b, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return &HttpHelperResponse{0, "", err, nil}
	}
	return &HttpHelperResponse{rsp.StatusCode, rsp.Status, nil, b}
}

// OK возвращает true, если во время выполнения запроса
// не произошло ошибок, а код HTTP-статуса ответа равен 200
func (r *HttpHelperResponse) OK() bool {
	return r.StatusCode == 200
}

// Bytes возвращает тело ответа как срез байт
func (r *HttpHelperResponse) Bytes() []byte {
	return r.response
}

// String возвращает тело ответа как строку
func (r *HttpHelperResponse) String() string {
	return string(r.response)
}

// JSON декодирует тело ответа из JSON и сохраняет
// результат по адресу, на который указывает v
func (r *HttpHelperResponse) JSON(v any) {
	// работает аналогично json.Unmarshal()
	// если при декодировании произошла ошибка,
	// она должна быть доступна через r.Err()
	rb := r.response
	if err := json.Unmarshal(rb, v); err != nil {
		r.err = err
	}
}

// Err возвращает ошибку, которая возникла при выполнении запроса
// или обработке ответа
func (r *HttpHelperResponse) Err() error {
	return r.err
}
