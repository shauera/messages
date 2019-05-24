package application

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

// HealthMonitor - inteface for any component that would like to check and report it's health status
type HealthMonitor interface {
	// GetSystemID - provides the id of the system that wants its health checked
	GetSystemID() string

	// IsHealthy - performs the actuall health check
	IsHealthy(ctx context.Context) bool
}

var registryLock sync.Mutex
var ctx context.Context

var healthRegistry = make(map[string]struct {
	monitor HealthMonitor
	status  bool
})

//InitHealth - initializes health scanning, which is used to register, aggregate and query for service health
func InitHealth(ctxParam context.Context) {
	ctx = ctxParam
	start()
}

// AddHealthMonitor - a thread safe method to add health monitors by any
// system that wants its health to be checked and reported back
func AddHealthMonitor(healthMonitor HealthMonitor) {
	registryLock.Lock()
	defer registryLock.Unlock()

	healthRegistry[healthMonitor.GetSystemID()] =
		struct {
			monitor HealthMonitor
			status  bool
		}{
			healthMonitor,
			false, // System is started unhealthy until proven otherwise
		}
}

func start() {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(config.GetDuration("service.healthScanInterval")):
				log.Debug("Scanning systems health status")
				for systemID, system := range healthRegistry {
					newStatus := system.monitor.IsHealthy(ctx)
					healthRegistry[systemID] = struct {
						monitor HealthMonitor
						status  bool
					}{
						monitor: system.monitor,
						status:  newStatus,
					}
				}
			}
		}
	}()
}

//healthResponse - template for rendering health status in HTTP responses
type healthResponse struct {
	Status string `json:"status"`
}

// GetHealth - handler for retrieving health state for the entire service
func GetHealth(response http.ResponseWriter, request *http.Request) {
	currentHealthStatus := "healthy"
	statusCode := http.StatusOK
	for _, system := range healthRegistry {
		if system.status == false {
			currentHealthStatus = "unhealthy"
			statusCode = http.StatusInternalServerError
			break
		}
	}

	response.WriteHeader(statusCode)
	json.NewEncoder(response).Encode(healthResponse{Status: currentHealthStatus})
}
