package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var version = "dev"
var configFile = "config.json"
var config = map[string]string{}
var configMu sync.RWMutex

type ConfigItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func loadConfig() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &config)
}

func saveConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(configFile)
	tmp, err := os.CreateTemp(dir, "config-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), configFile)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"version": version})
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENVIRONMENT")
	json.NewEncoder(w).Encode(map[string]string{"environment": env})
}

func createConfigHandler(w http.ResponseWriter, r *http.Request) {
	var item ConfigItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if item.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	configMu.Lock()
	config[item.Name] = item.Value
	err := saveConfig()
	configMu.Unlock()

	if err != nil {
		http.Error(w, "failed to save config", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func configItemHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/config/")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		configMu.RLock()
		value, found := config[name]
		configMu.RUnlock()

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(ConfigItem{Name: name, Value: value})

	case "DELETE":
		configMu.Lock()
		delete(config, name)
		err := saveConfig()
		configMu.Unlock()

		if err != nil {
			http.Error(w, "failed to save config", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"deleted": true})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	loadConfig()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/config", createConfigHandler)
	http.HandleFunc("/config/", configItemHandler)

	log.Println("server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
