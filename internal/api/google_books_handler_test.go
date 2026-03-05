package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GoogleBooksHandlerTest struct {
	suite.Suite
	mockGoogleBookAPI *mocks.MockGoogleBookAPIStore
	handler           *GoogleBookApiHandler
}

func (suite *GoogleBooksHandlerTest) SetupTest() {
	suite.mockGoogleBookAPI = new(mocks.MockGoogleBookAPIStore)
	var buf bytes.Buffer
	logger := log.New(&buf, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	suite.handler = NewGoogleBookApiHandler(suite.mockGoogleBookAPI, logger)
}

func TestGoogleBooksHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GoogleBooksHandlerTest))
}

// --- HandleSearchGoogleBooks Tests ---
func (s *GoogleBooksHandlerTest) TestHandleSearchGoogleBooks_MissingQuery() {
	req, _ := http.NewRequest(http.MethodGet, "/api/books", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleSearchGoogleBooks(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "missing search query")
}

func (s *GoogleBooksHandlerTest) TestHandleSearchGoogleBooks_ErrorFromAPI() {
	s.mockGoogleBookAPI.On("SearchGoogleBooks", "test").Return(nil, fmt.Errorf("api error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/books?q=test", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleSearchGoogleBooks(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "failed to fetch books")
	s.mockGoogleBookAPI.AssertExpectations(s.T())
}

func (s *GoogleBooksHandlerTest) TestHandleSearchGoogleBooks_SuccessWithValidBooks() {
	books := []store.GoogleBookBasicInfo{
		{
			ID: "book1",
			VolumeInfo: &store.VolumeInfo{
				Title: "Test Book 1",
			},
		},
		{
			ID: "book2",
			VolumeInfo: &store.VolumeInfo{
				Title: "Test Book 2",
			},
		},
	}
	s.mockGoogleBookAPI.On("SearchGoogleBooks", "test").Return(books, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/books?q=test", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleSearchGoogleBooks(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Test Book 1")
	s.Contains(w.Body.String(), "Test Book 2")
	s.mockGoogleBookAPI.AssertExpectations(s.T())
}

func (s *GoogleBooksHandlerTest) TestHandleSearchGoogleBooks_SkipsInvalidBooks() {
	books := []store.GoogleBookBasicInfo{
		{
			ID: "book1",
			VolumeInfo: &store.VolumeInfo{
				Title: "Valid Book",
			},
		},
		{
			ID:         "book2",
			VolumeInfo: nil, // This will cause mapping error
		},
	}
	s.mockGoogleBookAPI.On("SearchGoogleBooks", "test").Return(books, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/books?q=test", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleSearchGoogleBooks(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Valid Book")
	s.mockGoogleBookAPI.AssertExpectations(s.T())
}

func (s *GoogleBooksHandlerTest) TestHandleSearchGoogleBooks_EmptyResults() {
	books := []store.GoogleBookBasicInfo{}
	s.mockGoogleBookAPI.On("SearchGoogleBooks", "nonexistent").Return(books, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/books?q=nonexistent", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	s.handler.HandleSearchGoogleBooks(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "books")
	s.mockGoogleBookAPI.AssertExpectations(s.T())
}

// --- mapGoogleToBook Tests ---
func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_MissingVolumeInfo() {
	book := store.GoogleBookBasicInfo{
		ID:         "test-id",
		VolumeInfo: nil,
	}

	result, err := mapGoogleToBook(book, 0)

	s.Nil(result)
	s.NotNil(err)
	s.Contains(err.Error(), "missing volume info")
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_ValidBookWithAllFields() {
	thumbnail := "http://example.com/thumb.jpg"
	smallThumb := "http://example.com/small.jpg"

	book := store.GoogleBookBasicInfo{
		ID: "valid-id",
		VolumeInfo: &store.VolumeInfo{
			Title:         "Test Book",
			Authors:       []string{"Author 1", "Author 2"},
			Publisher:     "Test Publisher",
			PublishedDate: "2023-05-15",
			Description:   "A detailed description",
			PageCount:     350,
			IndustryIdentifiers: []store.IndustryIdentifier{
				{Type: "ISBN_10", Identifier: "1234567890"},
				{Type: "ISBN_13", Identifier: "9781234567890"},
			},
			ImageLinks: &store.ImageLinks{
				Thumbnail:      thumbnail,
				SmallThumbnail: smallThumb,
			},
		},
	}

	result, err := mapGoogleToBook(book, 5)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(int64(5), result.ID)
	s.Equal("Test Book", result.Title)
	s.Equal(2, len(result.Authors))
	s.Equal("Test Publisher", result.Publisher)
	s.NotNil(result.Description)
	s.Equal("A detailed description", *result.Description)
	s.NotNil(result.PageCount)
	s.Equal(350, *result.PageCount)
	s.NotNil(result.ISBN10)
	s.Equal("1234567890", *result.ISBN10)
	s.Equal("9781234567890", result.ISBN13)
	s.NotNil(result.Images.ThumbnailUrl)
	s.Equal(thumbnail, *result.Images.ThumbnailUrl)
	s.NotNil(result.Images.SmallUrl)
	s.Equal(smallThumb, *result.Images.SmallUrl)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_NoDescription() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title: "Test Book",
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.Nil(result.Description)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_NoPageCount() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title:     "Test Book",
			PageCount: 0,
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.Nil(result.PageCount)
}

// func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_NoImageLinks() {
// 	book := store.GoogleBookBasicInfo{
// 		ID: "test-id",
// 		VolumeInfo: &store.VolumeInfo{
// 			Title:      "Test Book",
// 			ImageLinks: nil,
// 		},
// 	}

// 	result, err := mapGoogleToBook(book, 0)

// 	s.NoError(err)
// 	s.NotNil(result)
// 	s.Nil(result.Images.ThumbnailUrl)
// 	s.Nil(result.Images.SmallUrl)
// }

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_OnlyISBN10() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title: "Test Book",
			IndustryIdentifiers: []store.IndustryIdentifier{
				{Type: "ISBN_10", Identifier: "1234567890"},
			},
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.ISBN10)
	s.Equal("1234567890", *result.ISBN10)
	s.Equal("", result.ISBN13)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_OnlyISBN13() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title: "Test Book",
			IndustryIdentifiers: []store.IndustryIdentifier{
				{Type: "ISBN_13", Identifier: "9781234567890"},
			},
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("9781234567890", result.ISBN13)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_NoIdentifiers() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title:               "Test Book",
			IndustryIdentifiers: nil,
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_InvalidPublishedDate() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title:         "Test Book",
			PublishedDate: "invalid-date",
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	// When date parsing fails, it should still return empty JSONDate
	s.Equal(store.JSONDate(time.Time{}), result.PublishedDate)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_OtherIdentifierTypes() {
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title: "Test Book",
			IndustryIdentifiers: []store.IndustryIdentifier{
				{Type: "ISBN_10", Identifier: "1234567890"},
				{Type: "OTHER_TYPE", Identifier: "should-be-ignored"},
			},
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.ISBN10)
	s.Equal("1234567890", *result.ISBN10)
}

// --- parseGoogleDate Tests ---
func (s *GoogleBooksHandlerTest) TestParseGoogleDate_EmptyString() {
	result, err := parseGoogleDate("")

	s.NoError(err)
	s.Equal(store.JSONDate(time.Time{}), result)
}

func (s *GoogleBooksHandlerTest) TestParseGoogleDate_FullDate() {
	result, err := parseGoogleDate("2023-05-15")

	s.NoError(err)
	s.NotNil(result)
	expectedTime, _ := time.Parse("2006-01-02", "2023-05-15")
	s.Equal(store.JSONDate(expectedTime), result)
}

func (s *GoogleBooksHandlerTest) TestParseGoogleDate_YearMonth() {
	result, err := parseGoogleDate("2023-05")

	s.NoError(err)
	expectedTime, _ := time.Parse("2006-01", "2023-05")
	s.Equal(store.JSONDate(expectedTime), result)
}

func (s *GoogleBooksHandlerTest) TestParseGoogleDate_YearOnly() {
	result, err := parseGoogleDate("2023")

	s.NoError(err)
	expectedTime, _ := time.Parse("2006", "2023")
	s.Equal(store.JSONDate(expectedTime), result)
}

func (s *GoogleBooksHandlerTest) TestParseGoogleDate_InvalidString() {
	result, err := parseGoogleDate("not-a-date")

	s.NotNil(err)
	s.Equal(store.JSONDate{}, result)
}

func (s *GoogleBooksHandlerTest) TestMapGoogleToBook_PartialPageCount() {
	// Test with non-zero PageCount to ensure it's set correctly
	book := store.GoogleBookBasicInfo{
		ID: "test-id",
		VolumeInfo: &store.VolumeInfo{
			Title:     "Test Book",
			PageCount: 100,
		},
	}

	result, err := mapGoogleToBook(book, 0)

	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.PageCount)
	s.Equal(100, *result.PageCount)
}

// Additional test using testify/assert for alternative assertion style
func TestParseGoogleDateDirectly(t *testing.T) {
	// Test with assert instead of suite for variety
	result, err := parseGoogleDate("2020-12-25")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result2, err2 := parseGoogleDate("invalid")
	assert.Error(t, err2)
	assert.Equal(t, store.JSONDate{}, result2)
}
