package storage

import (
	"github.com/icali-app/icali-tui/internal/config"
	"sync"
)

var (
	conf    = config.Get()
	once    sync.Once
	storage CheckedStorage
)

type Storage interface {
	Upload([]byte) error
	Download() ([]byte, error)
}

func Get() CheckedStorage {
	once.Do(func() {
		webdav := CreateWebDav(conf)
		encrypted := WithAgeEncryption(webdav, conf)
		storage = WithCollisionChecks(encrypted)
	})

	return storage
}
