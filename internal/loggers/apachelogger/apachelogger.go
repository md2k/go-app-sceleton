package apachelogger

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ApacheLogRecord struct {
	http.ResponseWriter

	ip                    string
	username              string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	referer               string
	userAgent             string
	xff                   string // X-Forwarded-For if exists
	xri                   string // X-Real-Ip if exists
	elapsedTime           time.Duration
}

const ApacheFormatPattern = "[%s > %s] %s - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" XFF:(%s) (%.4f)\n"

func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	fmt.Fprintf(out, ApacheFormatPattern, r.xri, r.ip, r.username, timeFormatted, r.method,
		r.uri, r.protocol, r.status, r.responseBytes, r.referer, r.userAgent,
		r.xff, r.elapsedTime.Seconds())
}

func (r *ApacheLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

func (r *ApacheLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type ApacheLoggingHandler struct {
	handler http.Handler
	out     io.Writer
}

func NewApacheLoggingHandler(out io.Writer, handler http.Handler) http.Handler {
	return &ApacheLoggingHandler{
		handler: handler,
		out:     out,
	}
}

func (h *ApacheLoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	username := "-"
	if r.URL.User != nil {
		if name := r.URL.User.Username(); name != "" {
			username = name
		}
	} else {
		name, _, ok := r.BasicAuth()
		if ok {
			username = name
		}
	}

	referer := r.Referer()
	if referer == "" {
		referer = "-"
	}

	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "-"
	}

	// Get X-Forwarded-For chain if exists
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		xff = "-"
	}

	// Trying to Get Real IP is specific headers is exists
	xri := "-"
	if r.Header.Get("True-Client-IP") != "" {
		xri = r.Header.Get("True-Client-IP")
	} else if r.Header.Get("CF-Connecting-Ip") != "" {
		xri = r.Header.Get("CF-Connecting-Ip")
	} else if r.Header.Get("X-Real-Ip") != "" {
		xri = r.Header.Get("X-Real-Ip")
	}

	record := &ApacheLogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
		username:       username,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		referer:        referer,
		userAgent:      userAgent,
		xff:            xff,
		xri:            xri,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.out)
}
