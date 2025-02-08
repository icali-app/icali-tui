package webdav

import (
	"github.com/icali-app/icali-tui/internal/config"
	"github.com/icali-app/icali-tui/internal/storage"
	"os"

	"github.com/studio-b12/gowebdav"
)

type WebDAVStorage struct {
	client     *gowebdav.Client
	localFile  string
	remoteFile string
}

func Create(config config.Config) storage.Storage {
	client := gowebdav.NewClient(
		config.WebDAV.URL,
		config.WebDAV.Username,
		config.WebDAV.Password,
	)

	return WebDAVStorage{
		client:     client,
		localFile:  "TODO.ics",
		remoteFile: config.WebDAV.RemotePath,
	}
}

func (w WebDAVStorage) Upload() error {
	// Open the local file
	data, err := os.ReadFile(w.localFile)
	if err != nil {
		return err
	}

	// Upload the file to the WebDAV server
	err = w.client.Write(w.remoteFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (w WebDAVStorage) Download() error {
	// Retrieve the file from the WebDAV server
	data, err := w.client.Read(w.remoteFile)
	if err != nil {
		return err
	}

	// Write the file locally
	err = os.WriteFile(w.localFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
