package storage

import (
	"bytes"
	"crypto/sha256"
	"errors"
	ics "github.com/arran4/golang-ical"
	"os"
)

type CheckedStorage struct {
	backend   Storage
	hash      [32]byte
	localFile string
}

func WithCollisionChecks(backend Storage) CheckedStorage {
	return CheckedStorage{
		backend:   backend,
		localFile: "idk.ics", // TODO
	}
}

func (c CheckedStorage) Upload(content []byte) error {
	allZero := [32]byte{}
	if allZero == c.hash {
		return c.Upload(content)
	}

	download, err := c.backend.Download()
	if err != nil {
		return err
	}

	remoteHash := sha256.Sum256(download)
	if remoteHash != c.hash {
		return errors.New("hash mismatch") // TODO Improve error messages
	}

	return c.Upload(content)
}

func (c CheckedStorage) DownloadAndSave() (*ics.Calendar, error) {
	download, err := c.backend.Download()
	if err != nil {
		return nil, err
	}

	c.hash = sha256.Sum256(download)

	reader := bytes.NewReader(download)

	// Parse the calendar
	cal, err := ics.ParseCalendar(reader)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(c.localFile, download, 0644)
	if err != nil {
		return nil, err
	}

	return cal, nil
}
