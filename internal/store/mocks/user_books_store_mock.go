package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockUserBooksStore struct {
	mock.Mock
}

func (mubs *MockUserBooksStore) GetUserBooksByUserID(userID int64, status *string, page, limit int) ([]*store.BasicUserBook, error) {
	args := mubs.Called(userID, status, page, limit)
	return args.Get(0).([]*store.BasicUserBook), args.Error(1)
}

func (mubs *MockUserBooksStore) AddUserBook(userID, bookID int64, status string) (*store.UserBook, error) {
	args := mubs.Called(userID, bookID, status)
	return args.Get(0).(*store.UserBook), args.Error(1)
}

func (mubs *MockUserBooksStore) UpdateUserBook(userID, userBookID int64, req store.UpdateUserBookRequest) (*store.UserBook, error) {
	args := mubs.Called(userID, userBookID, req)
	return args.Get(0).(*store.UserBook), args.Error(1)
}

func (mubs *MockUserBooksStore) DeleteUserBook(userID, userBookID int64) error {
	args := mubs.Called(userID, userBookID)
	return args.Error(0)
}
