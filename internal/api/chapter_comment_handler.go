package api

import (
	"database/sql"
	"errors"
	"fmt"
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
// @Param        chapter_id path int true "Chapter ID"
// @Param        request body AddChapterCommentRequest true "Register comment request"
// @Success      200 {object} store.ChapterComment
// @Failure      400 {object} HTTPError "Error: Invalid Request"
// @Failure      401 {object} HTTPError "Error: Unauthorized"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments  [post]
func (ch *ChapterCommentHandler) HandleAddComment(ctx *gin.Context) {
	chapterID, err := utils.ReadChapterIDParam(ctx)
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

// HandleUpdateComment godoc
// @Summary      Update a comment
// @Description  Updates a book's chapter comment. Expects a body with the edited comment. Returns the updated comment on success.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chapter_id path int true "Chapter ID"
// @Param        id path int true "Comment ID"
// @Param        request body AddChapterCommentRequest true "Edit comment request"
// @Success      200 {object} store.ChapterComment
// @Failure      401 {object} HTTPError "Error: Unauthorized"
// @Failure      404 {object} HTTPError "Error: User not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments/{id} [put]
func (ch *ChapterCommentHandler) HandleUpdateComment(ctx *gin.Context) {
	chapterID, err := utils.ReadChapterIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	fmt.Println(chapterID)
	// TODO: Get comment by ID
	//TODO: only aLlow a comment's edit to the owner of the comment
}

// HandleGetCommentById godoc
// @Summary      Get a comment by id
// @Description  Retrieves the details of a specific comment .
//
//	Provide a valid id as a path parameter. Returns the comment object on success.
//
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        chapter_id path int true "Chapter ID"
// @Param        id path int true "Comment ID"
// @Success      200 {object} store.ChapterComment
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      404 {object} HTTPError "Error: Comment not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments/{id} [get]
func (ch *ChapterCommentHandler) HandleGetCommentById(ctx *gin.Context) {
	// chapterID, err := utils.ReadChapterIDParam(ctx)
	// if err != nil {
	// 	ch.logger.Printf("ERROR: readChapterIDParam %v", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
	// 	return
	// }

	//TODO check if chapter exists
	id, err := utils.ReadIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	fmt.Println("id", id)

	comment, err := ch.chapterCommentStore.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		ch.logger.Printf("ERROR: getUsername %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, comment)

}
