package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ChapterCommentHandlerTestSuite struct {
	suite.Suite
	mockStore        *mocks.MockChapterCommentStore
	mockChapterStore *mocks.MockChapterStore
	handler          *ChapterCommentHandler
}

func (s *ChapterCommentHandlerTestSuite) SetupTest() {
	s.mockStore = new(mocks.MockChapterCommentStore)
	s.mockChapterStore = new(mocks.MockChapterStore)
	var buf bytes.Buffer
	logger := log.New(&buf, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	s.handler = NewChapterCommentHandler(s.mockStore, s.mockChapterStore, logger)
}

func TestChapterCommentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterCommentHandlerTestSuite))
}

// --- Add Comment ---
func (s *ChapterCommentHandlerTestSuite) TestHandleAddComment_InvalidChapterID() {
	req, _ := http.NewRequest(http.MethodPost, "/chapters/abc/comments", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "abc"}}

	s.handler.HandleAddComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleAddComment_InvalidJSON() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	req, _ := http.NewRequest(http.MethodPost, "/chapters/1/comments", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleAddComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleAddComment_ErrorFromStore() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, fmt.Errorf("boom"))
	// s.mockStore.On("AddComment", mock.Anything).Return(&store.ChapterComment{}, fmt.Errorf("boom"))

	body, _ := json.Marshal(map[string]string{"body": "hi"})
	req, _ := http.NewRequest(http.MethodPost, "/chapters/1/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleAddComment(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleAddComment_Success() {
	returned := &store.ChapterComment{Body: "ok", ID: 2, UserID: 1, ChapterID: 1}
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	s.mockStore.On("AddComment", mock.Anything).Return(returned, nil)

	body, _ := json.Marshal(map[string]string{"body": "ok"})
	req, _ := http.NewRequest(http.MethodPost, "/chapters/1/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleAddComment(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "ok")
	s.mockStore.AssertExpectations(s.T())
}

// --- Get Comment ---
func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentById_InvalidChapterID() {
	req, _ := http.NewRequest(http.MethodGet, "/chapters/abc/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "abc"}}

	s.handler.HandleGetCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentById_InvalidID() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "abc"},
	}

	s.handler.HandleGetCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentById_NotFound() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	s.mockStore.On("GetCommentByID", int64(1)).Return(&store.ChapterComment{}, sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleGetCommentById(ctx)

	s.Equal(http.StatusNotFound, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentById_ChapterMismatch() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	c := &store.ChapterComment{ID: 1, Body: "x", ChapterID: 2}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleGetCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentById_Success() {
	c := &store.ChapterComment{ID: 1, Body: "hello", ChapterID: 1}
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleGetCommentById(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "hello")
	s.mockStore.AssertExpectations(s.T())
}

// --- Update Comment ---
func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_InvalidChapterID() {
	req, _ := http.NewRequest(http.MethodPut, "/chapters/abc/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "abc"}}

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_InvalidID() {
	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "abc"},
	}

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_NotFound() {
	s.mockStore.On("GetCommentByID", int64(1)).Return(&store.ChapterComment{}, sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusNotFound, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_ChapterMismatch() {
	c := &store.ChapterComment{ID: 1, ChapterID: 2, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_InvalidJSON() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_Unauthorized() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 2}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	body, _ := json.Marshal(map[string]string{"body": "edited"})
	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_ErrorUpdating() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)
	s.mockStore.On("UpdateComment", mock.Anything).Return(fmt.Errorf("boom"))

	body, _ := json.Marshal(map[string]string{"body": "edited"})
	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleUpdateComment_Success() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)
	s.mockStore.On("UpdateComment", mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]string{"body": "edited"})
	req, _ := http.NewRequest(http.MethodPut, "/chapters/1/comments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleUpdateComment(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

// --- Delete Comment ---
func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_InvalidChapterID() {
	req, _ := http.NewRequest(http.MethodDelete, "/chapters/abc/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "abc"}}

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_InvalidID() {
	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/abc", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "abc"},
	}

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_NotFound() {
	s.mockStore.On("GetCommentByID", int64(1)).Return(&store.ChapterComment{}, sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusNotFound, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_ChapterMismatch() {
	c := &store.ChapterComment{ID: 1, ChapterID: 2, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_Unauthorized() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 2}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_ErrorDeleting() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)
	s.mockStore.On("DeleteCommentByID", int64(1)).Return(fmt.Errorf("boom"))

	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleDeleteCommentById_Success() {
	c := &store.ChapterComment{ID: 1, ChapterID: 1, UserID: 1}
	s.mockStore.On("GetCommentByID", int64(1)).Return(c, nil)
	s.mockStore.On("DeleteCommentByID", int64(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/chapters/1/comments/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "chapter_id", Value: "1"},
		gin.Param{Key: "id", Value: "1"},
	}
	ctx.Set("user", &store.User{ID: 1})

	s.handler.HandleDeleteCommentById(ctx)

	s.Equal(http.StatusOK, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

// --- Get Comments By Chapter ---
func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentsByChapterID_InvalidChapterID() {
	req, _ := http.NewRequest(http.MethodGet, "/chapters/abc/comments", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "abc"}}

	s.handler.HandleGetCommentsByChapterID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentsByChapterID_InvalidPagination() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments?page=0&limit=10", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}

	s.handler.HandleGetCommentsByChapterID(ctx)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentsByChapterID_ErrorFromStore() {
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	s.mockStore.On("GetCommentsByChapterID", int64(1), 1, 20).Return([]*store.ChapterComment{}, 0, fmt.Errorf("boom"))

	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}

	s.handler.HandleGetCommentsByChapterID(ctx)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.mockStore.AssertExpectations(s.T())
}

func (s *ChapterCommentHandlerTestSuite) TestHandleGetCommentsByChapterID_Success() {
	comments := []*store.ChapterComment{
		{ID: 1, Body: "c1", ChapterID: 1},
		{ID: 2, Body: "c2", ChapterID: 1},
	}
	total := 25
	s.mockChapterStore.On("GetChapterByID", int64(1)).Return(&store.Chapter{ID: 1}, nil)
	s.mockStore.On("GetCommentsByChapterID", int64(1), 1, 20).Return(comments, total, nil)

	req, _ := http.NewRequest(http.MethodGet, "/chapters/1/comments", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{gin.Param{Key: "chapter_id", Value: "1"}}

	s.handler.HandleGetCommentsByChapterID(ctx)

	s.Equal(http.StatusOK, w.Code)
	// total pages = (25 + 20 -1) / 20 = 2
	s.Contains(w.Body.String(), "c1")
	s.Contains(w.Body.String(), "total_pages")
	s.mockStore.AssertExpectations(s.T())
}
