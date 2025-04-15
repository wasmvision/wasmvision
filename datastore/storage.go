package datastore

// DataStorage is an interface for a data storage system.
// It defines methods for getting, setting, deleting, and checking the existence of data.
// The interface is designed to be implemented by different storage backends,
// such as in-memory storage, file-based storage, or database storage.
type DataStorage interface {
	Get(string, string) (string, bool)
	GetKeys(string) ([]string, bool)
	Set(string, string, string) error
	Delete(string, string)
	DeleteAll(string)
	Exists(string) bool
	Error() error
	Close() error
}
