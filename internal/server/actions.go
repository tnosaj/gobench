package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/gobench/internal"
	"github.com/tnosaj/gobench/internal/work"
)

func (s *GobenchServer) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status": "ok"}`)
}

func (s *GobenchServer) Prepare(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Prepare request")
	var err error
	s.Settings, err = s.unifySettings(r, "prepare")
	if err != nil {
		returnError(w, err, http.StatusInternalServerError)
		return
	}
	go work.Start(s.Settings, s.Strategy)
	return
}

func (s *GobenchServer) Run(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Run request")
	var err error
	s.Settings, err = s.unifySettings(r, "run")
	if err != nil {
		returnError(w, err, http.StatusInternalServerError)
		return
	}
	go work.Start(s.Settings, s.Strategy)
	return
}

func (s *GobenchServer) Cleanup(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("Cleanup request")
	var err error
	s.Settings, err = s.unifySettings(r, "cleanup")
	if err != nil {
		returnError(w, err, http.StatusInternalServerError)
		return
	}
	go work.Start(s.Settings, s.Strategy)
	return
}

func (s *GobenchServer) unifySettings(r *http.Request, action string) (internal.Settings, error) {
	settings := s.Settings
	httpsettings, err := getSettingsFromPost(r)
	if err != nil {
		return s.Settings, err
	}

	settings.Action = action
	if httpsettings.SqlMigrationFolder != "" {
		settings.SqlMigrationFolder = httpsettings.SqlMigrationFolder
	}
	if httpsettings.Strategy != "" {
		settings.Strategy = httpsettings.Strategy
	}
	if httpsettings.Concurrency > 0 {
		settings.Concurrency = httpsettings.Concurrency
	}
	if httpsettings.Initialdatasize > 0 {
		settings.Initialdatasize = httpsettings.Initialdatasize
	}
	if httpsettings.Duration > 0 {
		settings.Duration = httpsettings.Duration
	}
	if httpsettings.Rate > 0 {
		settings.Rate = httpsettings.Rate
	}
	logrus.Infof("Settings %+v", settings.PrintableSettings())
	return settings, nil
}

func getSettingsFromPost(r *http.Request) (HttpSettings, error) {
	untypedRequestBody, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return HttpSettings{}, fmt.Errorf("Body read error: %q", err)
	}

	var typedRequestBody HttpSettings
	err = json.Unmarshal(untypedRequestBody, &typedRequestBody)
	if err != nil {
		return HttpSettings{}, fmt.Errorf("untypedRequestBody Unmarshal error: %q", err)
	}
	return typedRequestBody, nil
}

func returnError(w http.ResponseWriter, err error, httpCode int) {
	logrus.Error(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	io.WriteString(w, fmt.Sprintf(`{"error": "%s"}`, err))
}
