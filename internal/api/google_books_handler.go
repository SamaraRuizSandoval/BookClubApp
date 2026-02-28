package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

type GoogleBookApiHandler struct {
	googleBookAPI store.GoogleBookAPI
	logger        *log.Logger
}

func NewGoogleBookApiHandler(googleAPiStore store.GoogleBookAPI, logger *log.Logger) *GoogleBookApiHandler {
	return &GoogleBookApiHandler{
		googleBookAPI: googleAPiStore,
		logger:        logger,
	}
}

type GoogleBooksResponse struct {
	GoogleBooks []*store.GoogleBooksAPISearch `json:"google_books"`
}

// HandleSearchGoogleBooks godoc
// @Summary      Search Google Books
// @Description  Searches for books in the Google Books API.
// @Tags         google_books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        query query string false "Search query (e.g. title, author)"
// @Success      200 {object} []store.Book "Successful response with list of books"
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /api/books [get]
func (gbh *GoogleBookApiHandler) HandleSearchGoogleBooks(ctx *gin.Context) {
	newBooks := make([]store.Book, 0)
	query := ctx.Query("query")
	if query == "" {
		gbh.logger.Println("ERROR: missing search query")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing search query",
		})
		return
	}

	books, err := gbh.googleBookAPI.SearchGoogleBooks(query)
	if err != nil {
		gbh.logger.Printf("ERROR: searching google books: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch books",
		})
		return
	}

	for index, book := range books {

		mapped, err := mapGoogleToBook(book, index)
		if err != nil {
			gbh.logger.Printf("ERROR: mapping google book to internal book: %v", err)
			continue
		}
		newBooks = append(newBooks, *mapped)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": newBooks,
	})
}

func mapGoogleToBook(book store.GoogleBookBasicInfo, index int) (*store.Book, error) {
	var Isbn10, Isbn13, thumbnail, smallThumbnail string

	if book.VolumeInfo == nil {
		return nil, fmt.Errorf("missing volume info for book: %s", book.ID)
	}
	var descriptionPtr *string
	if book.VolumeInfo.Description != "" {
		description := book.VolumeInfo.Description
		descriptionPtr = &description
	}

	var pageCountPtr *int
	if book.VolumeInfo.PageCount != 0 {
		pageCount := book.VolumeInfo.PageCount
		pageCountPtr = &pageCount
	}

	if book.VolumeInfo.ImageLinks != nil {
		thumbnail = book.VolumeInfo.ImageLinks.Thumbnail
		smallThumbnail = book.VolumeInfo.ImageLinks.SmallThumbnail
	}

	publishedDate, err := parseGoogleDate(book.VolumeInfo.PublishedDate)
	if err != nil {
		publishedDate = store.JSONDate{}
	}

	if book.VolumeInfo.IndustryIdentifiers != nil {
		for _, identifier := range book.VolumeInfo.IndustryIdentifiers {
			switch identifier.Type {
			case "ISBN_10":
				Isbn10 = identifier.Identifier
			case "ISBN_13":
				Isbn13 = identifier.Identifier
			}
		}
	}

	newBook := store.Book{
		ID:            int64(index),
		Title:         book.VolumeInfo.Title,
		Authors:       book.VolumeInfo.Authors,
		PublishedDate: publishedDate,
		Description:   descriptionPtr,
		PageCount:     pageCountPtr,
		ISBN10:        &Isbn10,
		ISBN13:        Isbn13,

		Images: store.BookImages{
			ThumbnailUrl: &thumbnail,
			SmallUrl:     &smallThumbnail,
		},
	}
	return &newBook, nil
}

func parseGoogleDate(dateStr string) (store.JSONDate, error) {
	if dateStr == "" {
		return store.JSONDate{}, nil
	}

	layouts := []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}

	var parsed time.Time
	var err error

	for _, layout := range layouts {
		parsed, err = time.Parse(layout, dateStr)
		if err == nil {
			return store.JSONDate(parsed), nil
		}
	}

	return store.JSONDate{}, err
}
