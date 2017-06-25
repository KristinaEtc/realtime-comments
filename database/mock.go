package database

import "time"

// MockDB is a struct which implements DataBase interface.
// It could be use when no DB is set.
type MockDB struct {
}

func initMockDB() (*MockDB, error) {
	log.Debug("Mock: init data")
	return &MockDB{}, nil
}

// GetData is a method of DataBase interface.
func (m *MockDB) GetData() error {
	log.Debug("Mock: get data")
	return nil
}

// InsertData is a method of DataBase interface.
func (m *MockDB) InsertData(data []byte, currTime time.Time) error {
	log.Debug("Mock: insert data")
	return nil
}

// Close is a method of DataBase interface. Returns nil.
func (m *MockDB) Close() error {
	log.Debug("Mock: close DB")
	return nil
}
