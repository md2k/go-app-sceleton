package cors

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"web-app-go/internal/config"
	"web-app-go/internal/utils"
)

type Header struct {
	Name  string
	Value string
}

func GetCORS(h http.Header) (headers []Header, err error) {
	cfg := config.GetConfig()
	origin := h.Get("Origin")
	if len(origin) <= 0 {
		err = errors.New("CORS is not set because header 'Origin' not set in request...")
		return
	}

	originUrl, err := url.Parse(origin)
	if err != nil {
		err = fmt.Errorf("CORS is not set because value of Origin's header is not correct. %s", err.Error())
		return
	}

	ohost := strings.Split(originUrl.Host, ":")
	if !utils.IsStirngIn(ohost[0], cfg.Cors.Origins) {
		err = fmt.Errorf("CORS is not set because request Origin is not in Allowed Origin List")
		return
	}

	// Set Access-Control-Allow-Origin
	headers = append(headers, Header{Name: "Access-Control-Allow-Origin", Value: fmt.Sprintf("%s://%s", originUrl.Scheme, originUrl.Host)})
	// Set Access-Control-Max-Age
	headers = append(headers, Header{Name: "Access-Control-Max-Age", Value: strconv.FormatFloat(cfg.Cors.MaxAge.Seconds(), 'f', 0, 64)})
	// Set Access-Control-Allow-Credentials
	headers = append(headers, Header{Name: "Access-Control-Allow-Credentials", Value: strconv.FormatBool(cfg.Cors.Credentials)})
	// Set Access-Control-Allow-Methods
	if len(cfg.Cors.Allow.Methods) > 0 {
		headers = append(headers, Header{Name: "Access-Control-Allow-Methods", Value: strings.Join(cfg.Cors.Allow.Methods, ",")})
	}
	// Set Access-Control-Allow-Headers
	allowedHeaders := strings.TrimRight(strings.Join(append(cfg.Cors.Allow.Headers, h.Get("Access-Control-Request-Headers")), ","), ",")
	if len(allowedHeaders) > 0 {
		headers = append(headers, Header{Name: "Access-Control-Allow-Headers", Value: allowedHeaders})
	}
	// Set Access-Control-Expose-Headers
	if len(cfg.Cors.Expose.Headers) > 0 {
		headers = append(headers, Header{Name: "Access-Control-Expose-Headers", Value: strings.Join(cfg.Cors.Expose.Headers, ",")})
	}

	return
}
