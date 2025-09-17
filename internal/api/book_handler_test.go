package api

// import (
// 	"github.com/SamaraRuizSandoval/BookClubApp/internal/routes"
// 	"github.com/gin-gonic/gin"
// )

// func setupRouter() *gin.Engine {
//   router := gin.Default()
//   router.GET("/ping", func(c *gin.Context) {
//     c.String(200, "pong")
//   })
//   return router
// }

// func postUser(router *gin.Engine) *gin.Engine {
//   router.POST("/user/add", func(c *gin.Context) {
//     var user User
//     c.BindJSON(&user)
//     c.JSON(200, user)
//   })
//   return router
// }

// func main() {
//   router := setupRouter()
//   router = postUser(r)
//   router.Run(":8080")
// }

// //Test for code example above:

// package main

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
