package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HttpHelper предоставляет удобный интерфейс
// для выполнения HTTP-запросов
type HttpHelper struct {
	client       *http.Client
	req          *http.Request
	uri          string
	cookies      map[string]string
	headers      map[string]string
	params, form map[string]string
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
	return &HttpHelper{&cl, nil, uri, make(map[string]string), make(map[string]string), make(map[string]string), make(map[string]string)}
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
	h.headers[key] = value
	return h
}

// Header устанавливает значение заголовка
func (h *HttpHelper) Cookies(key, value string) *HttpHelper {
	h.cookies[key] = value
	return h
}

// Param устанавливает значение URL-параметра
func (h *HttpHelper) Param(key, value string) *HttpHelper {
	h.params[key] = value
	return h
}

// Form устанавливает данные, которые будут закодированы
// как application/x-www-form-urlencoded и отправлены в теле запроса
// с соответствующим content-type
func (h *HttpHelper) Form(form map[string]string) *HttpHelper {
	for k, v := range form {
		h.form[k] = v
	}
	return h
}

// JSON устанавливает данные, которые будут закодированы
// как application/json и отправлены в теле запроса
// с соответствующим content-type
func (h *HttpHelper) JSON(v any) *HttpHelper {
	h.form = nil
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	h.req, err = http.NewRequest(http.MethodPost, h.uri, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	h.req.Header.Add("Content-Type", "application/json")
	h.req.Header.Add("Accept", "application/json")
	return h
}

// Get выполняет GET-запрос с настроенными ранее параметрами
func (h *HttpHelper) Get() *HttpHelperResponse {
	var err error
	h.req, err = http.NewRequest(http.MethodGet, h.uri, nil)
	if err != nil {
		return &HttpHelperResponse{0, "", err, nil}
	}
	p := url.Values{}
	for k, v := range h.params {
		p.Add(k, v)
	}
	h.req.URL.RawQuery = p.Encode()
	for k, v := range h.headers {
		h.req.Header.Add(k, v)
	}
	resp, err := h.client.Do(h.req)

	if err != nil {
		return &HttpHelperResponse{0, "", err, nil}
	} else {
		return NewHttpHelperResponse(resp)
	}
}

// Post выполняет POST-запрос с настроенными ранее параметрами
func (h *HttpHelper) Post() *HttpHelperResponse {
	if h.form != nil {
		f := url.Values{}
		for k, v := range h.form {
			f.Add(k, v)
		}
		if resp, err := h.client.PostForm(h.uri, f); err != nil {
			return &HttpHelperResponse{0, "", err, nil}
		} else {
			return NewHttpHelperResponse(resp)
		}
	} else {
		resp, err := h.client.Do(h.req)
		if err != nil {
			return &HttpHelperResponse{0, "", err, nil}
		} else {
			return NewHttpHelperResponse(resp)
		}

	}
}

// HttpHelperResponse представляет ответ на HTTP-запрос
type HttpHelperResponse struct {
	StatusCode int
	Status     string
	err        error
	response   []byte
	//response   *http.Response
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
	if r.StatusCode >= 200 && r.StatusCode <= 302 {
		return true
	}
	r.err = fmt.Errorf("HttpHelperResponse.OK: error request %v", r.Status)
	return false
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

	//log.Println("htmlHelper: valid ", json.Valid(r.Bytes()))
	r.err = nil
	rb := r.response
	if json.Valid(rb) {
		if err := json.Unmarshal(rb, v); err != nil {
			r.err = err
		}
	} else {
		r.err = fmt.Errorf("HttpHelperResponse.JSON: error %v is not valid json", r.String())
	}
}

// Err возвращает ошибку, которая возникла при выполнении запроса
// или обработке ответа
func (r *HttpHelperResponse) Err() error {
	if r.err != nil {
		return r.err
	}
	return nil
}
