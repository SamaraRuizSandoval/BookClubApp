package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	mockUserStore *mocks.MockUserStore
	userHandler   *UserHandler
}

func (uhs *UserHandlerTestSuite) SetupTest() {
	uhs.mockUserStore = new(mocks.MockUserStore)

	var buf bytes.Buffer
	logger := log.New(&buf, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	uhs.userHandler = NewUserHandler(uhs.mockUserStore, logger)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (uhs *UserHandlerTestSuite) TestHandleRegisterUser_InvalidJSON() {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	uhs.userHandler.HandleRegisterUser(ctx)

	uhs.Equal(http.StatusBadRequest, rec.Code)
	uhs.Contains(rec.Body.String(), "invalid request body")
}

func (uhs *UserHandlerTestSuite) TestHandleRegisterUser_MissingUsername() {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	body := map[string]string{
		"email":    "a@b.com",
		"password": "pass",
	}
	b, _ := json.Marshal(body)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	ctx.Request.Header.Set("Content-Type", "application/json")

	uhs.userHandler.HandleRegisterUser(ctx)

	uhs.Equal(http.StatusBadRequest, rec.Code)
	uhs.Contains(rec.Body.String(), "username is required")
}

func (uhs *UserHandlerTestSuite) TestHandleRegisterUser_InvalidEmail() {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	body := map[string]string{
		"username": "alice",
		"email":    "invalid-email",
		"password": "pass",
	}
	b, _ := json.Marshal(body)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	ctx.Request.Header.Set("Content-Type", "application/json")

	uhs.userHandler.HandleRegisterUser(ctx)

	uhs.Equal(http.StatusBadRequest, rec.Code)
	uhs.Contains(rec.Body.String(), "invalid email format")
}

func (uhs *UserHandlerTestSuite) TestHandleRegisterUser_MissingPassword() {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	body := map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
	}
	b, _ := json.Marshal(body)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	ctx.Request.Header.Set("Content-Type", "application/json")

	uhs.userHandler.HandleRegisterUser(ctx)

	uhs.Equal(http.StatusBadRequest, rec.Code)
	uhs.Contains(rec.Body.String(), "password is required")
}

func (uhs *UserHandlerTestSuite) TestHandleRegisterUser_Success() {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	reqBody := map[string]string{
		"username": "john",
		"email":    "john@example.com",
		"password": "strong-password",
	}
	b, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	ctx.Request.Header.Set("Content-Type", "application/json")

	expectedUser := &store.User{
		Username: "john",
		Email:    "john@example.com",
		Role:     "user",
	}

	uhs.mockUserStore.
		On("CreateUser", mock.MatchedBy(func(u *store.User) bool {
			return u != nil && u.Username == "john" && u.Email == "john@example.com"
		})).
		Return(expectedUser, nil)

	uhs.userHandler.HandleRegisterUser(ctx)

	uhs.Equal(http.StatusOK, rec.Code)
	uhs.Contains(rec.Body.String(), `"username":"john"`)
	uhs.Contains(rec.Body.String(), `"email":"john@example.com"`)

	uhs.mockUserStore.AssertExpectations(uhs.T())
}
