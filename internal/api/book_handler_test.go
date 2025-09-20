package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type BookHandlerTestSuite struct {
	suite.Suite
	mockStore *store.MockBookStore
	handler   *BookHandler
}

func (s *BookHandlerTestSuite) SetupTest() {
	s.mockStore = new(store.MockBookStore)
	s.handler = NewBookHandler(s.mockStore)
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}

func (s *BookHandlerTestSuite) TestAddBooks_InvalidRquest() {
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleAddBook(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleAddBook_Success() {
	testUrl := "test"
	testImages := &store.BookImages{ThumbnailUrl: &testUrl}
	expectedBook := &store.Book{
		ID:        1,
		Title:     "1984",
		Authors:   []string{"Test"},
		Publisher: "test",
		ISBN13:    "TEST",
		Images:    *testImages,
	}
	s.mockStore.On("AddBook", expectedBook).Return(expectedBook, nil)

	body, _ := json.Marshal(expectedBook)
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleAddBook(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.mockStore.AssertExpectations(s.T())
}
