package cache

type Cache interface {
	Set(key string, value []byte, ttl int64) error
	Get(key string) ([]byte, error)
}
