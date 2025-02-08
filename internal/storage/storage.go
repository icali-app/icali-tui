package storage

type Storage interface {
	Upload() error
	Download() error
}
