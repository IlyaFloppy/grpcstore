package server

//go:generate mockgen -destination=mocks/interfaces.go . IStorage
type IStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}
