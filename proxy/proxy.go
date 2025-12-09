package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

func AdProxy(r *mux.Router) {
	adAddr := os.Getenv("AD_URL")
	if adAddr == "" {
		adAddr = "http://localhost:8083"
	}

	adURL, err := url.Parse(adAddr)
	if err != nil {
		log.Println("Неверный адрес ad:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(adURL)

	r.PathPrefix("/ad").Handler(http.StripPrefix("/ad", proxy))

	log.Println("Проксируем /ad →", adAddr)
}
