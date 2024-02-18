package main

import (
	"os"
	"regexp"
	"strings"
)

type Backend struct {
	Host     string            `json:"host"`
	Location string            `json:"location"`
	Headers  map[string]string `json:"headers"`
}

func parseEnvVariables() map[string][]*Backend {
	re := regexp.MustCompile(`^NEXUS_([A-Z0-9]+)_([A-Z0-9]+)_(URL|HEADERS)$`)

	services := map[string][]*Backend{}

	for _, env := range os.Environ() {
		name, value, _ := strings.Cut(env, "=")
		matches := re.FindStringSubmatch(name)
		if matches == nil {
			continue // Skip if no match
		}

		service := matches[1]
		_, found := services[service]
		if !found {
			services[service] = nil
		}

		location := matches[2]

		backend := findBackend(services[service], location)
		if backend == nil {
			backend = &Backend{}
			backend.Location = location
			backend.Headers = make(map[string]string)
			services[service] = append(services[service], backend)
		}

		property := matches[3]
		if property == "URL" {
			backend.Host = value
		}
		if property == "HEADERS" {
			headersList := strings.Split(value, ";")
			for _, header := range headersList {
				k, v, _ := strings.Cut(header, "=")
				backend.Headers[k] = v
			}
		}
		logger.Info().Msgf("Loaded %s %s %s", service, location, property)
	}

	return services
}

func findBackend(backends []*Backend, location string) *Backend {
	for _, backend := range backends {
		if backend.Location == location {
			return backend
		}
	}
	return nil
}
