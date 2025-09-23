package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
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

var expectedBook = &store.Book{
	ID:        1,
	Title:     "1984",
	Authors:   []string{"Test"},
	Publisher: "test",
	ISBN13:    "TEST",
}

func (s *BookHandlerTestSuite) TestHandleAddBook_InvalidRequest() {
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleAddBook(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleAddBook_ErrorAddingBook() {
	s.mockStore.On("AddBook", expectedBook).Return(expectedBook, fmt.Errorf("test error"))

	body, _ := json.Marshal(expectedBook)
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleAddBook(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleAddBook_Success() {
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

func (s *BookHandlerTestSuite) TestGetBookByID_InvalidRequest() {
	req, _ := http.NewRequest(http.MethodGet, "/book/abc", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleGetBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleGetBook_InvalidID_ParseError() {
	req, _ := http.NewRequest(http.MethodGet, "/book/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	s.handler.HandleGetBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleGetBook_ErrorGettingBook() {
	s.mockStore.On("GetBookByID", mock.Anything).Return(expectedBook, fmt.Errorf("test error"))

	req, _ := http.NewRequest(http.MethodGet, "/book/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	s.handler.HandleGetBookByID(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleGetBook_NotFound() {
	s.mockStore.On("GetBookByID", int64(1)).Return(nil, nil)

	req, _ := http.NewRequest(http.MethodGet, "/book/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	s.handler.HandleGetBookByID(ctx)

	s.Equal(http.StatusNotFound, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleGetBook_Success() {
	s.mockStore.On("GetBookByID", int64(1)).Return(expectedBook, nil)

	req, _ := http.NewRequest(http.MethodGet, "/book/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	s.handler.HandleGetBookByID(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "1984")
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestUpdateBookByID_InvalidRequest() {
	req, _ := http.NewRequest(http.MethodPut, "/book/abc", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleUpdateBook_InvalidID() {
	req, _ := http.NewRequest(http.MethodPut, "/books/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleUpdateBook_ErrorGettingBook() {
	s.mockStore.On("GetBookByID", mock.Anything).Return(expectedBook, fmt.Errorf("test error"))

	req, _ := http.NewRequest(http.MethodPut, "/books/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleUpdateBook_NotFound() {
	s.mockStore.On("GetBookByID", int64(1)).Return(nil, nil)

	req, _ := http.NewRequest(http.MethodGet, "/book/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusNotFound, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleUpdateBook_InvalidUpdateRequest() {
	s.mockStore.On("GetBookByID", int64(1)).Return(expectedBook, nil)

	req, _ := http.NewRequest(http.MethodGet, "/book/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleUpdateBook_ErrorUpdatingBook() {
	s.mockStore.On("GetBookByID", int64(1)).Return(expectedBook, nil)
	s.mockStore.On("UpdateBook", mock.Anything).Return(fmt.Errorf("test error"))

	body, _ := json.Marshal(expectedBook)
	req, _ := http.NewRequest(http.MethodGet, "/book/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	s.handler.HandleUpdateBookByID(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestDeleteBookByID_InvalidRequest() {
	req, _ := http.NewRequest(http.MethodDelete, "/book/abc", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleDeleteBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleDeleteBook_InvalidID() {
	req, _ := http.NewRequest(http.MethodDelete, "/book/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	s.handler.HandleDeleteBookByID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BookHandlerTestSuite) TestHandleDeleteBook_ErrorBookNotFound() {
	s.mockStore.On("DeleteBookByID", mock.Anything).Return(fmt.Errorf("no rows in result set"))

	req, _ := http.NewRequest(http.MethodDelete, "/books/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	s.handler.HandleDeleteBookByID(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *BookHandlerTestSuite) TestHandleDeleteBook_Success() {
	s.mockStore.On("DeleteBookByID", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/books/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	s.handler.HandleDeleteBookByID(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.mockStore.AssertExpectations(s.T())
}
