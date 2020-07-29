package email

// Repository provides access to the storage.
type Repository interface {
	Add(email, token string) error
	Read() (map[string]string, error)
	Remove(key string) error
	Seek(email string) error
}

// Service provides email lists operations.
type Service interface {
	Add(email, token string) error
	Read() (map[string]string, error)
	Remove(key string) error
	Seek(email string) error
}
