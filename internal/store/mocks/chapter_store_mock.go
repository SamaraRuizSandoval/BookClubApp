package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockChapterStore struct {
	mock.Mock
}

func (mcs *MockChapterStore) GetChapterByID(id int64) (*store.Chapter, error) {
	args := mcs.Called(id)
	return args.Get(0).(*store.Chapter), args.Error(1)
}
