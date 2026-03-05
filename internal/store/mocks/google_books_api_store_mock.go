package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockGoogleBookAPIStore struct {
	mock.Mock
}

func (m *MockGoogleBookAPIStore) SearchGoogleBooks(query string) ([]store.GoogleBookBasicInfo, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]store.GoogleBookBasicInfo), args.Error(1)
}
