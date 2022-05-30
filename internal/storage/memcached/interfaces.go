package memcached

type IMemcachedClient interface {
	Close() error
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
