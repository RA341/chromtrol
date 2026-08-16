package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var servicesToToggle = []string{
	"svc-de",
	"svc-nginx",
	"svc-selkies",
	"svc-pulseaudio",
	"svc-xsettingsd",
	"svc-xorg",
	"svc-dbus",
	"svc-cron",
	"svc-docker",
	"svc-watchdog",
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Start in reverse order (infrastructure first)
	var errors []string
	for i := len(servicesToToggle) - 1; i >= 0; i-- {
		svc := servicesToToggle[i]
		cmd := exec.Command("s6-svc", "-u", "/run/service/"+svc)
		if output, err := cmd.CombinedOutput(); err != nil {
			errors = append(errors, fmt.Sprintf("failed to start %s: %v (%s)", svc, err, strings.TrimSpace(string(output))))
		}
		time.Sleep(50 * time.Millisecond) // brief pause for initialization order
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "Started display and browser environment",
	}
	if len(errors) > 0 {
		resp["status"] = "partial_error"
		resp["errors"] = errors
		w.WriteHeader(http.StatusMultiStatus)
	}

	json.NewEncoder(w).Encode(resp)
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Stop in forward order (applications and proxies first)
	var errors []string
	for _, svc := range servicesToToggle {
		cmd := exec.Command("s6-svc", "-d", "/run/service/"+svc)
		if output, err := cmd.CombinedOutput(); err != nil {
			errors = append(errors, fmt.Sprintf("failed to stop %s: %v (%s)", svc, err, strings.TrimSpace(string(output))))
		}
	}

	// Kill any orphaned sleep processes left behind by stopped service scripts
	exec.Command("pkill", "-f", "sleep").Run()

	resp := map[string]interface{}{
		"status":  "success",
		"message": "Stopped display and browser environment",
	}
	if len(errors) > 0 {
		resp["status"] = "partial_error"
		resp["errors"] = errors
		w.WriteHeader(http.StatusMultiStatus)
	}

	json.NewEncoder(w).Encode(resp)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cmd := exec.Command("s6-svstat", "/run/service/svc-de")
	output, err := cmd.CombinedOutput()

	statusText := strings.TrimSpace(string(output))
	isRunning := strings.Contains(statusText, "up")

	resp := map[string]interface{}{
		"status":     "success",
		"running":    isRunning,
		"raw_status": statusText,
	}
	if err != nil {
		resp["status"] = "error"
		resp["message"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(resp)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)
	http.HandleFunc("/status", handleStatus)

	log.Printf("Control plane server running on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
