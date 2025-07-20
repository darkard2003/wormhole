package storageservice

type StorageInterface interface {
	StoreBlob(data []byte) (string, error)
	GetBlob(id string) ([]byte, error)
}
