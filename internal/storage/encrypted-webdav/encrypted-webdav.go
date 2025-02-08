package encrypted

import (
	"bytes"
	"github.com/icali-app/icali-tui/internal/config"
	"github.com/icali-app/icali-tui/internal/storage"
	"io"
	"os"

	"filippo.io/age"
	"github.com/studio-b12/gowebdav"
)

type EncryptedWebDAVStorage struct {
	client     *gowebdav.Client
	localFile  string
	remoteFile string
	identity   *age.ScryptIdentity
	recipient  *age.ScryptRecipient
}

func Create(config config.Config) storage.Storage {
	client := gowebdav.NewClient(
		config.WebDAV.URL,
		config.WebDAV.Username,
		config.WebDAV.Password,
	)

	identity, err := age.NewScryptIdentity(config.Encryption.Password)
	if err != nil {
		return nil
	}

	recipient, err := age.NewScryptRecipient(config.Encryption.Password)
	if err != nil {
		return nil
	}

	return EncryptedWebDAVStorage{
		client:     client,
		localFile:  "TODO.ics",
		remoteFile: config.WebDAV.RemotePath,
		identity:   identity,
		recipient:  recipient,
	}
}

func (w EncryptedWebDAVStorage) Upload() error {
	// TODO Maybe find more ideal way to implement this using reader/writers?

	// Open the local file
	unencryptedBuffer, err := os.ReadFile(w.localFile)
	if err != nil {
		return err
	}

	encryptedBuffer := &bytes.Buffer{}
	encryptionStream, err := age.Encrypt(encryptedBuffer, w.recipient)
	if err != nil {
		return err
	}

	_, err = encryptionStream.Write(unencryptedBuffer)
	if err != nil {
		return err
	}

	err = w.client.Write(w.remoteFile, encryptedBuffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (w EncryptedWebDAVStorage) Download() error {
	// TODO Maybe find more ideal way to implement this using reader/writers?

	// Retrieve the file from the WebDAV server
	encryptedReader, err := w.client.ReadStream(w.remoteFile)
	if err != nil {
		return err
	}

	decryptedReader, err := age.Decrypt(encryptedReader, w.identity)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(decryptedReader)
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
