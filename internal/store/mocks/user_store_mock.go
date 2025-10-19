package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (mus *MockUserStore) CreateUser(user *store.User) (*store.User, error) {
	args := mus.Called(user)
	return args.Get(0).(*store.User), args.Error(1)
}

func (mus *MockUserStore) GetUserByUsername(username string) (*store.User, error) {
	args := mus.Called(username)
	return args.Get(0).(*store.User), args.Error(1)
}

func (mus *MockUserStore) UpdateUser(user *store.User) error {
	//TODO
	return nil
}

func (mus *MockUserStore) GetUserToken(scope, plainTextToken string) (*store.User, error) {
	args := mus.Called(scope, plainTextToken)
	return args.Get(0).(*store.User), args.Error(1)
}
