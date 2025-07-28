package localstorage

type LocalStorage struct {
	location string
}

func NewLocalStorage(location string) *LocalStorage {
	return &LocalStorage{
		location: location,
	}
}
