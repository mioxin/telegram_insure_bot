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

type pair struct {
	k, v string
}

// HttpHelper предоставляет удобный интерфейс
// для выполнения HTTP-запросов
type HttpHelper struct {
	client       *http.Client
	req          *http.Request
	uri          string
	cookies      []pair
	headers      []pair
	params, form []pair
}

// NewHttpHelper создает новый экземпляр HttpHelper
func NewHttpHelper() *HttpHelper {
	//req.ProxyHeader.Set("User-Agent", "Wget/1.9.1")
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 10
	transport.GetProxyConnectHeader = func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error) {
		return http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"}}, nil
	}
	cl := http.Client{
		Timeout:   time.Duration(3 * time.Second),
		Transport: transport,
	}
	uri := ""
	cook, head, param := []pair{}, []pair{}, []pair{}
	return &HttpHelper{&cl, nil, uri, cook, head, param, []pair{}}
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

// Header устанавливает значение заголовка
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
		h.form = append(h.form, pair{k, v})
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
	for _, par := range h.params {
		p.Add(par.k, par.v)
	}
	h.req.URL.RawQuery = p.Encode()
	for _, hd := range h.headers {
		h.req.Header.Add(hd.k, hd.v)
	}
	resp, err := h.client.Do(h.req)

	if err != nil {
		return &HttpHelperResponse{0, "", err, resp}
	} else {
		return &HttpHelperResponse{resp.StatusCode, resp.Status, nil, resp}
	}
}

// Post выполняет POST-запрос с настроенными ранее параметрами
func (h *HttpHelper) Post() *HttpHelperResponse {
	if h.form != nil {
		f := url.Values{}
		for _, par := range h.form {
			f.Add(par.k, par.v)
		}
		if resp, err := h.client.PostForm(h.uri, f); err != nil {
			return &HttpHelperResponse{0, "", err, resp}
		} else {
			return &HttpHelperResponse{resp.StatusCode, resp.Status, nil, resp}
		}
	} else {
		resp, err := h.client.Do(h.req)
		if err != nil {
			return &HttpHelperResponse{0, "", err, resp}
		} else {
			return &HttpHelperResponse{resp.StatusCode, resp.Status, nil, resp}
		}

	}
}

// HttpHelperResponse представляет ответ на HTTP-запрос
type HttpHelperResponse struct {
	StatusCode int
	Status     string
	err        error
	response   *http.Response
}

// OK возвращает true, если во время выполнения запроса
// не произошло ошибок, а код HTTP-статуса ответа равен 200
func (r *HttpHelperResponse) OK() bool {
	if r.StatusCode >= 200 && r.StatusCode <= 302 {
		return true
	}
	r.err = fmt.Errorf("HttpHelperResponse.OK: error request %v, %v", r.response.Request.URL.String(), r.Status)
	return false
}

// Bytes возвращает тело ответа как срез байт
func (r *HttpHelperResponse) Bytes() []byte {
	if b, err := io.ReadAll(r.response.Body); err == nil {
		defer r.response.Body.Close()
		return b
	} else {
		r.err = err
	}
	return []byte{}
}

// String возвращает тело ответа как строку
func (r *HttpHelperResponse) String() string {
	if b, err := io.ReadAll(r.response.Body); err != nil {
		r.err = err
		return ""
	} else {
		defer r.response.Body.Close()
		return string(b)
	}
}

// JSON декодирует тело ответа из JSON и сохраняет
// результат по адресу, на который указывает v
func (r *HttpHelperResponse) JSON(v any) {
	// работает аналогично json.Unmarshal()
	// если при декодировании произошла ошибка,
	// она должна быть доступна через r.Err()

	//log.Println("htmlHelper: valid ", json.Valid(r.Bytes()))
	rb := r.Bytes()
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
