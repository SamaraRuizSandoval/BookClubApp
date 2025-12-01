package api

import (
	"database/sql"
	"errors"
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
		status = "wishlist"
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

// HandleUpdateUserBook godoc
// @Summary      Partially update a user-book relationship
// @Description  Updates only the fields provided in the JSON request.
// @Tags         user_books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "UserBook ID"
// @Param        data body store.UpdateUserBookRequest true "Fields to update"
// @Success      200 {object} store.UserBook
// @Failure      400 {object} HTTPError
// @Failure      404 {object} HTTPError
// @Failure      500 {object} HTTPError
// @Router       /user-books/{id} [patch]
func (h *UserBooksHandler) HandleUpdateUserBook(ctx *gin.Context) {
	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	idParam := ctx.Param("id")
	userBookID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// TODO: Only allow updating own user-books

	var req store.UpdateUserBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	if req.Status != nil {
		if *req.Status != "wishlist" &&
			*req.Status != "reading" &&
			*req.Status != "completed" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
	}

	if req.PagesRead != nil && *req.PagesRead < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pages_read must be >= 0"})
		return
	}

	if req.PercentageRead != nil {
		if *req.PercentageRead < 0 || *req.PercentageRead > 100 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "percentage_read must be between 0 and 100"})
			return
		}
	}

	// Delegate to store
	updated, err := h.userBooksStore.UpdateUserBook(
		user.ID, userBookID, req,
	)
	if err != nil {
		h.logger.Println("ERROR UpdateUserBook:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// HandleDeleteUserBook godoc
// @Summary      Delete a book from a user's shelf
// @Description  Removes the user-book entry entirely.
// @Tags         user_books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "UserBook ID"
// @Success      204 "Deleted successfully"
// @Failure      400 {object} HTTPError "Invalid ID"
// @Failure      404 {object} HTTPError "Not found"
// @Failure      500 {object} HTTPError "Internal server error"
// @Router       /user-books/{id} [delete]
func (h *UserBooksHandler) HandleDeleteUserBook(ctx *gin.Context) {
	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	idParam := ctx.Param("id")
	userBookID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.userBooksStore.DeleteUserBook(user.ID, userBookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user book not found"})
			return
		}

		h.logger.Println("ERROR DeleteUserBook:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
