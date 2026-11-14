// Command sample-orders-api is the first tenant app of the ithaca IDP (idp-infra docs/order-04-first-app.md). It
// proves the whole path: code, image, Spinnaker pipeline, portal secret in the environment, route over HTTPS.
// GET / answers 200 with JSON naming each expected environment variable as "set" or "unset". It never prints a value.
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// secretKeys are the keys of the portal secret orders-config (Secrets page) that this app expects.
var secretKeys = []string{"ORDERS_API_KEY", "ORDERS_GREETING"}

func status() map[string]string {
	out := make(map[string]string, len(secretKeys))
	for _, k := range secretKeys {
		if os.Getenv(k) != "" {
			out[k] = "set"
		} else {
			out[k] = "unset"
		}
	}
	return out
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"service": "sample-orders-api", "secret": status()})
	})
	return mux
}

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	srv := &http.Server{Addr: addr, Handler: newMux(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
