package localstorage

import storageservice "github.com/darkard2003/wormhole/services/storage_service"

func ToStorageError(err error, filepath string) (storageError *storageservice.StorageError) {
	return storageservice.NewStorageError("Something Went Wrong", filepath)
}
