package regexphandler

import (
	"log"
	"net/http"
	"regexp"
	"web-app-go/internal/cors"
)

// Regular expression handler
type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

// Handler Regexp functionality
func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// ensure all requests will have CORS headers when applicable
	headers, _ := cors.GetCORS(r.Header.Clone())
	for _, header := range headers {
		w.Header().Set(header.Name, header.Value)
	}
	// Return immediately on OPTIONS with CORS headers
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("")); err != nil {
			log.Printf("[WARNING ] Error to write response: ", err.Error())
		}
		return
	}

	// process routes further
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}

	// if no match by any rule then return 404
	http.NotFound(w, r)
}
