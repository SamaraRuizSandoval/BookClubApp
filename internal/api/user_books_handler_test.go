package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UserBooksHandlerTestSuite struct {
	suite.Suite
	UserBooksHandler *UserBooksHandler
	MockStore        *mocks.MockUserBooksStore
}

func (suite *UserBooksHandlerTestSuite) SetupTest() {
	suite.MockStore = new(mocks.MockUserBooksStore)
	var buf bytes.Buffer
	logger := log.New(&buf, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	suite.UserBooksHandler = NewUserBooksHandler(suite.MockStore, logger)
}

func (suite *UserBooksHandlerTestSuite) TestHandleGetUserBooks_InvalidPagination() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/?page=0", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleGetUserBooks(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleGetUserBooks_StoreError() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/?page=1&limit=20", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.MockStore.On("GetUserBooksByUserID", int64(1), (*string)(nil), 1, 20).Return([]*store.BasicUserBook(nil), errors.New("boom"))

	suite.UserBooksHandler.HandleGetUserBooks(ctx)
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleGetUserBooks_Success() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/?page=2&limit=5", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 42})

	ub := &store.BasicUserBook{ID: 10, UserID: 42, Status: "wishlist", UpdatedAt: store.JSONDate(time.Now())}
	suite.MockStore.On("GetUserBooksByUserID", int64(42), (*string)(nil), 2, 5).Return([]*store.BasicUserBook{ub}, nil)

	suite.UserBooksHandler.HandleGetUserBooks(ctx)
	suite.Equal(http.StatusOK, w.Code)

	var resp UserBooksResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(2, resp.Page)
	suite.Equal(5, resp.Limit)
	suite.Len(resp.UserBooks, 1)
	suite.Equal(int64(10), resp.UserBooks[0].ID)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleAddUserBook_MissingBookID() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleAddUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleAddUserBook_InvalidBookID() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/?book_id=abc", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleAddUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleAddUserBook_InvalidStatus() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/?book_id=1&status=invalid", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleAddUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleAddUserBook_StoreError() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/?book_id=2", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 7})

	suite.MockStore.On("AddUserBook", int64(7), int64(2), "wishlist").Return((*store.UserBook)(nil), errors.New("boom"))

	suite.UserBooksHandler.HandleAddUserBook(ctx)
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleAddUserBook_SuccessWithDefaultStatus() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/?book_id=3", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 8})

	ub := &store.UserBook{ID: 5, UserID: 8, BookID: 3, Status: "wishlist", UpdatedAt: store.JSONDate(time.Now())}
	suite.MockStore.On("AddUserBook", int64(8), int64(3), "wishlist").Return(ub, nil)

	suite.UserBooksHandler.HandleAddUserBook(ctx)
	suite.Equal(http.StatusOK, w.Code)

	var got store.UserBook
	err := json.Unmarshal(w.Body.Bytes(), &got)
	suite.NoError(err)
	suite.Equal(int64(5), got.ID)
	suite.Equal("wishlist", got.Status)
	suite.MockStore.AssertExpectations(suite.T())
}

func TestUserBooksHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserBooksHandlerTestSuite))
}
