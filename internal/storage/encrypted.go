package storage

import (
	"bytes"
	"filippo.io/age"
	"github.com/icali-app/icali-tui/internal/config"
	"io"
)

type AgeEncryptedStorage struct {
	backend   Storage
	identity  *age.ScryptIdentity
	recipient *age.ScryptRecipient
}

func WithAgeEncryption(backend Storage, config config.Config) Storage {
	identity, err := age.NewScryptIdentity(config.Encryption.Password)
	if err != nil {
		return nil
	}

	recipient, err := age.NewScryptRecipient(config.Encryption.Password)
	if err != nil {
		return nil
	}

	return AgeEncryptedStorage{
		backend:   backend,
		identity:  identity,
		recipient: recipient,
	}
}

func (w AgeEncryptedStorage) Upload(content []byte) error {
	encryptedBuffer := &bytes.Buffer{}
	encryptionStream, err := age.Encrypt(encryptedBuffer, w.recipient)
	if err != nil {
		return err
	}

	_, err = encryptionStream.Write(content)
	if err != nil {
		return err
	}

	return w.Upload(encryptedBuffer.Bytes())
}

func (w AgeEncryptedStorage) Download() ([]byte, error) {
	encryptedContent, err := w.backend.Download()
	if err != nil {
		return nil, err
	}

	encryptedReader := bytes.NewReader(encryptedContent)

	decryptedReader, err := age.Decrypt(encryptedReader, w.identity)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(decryptedReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
