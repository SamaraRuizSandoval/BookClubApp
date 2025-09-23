package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookStore store.BookStore
}

func NewBookHandler(bookStore store.BookStore) *BookHandler {
	return &BookHandler{
		bookStore: bookStore,
	}
}

func (bh *BookHandler) HandleGetBookByID(ctx *gin.Context) {
	paramsBookId := ctx.Param("id")
	if paramsBookId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	bookID, err := strconv.ParseInt(paramsBookId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch book"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	addedBook, err := bh.bookStore.AddBook(&book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add book"})
		return
	}

	ctx.JSON(http.StatusOK, addedBook)
}

func (bh *BookHandler) HandleUpdateBookByID(ctx *gin.Context) {
	paramsBookId := ctx.Param("id")
	if paramsBookId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	bookID, err := strconv.ParseInt(paramsBookId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existingBook, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch book"})
		return
	}

	if existingBook == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var updatebookRequest store.Book
	if err := ctx.ShouldBindJSON(&updatebookRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updatebookRequest.ID = existingBook.ID
	err = bh.bookStore.UpdateBook(&updatebookRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "update book"})
		return
	}

	updatedBook, err := bh.bookStore.GetBookByID(bookID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch book"})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

func (bh *BookHandler) HandleDeleteBookByID(ctx *gin.Context) {
	paramsBookId := ctx.Param("id")
	if paramsBookId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	bookID, err := strconv.ParseInt(paramsBookId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := bh.bookStore.DeleteBookByID(bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
