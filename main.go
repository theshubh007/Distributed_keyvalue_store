package main

import (
	"distributed_keyvalue_store/kvstore"
	"fmt"
	"net/http"
)

func main() {
	//numShards := 5: Defines the number of shards (data partitions) for the key-value store.
	// numReplicas := 2: Defines the number of replicas for fault tolerance.
	numShards := 5  
	numReplicas := 2
	store := kvstore.NewKeyValueStore(numShards, numReplicas)

	// This http server handles each incoming request in a separate goroutine by default.
	http.HandleFunc("/set", store.SetHandler)
	http.HandleFunc("/get", store.GetHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
