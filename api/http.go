package api

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	js "github.com/midnightrun/hexagonal-architecture-url-shortener-example/serializer/json"
	msg "github.com/midnightrun/hexagonal-architecture-url-shortener-example/serializer/msgpack"
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

var log shortener.Logger

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
	log             shortener.Logger
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msg.Redirect{}
	}

	return &js.Redirect{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	log.Infof("Get handler triggered")

	code := chi.URLParam(r, "code")

	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if err == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			log.Warnf("Get handler: %s", err)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Warnf("Get handler: %s", err)

		return
	}

	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Warnf("Post handler: %s", err)

		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Warnf("Post handler: %s", err)

		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if err == shortener.ErrReadirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Warnf("Post handler: %s", err)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Warnf("Post handler: %s", err)

		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Warnf("Post handler: %s", err)

		return
	}

	setupResponse(w, contentType, responseBody, http.StatusCreated)
}

func NewHandler(redirectService shortener.RedirectService, log shortener.Logger) RedirectHandler {
	return &handler{
		redirectService: redirectService,
		log:             log}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)

	if err != nil {
		log.Warnf("setupResponse: %s", err)
	}
}
