package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	var lb LoadBalancer

	services := parseEnvVariables()

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "",
	})
	loggingTransport := &LoggingRoundTripper{
		loadBalancer: &lb,
		services:     services,
	}

	proxy.Transport = loggingTransport

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	})

	// Load balancing handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	logger.Info().Msgf("Load balancer listening on port %s...", port)
	logger.Err(http.ListenAndServe(":"+port, nil)).Msg("")
}
