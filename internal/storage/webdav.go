package storage

import (
	"github.com/icali-app/icali-tui/internal/config"
	"github.com/studio-b12/gowebdav"
)

type WebDAVStorage struct {
	client     *gowebdav.Client
	remoteFile string
}

func CreateWebDav(config config.Config) Storage {
	client := gowebdav.NewClient(
		config.WebDAV.URL,
		config.WebDAV.Username,
		config.WebDAV.Password,
	)

	return WebDAVStorage{
		client:     client,
		remoteFile: config.WebDAV.RemotePath,
	}
}

func (w WebDAVStorage) Upload(content []byte) error {
	// Upload the file to the WebDAV server
	return w.client.Write(w.remoteFile, content, 0644)
}

func (w WebDAVStorage) Download() ([]byte, error) {
	// Retrieve the file from the WebDAV server
	return w.client.Read(w.remoteFile)
}
