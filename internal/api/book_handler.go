package api

import (
	"fmt"
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	bookID, err := strconv.ParseInt(paramsBookId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	ctx.String(http.StatusOK, "This is the book id: %d\n", bookID)
}

func (bh *BookHandler) HandleAddBook(ctx *gin.Context) {
	var book store.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	addedBook, err := bh.bookStore.AddBook(&book)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add book"})
		return
	}

	ctx.JSON(http.StatusOK, addedBook)
}
