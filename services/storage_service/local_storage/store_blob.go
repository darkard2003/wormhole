package localstorage

import (
	"os"

	"github.com/darkard2003/wormhole/services/db"
	"github.com/google/uuid"
)

func (s *LocalStorage) StoreBlob(data []byte) (string, error) {
	id := uuid.New()
	storagePath := s.location + "/" + id.String()

	err := os.WriteFile(storagePath, data, os.ModePerm)

	if err != nil {
		return "", ToStorageError(err, storagePath)
	}

	return id.String(), nil
}

func (s *LocalStorage) GetBlob(id string) ([]byte, error) {
	storagePath := s.location + "/" + id
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil, db.NewNotFoundError("blob", storagePath)
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, ToStorageError(err, storagePath)
	}

	return data, nil
}
