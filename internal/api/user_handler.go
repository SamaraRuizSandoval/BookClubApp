package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HTTPError struct {
	Error string `json:"error"`
}

func (uh *UserHandler) validateRegisterRequest(req *RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// HandleGetUserByUsername godoc
// @Summary      Get a user by username
// @Description  Retrieves the details of a user by their username.
//
//	Provide a valid username as a path parameter. Returns the user object on success.
//
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username query string true "Username" example(johndoe)
// @Success      200 {object} store.User
// @Failure      400 {object} HTTPError "Error: Invalid or missing username"
// @Failure      404 {object} HTTPError "Error: User not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /users [get]
func (uh *UserHandler) HandleGetUserByUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	fmt.Println("USERNAME:", username)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	user, err := uh.userStore.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		uh.logger.Printf("ERROR: getUsername %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// RegisterUser godoc
// @Summary      Register a new user account
// @Description  Registers a new user in the system. Expects a JSON body containing username, email, and password. Returns the created account object on success.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body RegisterUserRequest true "Register user request"
// @Success      200 {object} store.User
// @Failure      400 {object} HTTPError "Error: Invalid Request"
// @Failure      409 {object} HTTPError "Error: Email or Username already exists"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /users [post]
func (uh *UserHandler) RegisterUser(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		uh.logger.Printf("ERROR: decodingRegisterUser %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := uh.validateRegisterRequest(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "user",
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hasing password %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	newUser, err := uh.userStore.CreateUser(user)
	if err != nil {
		switch err {
		case store.ErrEmailAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		case store.ErrUsernameAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		default:
			uh.logger.Printf("ERROR: createUser %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, newUser)
}

// RegisterAdminAccount godoc
// @Summary      Register a new admin account
// @Description  Registers a new admin in the system (only allowed by another admin). `Expects a JSON body containing username, email, and password. Returns the created account object on success.
// @Tags         admins
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body RegisterUserRequest true "Register user request"
// @Success      200 {object} store.User
// @Failure      400 {object} HTTPError "Error: Invalid Request"
// @Failure      409 {object} HTTPError "Error: Email or Username already exists"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /admins [post]
func (uh *UserHandler) RegisterAdminAccount(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		uh.logger.Printf("ERROR: decodingRegisterUser %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := uh.validateRegisterRequest(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "admin",
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hasing password %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	newUser, err := uh.userStore.CreateUser(user)
	if err != nil {
		switch err {
		case store.ErrEmailAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		case store.ErrUsernameAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		default:
			uh.logger.Printf("ERROR: createUser %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, newUser)
}
