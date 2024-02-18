package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

func buildLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if os.Getenv("ENV") == "DEVELOPMENT" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelWarnValue = "WARNING"
	zerolog.LevelErrorValue = "ERROR"

	return logger
}

var logger = buildLogger()

func parseHTTPBody(r interface{}) (map[string]interface{}, error) {
	switch v := r.(type) {
	case *http.Response:
		var jsonRequest map[string]interface{}
		requestBody, err := io.ReadAll(v.Body)
		if err != nil {
			logger.Error().Msgf("%s", err)
		} else {
			if err := json.Unmarshal(requestBody, &jsonRequest); err != nil {
				logger.Error().Msgf("%s", err)
			}
		}
		// Restore request body since ioutil.ReadAll consumes it
		v.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		return jsonRequest, nil
	case *http.Request:
		var jsonRequest map[string]interface{}
		requestBody, err := io.ReadAll(v.Body)
		if err != nil {
			logger.Error().Msgf("%s", err)
		} else {
			if err := json.Unmarshal(requestBody, &jsonRequest); err != nil {
				logger.Error().Msgf("%s", err)
			}
		}
		// Restore request body since ioutil.ReadAll consumes it
		v.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		return jsonRequest, nil
	}

	return nil, errors.New("unsupported HTTP object type")
}
