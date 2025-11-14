package mocks

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockChapterCommentStore struct {
	mock.Mock
}

func (mccs *MockChapterCommentStore) AddComment(comment *store.ChapterComment, chapterID int64, userID int64) (*store.ChapterComment, error) {
	args := mccs.Called(comment)
	return args.Get(0).(*store.ChapterComment), args.Error(1)
}

func (mccs *MockChapterCommentStore) UpdateComment(comment *store.ChapterComment) error {
	args := mccs.Called(comment)
	return args.Error(0)
}

func (mccs *MockChapterCommentStore) GetCommentByID(id int64) (*store.ChapterComment, error) {
	args := mccs.Called(id)
	return args.Get(0).(*store.ChapterComment), args.Error(1)
}

func (mccs *MockChapterCommentStore) DeleteCommentByID(id int64) error {
	args := mccs.Called(id)
	return args.Error(0)
}

func (mccs *MockChapterCommentStore) GetCommentsByChapterID(chapterID int64, page, limit int) ([]*store.ChapterComment, int, error) {
	args := mccs.Called(chapterID, page, limit)
	return args.Get(0).([]*store.ChapterComment), args.Int(1), args.Error(2)
}
