package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoggingRoundTripper struct {
	loadBalancer *LoadBalancer
	services     map[string][]*Backend
}

func (lrt *LoggingRoundTripper) ParseServiceHeader(r *http.Request) (string, error) {
	service := strings.ToUpper(r.Header.Get("x-nexus-service"))

	_, found := lrt.services[service]
	if !found {
		errMsg := fmt.Sprintf("requested service %s doesn't exists", service)
		logger.Error().Msgf(errMsg)
		return "", errors.New(errMsg)
	}

	return service, nil
}

func (lrt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Save Request body as JSON
	reqBody, err := parseHTTPBody(r)
	if err != nil {
		logger.Error().Msgf("%s", err)
		return nil, err
	}

	service, err := lrt.ParseServiceHeader(r)
	if err != nil {
		return nil, err
	}

	maxRetries := 5
	retryInterval := 1 * time.Second

	var rs *http.Response
	bc, err := io.ReadAll(r.Body)
	for retries := 0; retries < maxRetries; retries++ {
		b, err := lrt.loadBalancer.getNextBackend(service, lrt.services[service])
		if err != nil {
			logger.Error().Msgf("%s", err)
			return nil, err
		}
		r.Host = b.Host
		r.URL.Host = b.Host
		r.Header.Set("x-nexus-service", service)
		r.Header.Set("x-nexus-location", strings.ToLower(b.Location))
		for k, v := range b.Headers {
			r.Header.Set(k, v)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bc))
		start := time.Now()
		rs, err = http.DefaultTransport.RoundTrip(r)
		if err != nil {
			logger.Error().Msgf("%s", err)
		}
		duration := time.Since(start)
		// Save Response body as JSON
		resBody, err := parseHTTPBody(rs)
		if err != nil {
			logger.Error().Msgf("%s", err)
		}

		logger.Info().Interface("request", reqBody).Interface("response", resBody).
			Str("service", r.Header.Get("x-nexus-service")).
			Str("location", r.Header.Get("x-nexus-location")).
			Msgf("%s - %s - %d - %s", r.Method, r.URL.Path, rs.StatusCode, duration.String())
		if err != nil || rs.StatusCode >= 300 {
			time.Sleep(retryInterval)
		} else {
			break
		}
	}
	if err != nil {
		errMSg := fmt.Sprintf("Unable to fulfill request after %d retries", maxRetries)
		logger.Error().Msg(errMSg)
		return nil, errors.New(errMSg)
	}

	return rs, nil
}
