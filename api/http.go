package api

import (
	"log"
	"net/http"

	js "github.com/midnightrun/hexagonal-architecture-url-shortener-example/serializer/json"
	msg "github.com/midnightrun/hexagonal-architecture-url-shortener-example/serializer/msgpack"
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msg.Redirect{}
	}

	return &js.Redirect{}
}

func (h *handler) Get(http.ResponseWriter, *http.Request) {
}

func (h *handler) Post(http.ResponseWriter, *http.Request) {
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
