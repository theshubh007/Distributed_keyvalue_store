package main

import (
	"distributed_keyvalue_store/kvstore"
	"distributed_keyvalue_store/logger"
	"net/http"
)

func main() {
	//numShards := 5: Defines the number of shards (data partitions) for the key-value store.
	// numReplicas := 2: Defines the number of replicas for fault tolerance.
	numShards := 5
	numReplicas := 2
	store := kvstore.NewKeyValueStore(numShards, numReplicas)

	// This http server handles each incoming request in a separate goroutine by default.
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/set", store.SetHandler)
	http.HandleFunc("/get", store.GetHandler)
	http.HandleFunc("/dashboard", store.DashboardHandler)

	logger.Info.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
