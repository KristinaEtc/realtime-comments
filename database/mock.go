package database

import "time"

// MockDB is a struct which implements DataBase interface.
// It could be use when no DB is set.
type MockDB struct {
	data []byte
}

func initMockDB() (*MockDB, error) {
	log2.Debug("Mock: init data")
	return &MockDB{
		data: []byte("test-data-mock "),
	}, nil
}

// GetData is a method of DataBase interface.
func (m *MockDB) GetData() ([]byte, error) {
	log2.Debug("Mock: get data")
	return m.data, nil
}

// InsertData is a method of DataBase interface.
func (m *MockDB) InsertData(data []byte, currTime time.Time) error {
	log2.Debug("Mock: insert data")
	return nil
}

// Close is a method of DataBase interface. Returns nil.
func (m *MockDB) Close() error {
	log2.Debug("Mock: close DB")
	return nil
}
