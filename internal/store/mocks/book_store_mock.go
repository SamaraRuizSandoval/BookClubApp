package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockBookStore struct {
	mock.Mock
}

func (m *MockBookStore) AddBook(book *store.Book) (*store.Book, error) {
	args := m.Called(book)
	return args.Get(0).(*store.Book), args.Error(1)
}

func (m *MockBookStore) GetBookByID(id int64) (*store.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.Book), args.Error(1)
}

func (m *MockBookStore) UpdateBook(book *store.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookStore) DeleteBookByID(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookStore) GetAllBooks(page, limit int) ([]*store.Book, int, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*store.Book), args.Int(1), args.Error(2)
}
