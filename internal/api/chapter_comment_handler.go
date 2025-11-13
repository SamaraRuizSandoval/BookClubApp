package api

import (
	"log"
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type ChapterCommentHandler struct {
	chapterCommentStore store.ChapterCommentStore
	logger              *log.Logger
}

func NewChapterCommentHandler(commentStore store.ChapterCommentStore, logger *log.Logger) *ChapterCommentHandler {
	return &ChapterCommentHandler{
		chapterCommentStore: commentStore,
		logger:              logger,
	}
}

type AddChapterCommentRequest struct {
	Body string `json:"body" example:"I loved this chapter"`
}

// HandleAddComment godoc
// @Summary      Add a comment to a book's chapter
// @Description  Adds a user comment on a book's chapter. Expects a JSON body containing the body of the comment. Returns the created comment object on success.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Chapter ID"
// @Param        request body AddChapterCommentRequest true "Register comment request"
// @Success      200 {object} store.ChapterComment
// @Failure      400 {object} HTTPError "Error: Invalid Request"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{id}/comments  [post]
func (ch *ChapterCommentHandler) HandleAddComment(ctx *gin.Context) {
	chapterID, err := utils.ReadIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	//TODO: Check that chapter exist

	var req AddChapterCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ch.logger.Printf("ERROR: decodingChapterComment %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	userID := int64(user.ID)

	chapterComment := store.ChapterComment{
		Body: req.Body,
	}

	addedComment, err := ch.chapterCommentStore.AddComment(&chapterComment, chapterID, userID)
	if err != nil {
		ch.logger.Printf("ERROR: addBook %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, addedComment)
}
