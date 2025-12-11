package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

func ControlProxy(r *mux.Router) {
	controlAddr := os.Getenv("CONTROL_URL")
	if controlAddr == "" {
		controlAddr = "http://localhost:8083"
	}

	controlURL, err := url.Parse(controlAddr)
	if err != nil {
		log.Println("Неверный адрес control:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(controlURL)

	r.PathPrefix("/control").Handler(http.StripPrefix("/control", proxy))

	log.Println("Проксируем /control →", controlAddr)
}
