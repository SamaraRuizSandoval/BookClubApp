package store

import "github.com/stretchr/testify/mock"

type MockBookStore struct {
	mock.Mock
}

func (m *MockBookStore) AddBook(book *Book) (*Book, error) {
	args := m.Called(book)
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockBookStore) GetBookByID(id int64) (*Book, error) {
	args := m.Called(id)
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockBookStore) UpdateBook(book *Book) error {
	//TODO
	return nil
}
