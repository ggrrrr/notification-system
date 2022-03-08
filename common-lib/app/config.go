package app

import (
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type appConfig struct {
	name         string
	shuttingDown bool
	httpServer   *http.Server
	httpRouter   *mux.Router
	osSignal     chan os.Signal
	isReady      *atomic.Value
}

var (
	cfg appConfig
)

func init() {
	cfg = appConfig{
		name:         "_not_set_",
		shuttingDown: false,
	}
}

// Set app name and HTTP listen address
func Configure(name string, osSignal chan os.Signal) {
	addr := config.GetString(name, "listen.addr")
	if addr == "" {
		addr = ":8080"
	}
	cfg = appConfig{
		name:     name,
		osSignal: osSignal,
		isReady:  &atomic.Value{},
	}
	cfg.isReady.Store(true)

	configure(addr)
}

func configure(addr string) {
	cfg.httpRouter = mux.NewRouter()
	cfg.httpRouter.NotFoundHandler = handle404()
	cfg.httpRouter.MethodNotAllowedHandler = method404()
	cfg.httpRouter.HandleFunc("/readyz", readyz())
	cfg.httpRouter.HandleFunc("/healthz", readyz())

	cfg.httpServer = &http.Server{
		Addr:    addr,
		Handler: cfg.httpRouter,
	}
}

// Blocking code
func Start() {
	logrus.Infof("listenAndServe: %v", cfg.httpServer.Addr)
	err := cfg.httpServer.ListenAndServe()
	logrus.Infof("http stoped %v", err)
	cfg.shuttingDown = true
}

// Application name
func GetName() string {
	return cfg.name
}

// will be false only when when shutting down
func ShuttingDown() bool {
	return cfg.shuttingDown
}

// Call this when fatal error
func Panic(err error) {
	logrus.Errorf("app panic from err: %v", err)
	logrus.Errorf("closing osSignal to exit")
	cfg.shuttingDown = true
	close(cfg.osSignal)
}
