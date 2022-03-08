package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTP_HEADER_ string
type HTTP_CONTENT_TYPE_ string

const (
	HTTP_CONTENT_TYPE HTTP_HEADER_       = "Content-Type"
	HTTP_SERVER       HTTP_HEADER_       = "Server"
	HTTP_CT_JSON      HTTP_CONTENT_TYPE_ = "application/json"
)

type ResponseData struct {
	Status    int       `json:"status"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func HandleFunc(path string,
	f func(http.ResponseWriter,
		*http.Request)) *mux.Route {
	funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "/")
	logrus.Infof("endpoint Handle: uri:%s -> %10s", path, funcName[len(funcName)-1])
	return cfg.httpRouter.HandleFunc(path, f)
}

func HttpError(w http.ResponseWriter, code int, msg string, err error) {
	setResponseHeader(w)
	w.WriteHeader(code)
	body := &ResponseData{
		Status:    code,
		Message:   msg,
		Error:     fmt.Sprintf("%+v", err),
		Timestamp: time.Now(),
	}
	json.
		NewEncoder(w).
		Encode(body)

}

func HttpOk(w http.ResponseWriter, msg string) {
	setResponseHeader(w)
	w.WriteHeader(http.StatusAccepted)
	body := &ResponseData{
		Status:    http.StatusAccepted,
		Message:   msg,
		Error:     "",
		Timestamp: time.Now(),
	}
	json.
		NewEncoder(w).
		Encode(body)

}
