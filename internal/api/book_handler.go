package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/utils"
	"github.com/gin-gonic/gin"
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

func (bh *BookHandler) HandleAddBook(ctx *gin.Context) {
	var book store.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		bh.logger.Printf("ERROR: decodingAddBook %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	addedBook, err := bh.bookStore.AddBook(&book)
	// TODO: Check if book already exists
	if err != nil {
		bh.logger.Printf("ERROR: addBook %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, addedBook)
}

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

	var updatebookRequest store.Book
	if err := ctx.ShouldBindJSON(&updatebookRequest); err != nil {
		bh.logger.Printf("ERROR: updateBookByID %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatebookRequest.ID = existingBook.ID
	err = bh.bookStore.UpdateBook(&updatebookRequest)
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
