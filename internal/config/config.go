package config

import "time"

type Config struct {
	LoggerConfig  LoggerConfig  `yaml:"logger"`
	ServerConfig  ServerConfig  `yaml:"server"`
	StorageConfig StorageConfig `yaml:"storage"`
}

type LoggerConfig struct {
	Level   string `yaml:"level"`
	Console bool   `yaml:"console"`
}

type ServerConfig struct {
	Address          string        `yaml:"address"`
	KeepAliveTime    time.Duration `yaml:"keep_alive_time"`
	KeepAliveTimeout time.Duration `yaml:"keep_alive_timeout"`
	WriteBufferSize  int           `yaml:"write_buffer_size"`
	ReadBufferSize   int           `yaml:"read_buffer_size"`
}

type StorageConfig struct {
	UseMemcached           bool                   `yaml:"use_memcached"`
	InMemoryStorageConfig  InMemoryStorageConfig  `yaml:"inmemory"`
	MemcachedStorageConfig MemcachedStorageConfig `yaml:"memcached"`
}

type InMemoryStorageConfig struct{}

type MemcachedStorageConfig struct {
	Address  string `yaml:"address"`
	UsePool  bool   `yaml:"use_pool"`
	PoolSize int    `yaml:"pool_size"`
}
