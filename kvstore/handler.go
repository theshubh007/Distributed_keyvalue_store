package kvstore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (kv *KeyValueStore) SetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttlStr := r.URL.Query().Get("ttl")

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		http.Error(w, "Invalid TTL value", http.StatusBadRequest)
		return
	}

	kv.Set(key, value, ttl)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (kv *KeyValueStore) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, ok := kv.Get(key)
	if !ok {
		http.Error(w, "Key not found or expired", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}

func (kv *KeyValueStore) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	status := kv.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
