package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

type BookHandler struct {
	bookStore store.BookStore
	logger    *log.Logger
}

func NewBookHandler(bookStore store.BookStore, logger *log.Logger) *BookHandler {
	return &BookHandler{
		bookStore: bookStore,
		logger:    logger,
	}
}

// HandleGetBookByID godoc
// @Summary      Get a book by id
// @Description  Retrieves the details of a book by their id.
//
//	Provide a valid id. Returns the book object on success.
//
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id path int true "The id of the book"
// @Success      200 {object} store.Book
// @Failure      404 {object} HTTPError "Error: Book not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /books/{id} [get]
func (bh *BookHandler) HandleGetBookByID(ctx *gin.Context) {
	bookID, err := utils.ReadIDParam(ctx)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		bh.logger.Printf("ERROR: getBookByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if book == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

type AddBookRequest struct {
	Title         string           `json:"title" example:"The Hobbit"`
	Authors       []string         `json:"authors" example:"J.R.R. Tolkien"`
	Publisher     string           `json:"publisher" example:"George Allen & Unwin"`
	PublishedDate store.JSONDate   `json:"published_date" example:"1937-09-21"`
	Description   *string          `json:"description,omitempty" example:"A fantasy novel..."`
	PageCount     *int             `json:"page_count,omitempty" example:"310"`
	ISBN13        string           `json:"isbn_13" example:"9780261102217"`
	ISBN10        *string          `json:"isbn_10,omitempty" example:"0261102214"`
	Images        store.BookImages `json:"book_images"`
	Chapters      []store.Chapter  `json:"chapters"`
}

// HandleAddBook godoc
// @Summary      Add a book
// @Description  Registers a book in the system.
//
//	Expects a body with the book information. Returns the created book on success.
//
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body AddBookRequest true "Add book request"
// @Success      200 {object} store.Book
// @Failure      401 {object} HTTPError "Error: Unauthorized"
// @Failure      409 {object} HTTPError "Error: Duplicate record"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /books [post]
func (bh *BookHandler) HandleAddBook(ctx *gin.Context) {
	var req AddBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		bh.logger.Printf("ERROR: decodingAddBook %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	book := store.Book{
		Title:         req.Title,
		Authors:       req.Authors,
		Publisher:     req.Publisher,
		PublishedDate: req.PublishedDate,
		Description:   req.Description,
		PageCount:     req.PageCount,
		ISBN13:        req.ISBN13,
		ISBN10:        req.ISBN10,
		Images:        req.Images,
		Chapters:      req.Chapters,
	}

	addedBook, err := bh.bookStore.AddBook(&book)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, gin.H{"error": "book with this ISBN already exists"})
				return
			}
		}
		bh.logger.Printf("ERROR: addBook %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, addedBook)
}

// HandleUpdateBookByID godoc
// @Summary      Update a book
// @Description  Updates a book's information in the system.
//
//	Expects a body with the book information. Returns the updated book on success.
//
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Book ID"
// @Param        request body AddBookRequest true "Add book request"
// @Success      200 {object} store.Book
// @Failure      401 {object} HTTPError "Error: Unauthorized"
// @Failure      404 {object} HTTPError "Error: User not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /books/{id} [put]
func (bh *BookHandler) HandleUpdateBookByID(ctx *gin.Context) {
	bookID, err := utils.ReadIDParam(ctx)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	existingBook, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		bh.logger.Printf("ERROR: getBookByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if existingBook == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var req AddBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		bh.logger.Printf("ERROR: updateBookByID %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	book := store.Book{
		Title:         req.Title,
		Authors:       req.Authors,
		Publisher:     req.Publisher,
		PublishedDate: req.PublishedDate,
		Description:   req.Description,
		PageCount:     req.PageCount,
		ISBN13:        req.ISBN13,
		ISBN10:        req.ISBN10,
		Images:        req.Images,
		Chapters:      req.Chapters,
	}

	err = bh.bookStore.UpdateBook(&book)
	if err != nil {
		bh.logger.Printf("ERROR: updateBookByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	updatedBook, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		bh.logger.Printf("ERROR: getBookByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

func (bh *BookHandler) HandleDeleteBookByID(ctx *gin.Context) {
	bookID, err := utils.ReadIDParam(ctx)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	if err := bh.bookStore.DeleteBookByID(bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		} else {
			bh.logger.Printf("ERROR: deleteBookByID %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}

type PaginatedBooksResponse struct {
	Books      []*store.Book `json:"books"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalItems int           `json:"total_items"`
	TotalPages int           `json:"total_pages"`
}

// HandleGetAllBooks godoc
// @Summary      Get all books
// @Description  Retrieves all books with pagination.
//
//	Provide a valid page and limit parameters. Returns the paginated books object on success.
//
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(20)
// @Success      200 {object} PaginatedBooksResponse
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      404 {object} HTTPError "Error: Book not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /books [get]
func (bh *BookHandler) HandleGetAllBooks(ctx *gin.Context) {
	page, limit, err := utils.ReadPaginationParams(ctx)
	if err != nil {
		bh.logger.Printf("ERROR: readPaginationParams %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	books, total, err := bh.bookStore.GetAllBooks(page, limit)
	if err != nil {
		bh.logger.Printf("ERROR: getAllBooks %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	totalPages := (total + limit - 1) / limit

	ctx.JSON(http.StatusOK, PaginatedBooksResponse{
		Books:      books,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
