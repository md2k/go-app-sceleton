package quotes

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
	"web-app-go/internal/config"
)

var quotesTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	MaxIdleConns:          100,
	MaxIdleConnsPerHost:   100,
	TLSHandshakeTimeout:   5 * time.Second,
	DisableKeepAlives:     true,
	ExpectContinueTimeout: 1 * time.Second,
}

var httpClient *http.Client = &http.Client{
	Transport: quotesTransport,
	Timeout:   time.Second * 10,
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// because i'm lazzy add custom handler to config loader, i'm going to load config and pars Quotes.URL here
	// anyway config object are atomic storage with safe-thread read feature.
	cfg := config.GetConfig()
	quoteUrl, _ := url.Parse(cfg.Quotes.Url) // skip errors, it is no prod app
	quoteUrl.Path = cfg.Quotes.Path

	request, _ := http.NewRequest(r.Method, quoteUrl.String(), nil)
	request.Header.Add("Accept", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	quoteStruct := struct {
		Author    string `json:"author,omitempty"`
		Id        int    `json:"id,omitempty"`
		Quote     string `json:"quote,omitempty"`
		Permalink string `json:"permalink,omitempty"`
	}{}

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		w.Write([]byte(http.StatusText(resp.StatusCode)))
		return
	}

	//var body []byte
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&quoteStruct); err != nil {
		log.Printf("[ERROR   ] Unable to parse json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	tmpl := template.Must(template.ParseFiles("html/index.tmpl"))
	tmpl.Execute(w, quoteStruct)
}
