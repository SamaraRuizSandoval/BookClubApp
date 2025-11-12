package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/tokens"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

// HandleCreateToken godoc
// @Summary      Login
// @Description  Authenticates a user in the system. Expects a JSON body containing username and password. Returns a bearer token on success.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body createTokenRequest true "Authentication Request"
// @Success      200 {object} tokens.Token
// @Failure      401 {object} HTTPError "Error: Invalid Credentials"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /tokens/authentication [post]
func (th *TokenHandler) HandleCreateToken(ctx *gin.Context) {
	var req createTokenRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		th.logger.Printf("ERROR: decodingCreateToken %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		th.logger.Printf("ERROR: GetUserByUsername %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: PasswordHash.Matches %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("ERROR: CreatingToken %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"auth_token": token})
}
