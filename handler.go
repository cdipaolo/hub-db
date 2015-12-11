package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cdipaolo/hub-db/model"
)

func HandleStatus(r http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(r, req)
		return
	}

	status := "running"
	if !running {
		status = "stopped"
	}

	bytes, err := json.Marshal(model.HealthCheckResponse{
		Status:       status,
		SaveTo:       Config.DumpPath,
		Uptime:       time.Now().Sub(startTime).String(),
		PagesCrawled: pageCount,
		Errors:       errors,
	})
	if err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		r.Write([]byte(fmt.Sprintf(`{
			"message": "couldn't marshal health check response to JSON",
			"error": %v
		}`, err)))
		return
	}

	// reset error store after dump
	errors = []error{}

	r.WriteHeader(http.StatusOK)
	r.Write(bytes)

	log.Printf("GET / [PagesCrawled = %v, Errors = %v]\n", pageCount, len(errors))
}
