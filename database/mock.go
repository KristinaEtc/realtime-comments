package database

import "time"

// MockDB is a struct which implements DataBase interface.
// It could be use when no DB is set.
type MockDB struct {
	data []byte
}

func initMockDB() (*MockDB, error) {
	log.Debug("Mock: init data")
	return &MockDB{
		data: []byte("test-data-mock "),
	}, nil
}

// GetLastComments is a method of DataBase interface.
func (m *MockDB) GetLastComments(videoID int64, id int) ([]Comment, error) {
	log.Debug("Mock: get data")

	var c = []Comment{
		Comment{
			ID:                1,
			Username:          "kolia",
			CommentBody:       "comment",
			VideoID:           2,
			VideoTimestamp:    4,
			CalendarTimestamp: time.Now(),
		},
	}

	return c, nil
}

// InsertData is a method of DataBase interface.
func (m *MockDB) InsertData(data Comment) error {
	log.Debug("Mock: insert data")
	return nil
}

// Close is a method of DataBase interface. Returns nil.
func (m *MockDB) Close() error {
	log.Debug("Mock: close DB")
	return nil
}
