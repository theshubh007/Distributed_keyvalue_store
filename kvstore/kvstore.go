package kvstore

import (
	"distributed_keyvalue_store/hash"
	"distributed_keyvalue_store/logger"
	"sync"
	"time"
)

type KeyValue struct {
	Value      string
	Expiration time.Time
}

type Shard struct {
	data map[string]KeyValue
	// Mutex (mutual exclusion) is used to prevent race conditions when accessing shared resources. This ensures that only one goroutine can access the resource at a time.
	mu sync.RWMutex
}

// shards []*Shard: This slice holds multiple instances of the Shard
// replicas int: This field specifies the number of replicas for each shard
// Replicas provide fault tolerance by ensuring that if one server hosting a shard fails, another replica of that shard can serve requests.

type KeyValueStore struct {
	shards   []*Shard
	replicas int
}

func NewKeyValueStore(numShards, numReplicas int) *KeyValueStore {
	store := &KeyValueStore{
		shards:   make([]*Shard, numShards),
		replicas: numReplicas,
	}

	for i := 0; i < numShards; i++ {
		store.shards[i] = &Shard{data: make(map[string]KeyValue)}
	}

	return store
}

func (kv *KeyValueStore) GetShardIndex(key string) int {
	hash := hash.FnvHash(key)
	return int(hash) % len(kv.shards)
}

func (kv *KeyValueStore) Set(key, value string, ttl time.Duration) {
	shardIndex := kv.GetShardIndex(key)
	shard := kv.shards[shardIndex]
	shard.mu.Lock()
	defer shard.mu.Unlock()

	expiration := time.Now().Add(ttl)
	shard.data[key] = KeyValue{Value: value, Expiration: expiration}
	logger.Info.Printf("Set key=%s, value=%s, ttl=%s", key, value, ttl)
}

func (kv *KeyValueStore) Get(key string) (string, bool) {
	shardIndex := kv.GetShardIndex(key)
	shard := kv.shards[shardIndex]
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	item, ok := shard.data[key]
	if !ok || time.Now().After(item.Expiration) {
		return "", false
	}
	logger.Info.Printf("Get key=%s, value=%s", key, item.Value)
	return item.Value, true
}

func (kv *KeyValueStore) GetStatus() map[string]interface{} {
	status := make(map[string]interface{})
	status["numShards"] = len(kv.shards)
	status["numReplicas"] = kv.replicas

	shardStatus := make([]map[string]interface{}, len(kv.shards))
	for i, shard := range kv.shards {
		shard.mu.RLock()
		shardStatus[i] = map[string]interface{}{
			"numKeys": len(shard.data),
		}
		shard.mu.RUnlock()
	}

	status["shards"] = shardStatus
	return status
}
