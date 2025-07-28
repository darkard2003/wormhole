package storageservice

import "fmt"

type StorageError struct {
	FilePath string
	Message  string
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("Storage error: %s %s", e.Message, e.FilePath)
}

func NewStorageError(message, filepath string) *StorageError {
	return &StorageError{
		filepath,
		message,
	}
}

func NewInternalStorageError(filepath string) *StorageError {
	return &StorageError{
		filepath,
		"Something Went Wrong",
	}
}

func NewNotFoundStorageError(file string) *StorageError {
	return NewStorageError("File Not Found", file)
}
