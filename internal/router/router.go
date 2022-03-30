package router

import (
	"log"
	"net/http"
	"regexp"
	"web-app-go/internal/router/quotes"
	"web-app-go/internal/router/regexphandler"
)

func NewRouter() http.Handler {
	mux := &regexphandler.RegexpHandler{}
	mux.HandleFunc(regexp.MustCompile(`^/health$`), healthCheckHandler)
	mux.HandleFunc(regexp.MustCompile(`^/$`), quotes.Handler)
	mux.HandleFunc(regexp.MustCompile(`^/`), notImplemented)
	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		notImplemented(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		log.Printf("[WARNING ] Error to write response: ", err.Error())
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte(http.StatusText(http.StatusNotImplemented))); err != nil {
		log.Printf("[WARNING ] Error to write response: ", err.Error())
	}
}
