package persistence

type Database interface {
	// Initialize() error
	// Save(resp.Value) error
	Close() error
}
