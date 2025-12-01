package api

import (
	"bytes"
	"database/sql"
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

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_InvalidID() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	reqBody := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("PATCH", "/user-books/abc", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_InvalidJSON() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	// invalid JSON
	reqBody := bytes.NewBufferString(`{invalid`)
	req, _ := http.NewRequest("PATCH", "/user-books/10", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_InvalidStatusValue() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	body := map[string]interface{}{"status": "bad"}
	b, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PATCH", "/user-books/10", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_NegativePagesRead() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	body := map[string]interface{}{"pages_read": -5}
	b, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PATCH", "/user-books/10", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_PercentageOutOfRange() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	body := map[string]interface{}{"percentage_read": 150}
	b, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PATCH", "/user-books/10", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_StoreError() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	status := "reading"
	body := map[string]interface{}{"status": status}
	b, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PATCH", "/user-books/12", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 9})
	ctx.Params = gin.Params{
		gin.Param{Key: "id", Value: "12"},
	}

	reqStruct := store.UpdateUserBookRequest{Status: &status}
	suite.MockStore.On("UpdateUserBook", int64(9), int64(12), reqStruct).Return((*store.UserBook)(nil), errors.New("fail"))

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleUpdateUserBook_Success() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	status := "completed"
	pages := 100
	perc := 100.0
	body := map[string]interface{}{"status": status, "pages_read": pages, "percentage_read": perc}
	b, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PATCH", "/user-books/77", reqBody)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 11})
	ctx.Params = gin.Params{
		gin.Param{Key: "id", Value: "77"},
	}

	reqStruct := store.UpdateUserBookRequest{Status: &status, PagesRead: &pages, PercentageRead: &perc}
	ub := &store.UserBook{ID: 77, UserID: 11, BookID: 5, Status: status, UpdatedAt: store.JSONDate(time.Now())}
	suite.MockStore.On("UpdateUserBook", int64(11), int64(77), reqStruct).Return(ub, nil)

	suite.UserBooksHandler.HandleUpdateUserBook(ctx)
	suite.Equal(http.StatusOK, w.Code)
	var got store.UserBook
	err := json.Unmarshal(w.Body.Bytes(), &got)
	suite.NoError(err)
	suite.Equal(int64(77), got.ID)
	suite.Equal(status, got.Status)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleDeleteUserBook_InvalidID() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("DELETE", "/user-books/abc", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 1})

	suite.UserBooksHandler.HandleDeleteUserBook(ctx)
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserBooksHandlerTestSuite) TestHandleDeleteUserBook_NotFound() {
	suite.MockStore.On("DeleteUserBook", int64(2), int64(55)).Return(sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodDelete, "/user-books/55", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req
	ctx.Set("user", &store.User{ID: 2})
	ctx.Params = gin.Params{
		gin.Param{Key: "id", Value: "55"},
	}

	suite.UserBooksHandler.HandleDeleteUserBook(ctx)
	suite.Equal(http.StatusNotFound, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleDeleteUserBook_StoreError() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/user-books/55", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 2})
	ctx.Params = gin.Params{
		gin.Param{Key: "id", Value: "55"},
	}

	suite.MockStore.On("DeleteUserBook", int64(2), int64(55)).Return(errors.New("boom"))

	suite.UserBooksHandler.HandleDeleteUserBook(ctx)
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func (suite *UserBooksHandlerTestSuite) TestHandleDeleteUserBook_Success() {
	suite.MockStore.On("DeleteUserBook", int64(2), int64(55)).Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/user-books/55", nil)
	ctx.Request = req
	ctx.Set("user", &store.User{ID: 2})
	ctx.Params = gin.Params{
		gin.Param{Key: "id", Value: "55"},
	}

	suite.UserBooksHandler.HandleDeleteUserBook(ctx)
	suite.Equal(http.StatusOK, w.Code)
	suite.MockStore.AssertExpectations(suite.T())
}

func TestUserBooksHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserBooksHandlerTestSuite))
}
