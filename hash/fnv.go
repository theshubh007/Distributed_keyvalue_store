package hash

const prime = 16777619
const offset = 2166136261

func FnvHash(data string) uint32 {
	hash := uint32(offset)
	for i := 0; i < len(data); i++ {
		hash ^= uint32(data[i])
		hash *= prime
	}
	return hash
}
