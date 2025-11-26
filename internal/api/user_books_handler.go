package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserBooksHandler struct {
	userBooksStore store.UserBooksStore
	logger         *log.Logger
}

func NewUserBooksHandler(userBooksStore store.UserBooksStore, logger *log.Logger) *UserBooksHandler {
	return &UserBooksHandler{
		userBooksStore: userBooksStore,
		logger:         logger,
	}
}

type UserBooksResponse struct {
	UserBooks []*store.BasicUserBook `json:"user_books"`
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
}

// HandleGetUserBooks godoc
// @Summary      Get a user's books
// @Description  Retrieves the books for a given user. Optional `status` query parameter filters by reading status.
// @Tags         user_books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status query string false "Filter by status (wishlist|reading|completed)"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(20)
// @Success      200 {object} UserBooksResponse
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /users/{user_id}/books [get]
func (h *UserBooksHandler) HandleGetUserBooks(ctx *gin.Context) {
	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	page, limit, err := utils.ReadPaginationParams(ctx)
	if err != nil {
		h.logger.Printf("ERROR: readPaginationParams %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	statusQuery := ctx.Query("status")
	var statusPtr *string
	if statusQuery != "" {
		statusPtr = &statusQuery
	}

	userBooks, err := h.userBooksStore.GetUserBooksByUserID(user.ID, statusPtr, page, limit)
	if err != nil {
		h.logger.Printf("ERROR: GetUserBooksByUserID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, UserBooksResponse{
		UserBooks: userBooks,
		Page:      page,
		Limit:     limit,
	})
}

// HandleAddUserBook godoc
// @Summary      Add a book to a user's collection
// @Description  Adds a book to a user's collection with a specified status.
// @Tags         user_books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        book_id query int true "Book ID"
// @Param        status query string false "Filter by status (wishlist|reading|completed)"
// @Success      200 {object} UserBooksResponse
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /users/{user_id}/books [post]
func (h *UserBooksHandler) HandleAddUserBook(ctx *gin.Context) {
	// Implementation for adding a user book goes here
	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	bookID := ctx.Query("book_id")
	if bookID == "" {
		h.logger.Printf("ERROR: readIDParam %v", bookID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing book id"})
		return
	}

	bookIdVal, err := strconv.ParseInt(bookID, 10, 64)
	if err != nil {
		h.logger.Printf("ERROR: parseInt %v", bookID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid type book id"})
		return
	}

	status := ctx.Query("status")
	if status == "" {
		status = "wishlist" // default status
	}

	if status != "wishlist" && status != "reading" && status != "completed" {
		h.logger.Printf("ERROR: invalid status value %v", status)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	userBook, err := h.userBooksStore.AddUserBook(user.ID, bookIdVal, status)
	if err != nil {
		h.logger.Printf("ERROR: AddUserBook %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, userBook)
}
