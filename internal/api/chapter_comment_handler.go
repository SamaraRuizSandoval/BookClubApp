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

type ChapterCommentHandler struct {
	chapterCommentStore store.ChapterCommentStore
	ChapterStore        store.ChapterStore
	logger              *log.Logger
}

func NewChapterCommentHandler(commentStore store.ChapterCommentStore, chapterStore store.ChapterStore, logger *log.Logger) *ChapterCommentHandler {
	return &ChapterCommentHandler{
		chapterCommentStore: commentStore,
		ChapterStore:        chapterStore,
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

	_, err = ch.ChapterStore.GetChapterByID(chapterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chapter not found"})
			return
		}
		ch.logger.Printf("Error retrieving chapter by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

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
// @Summary      Update a comment to a book's chapter
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
// @Failure      404 {object} HTTPError "Error: Comment not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments/{id} [put]
func (ch *ChapterCommentHandler) HandleUpdateComment(ctx *gin.Context) {
	chapterID, err := utils.ReadChapterIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	id, err := utils.ReadIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	existingComment, err := ch.chapterCommentStore.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		ch.logger.Printf("ERROR: GetCommentByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if existingComment.ChapterID != chapterID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment does not belong to the specified chapter"})
		return
	}

	var req AddChapterCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ch.logger.Printf("ERROR: decodingChapterComment %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	comment := store.ChapterComment{
		Body: req.Body,
		ID:   id,
	}

	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	userID := int64(user.ID)
	if existingComment.UserID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to edit this comment"})
		return
	}

	err = ch.chapterCommentStore.UpdateComment(&comment)
	if err != nil {
		ch.logger.Printf("ERROR: updateComment %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	updatedComment, err := ch.chapterCommentStore.GetCommentByID(id)
	if err != nil {
		ch.logger.Printf("ERROR: GetCommentByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, updatedComment)
}

// HandleGetCommentById godoc
// @Summary      Get a comment to a book's chapter by id
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
	chapterID, err := utils.ReadChapterIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	_, err = ch.ChapterStore.GetChapterByID(chapterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chapter not found"})
			return
		}
		ch.logger.Printf("Error retrieving chapter by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	id, err := utils.ReadIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	comment, err := ch.chapterCommentStore.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		ch.logger.Printf("ERROR: GetCommentByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if comment.ChapterID != chapterID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment does not belong to the specified chapter"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// HandleDeleteCommentById godoc
// @Summary      Delete a comment to a book's chapter by id
// @Description  Deletes a book's chapter comment. Returns the deleted comment on success.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        chapter_id path int true "Chapter ID"
// @Param        id path int true "Comment ID"
// @Success      200
// @Failure      401 {object} HTTPError "Error: Unauthorized"
// @Failure      404 {object} HTTPError "Error: Comment not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments/{id} [delete]
func (ch *ChapterCommentHandler) HandleDeleteCommentById(ctx *gin.Context) {
	chapterID, err := utils.ReadChapterIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	id, err := utils.ReadIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	existingComment, err := ch.chapterCommentStore.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		ch.logger.Printf("ERROR: GetCommentByID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if existingComment.ChapterID != chapterID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment does not belong to the specified chapter"})
		return
	}

	userValue, _ := ctx.Get("user")
	user := userValue.(*store.User)

	userID := int64(user.ID)
	if existingComment.UserID != userID { // TODO: Add admin check
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to edit this comment"})
		return
	}

	err = ch.chapterCommentStore.DeleteCommentByID(id)
	if err != nil {
		ch.logger.Printf("ERROR: updateComment %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

type PaginatedCommentsResponse struct {
	Items      []*store.ChapterComment `json:"items"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalItems int                     `json:"total_items"`
	TotalPages int                     `json:"total_pages"`
}

// HandleGetCommentsByChapterID godoc
// @Summary      Get the comments of a book's chapter
// @Description  Retrieves the comments of a specific chapter.
//
//	Provide a valid chapter_id as a path and  parameter. Returns the paginated comments object on success.
//
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        chapter_id path int true "Chapter ID"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(20)
// @Success      200 {object} PaginatedCommentsResponse
// @Failure      400 {object} HTTPError "Error: Invalid or missing id"
// @Failure      404 {object} HTTPError "Error: Comment not found"
// @Failure      500 {object} HTTPError "Error: Internal server error"
// @Router       /chapters/{chapter_id}/comments [get]
func (ch *ChapterCommentHandler) HandleGetCommentsByChapterID(ctx *gin.Context) {
	chapterID, err := utils.ReadChapterIDParam(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readChapterIDParam %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	_, err = ch.ChapterStore.GetChapterByID(chapterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chapter not found"})
			return
		}
		ch.logger.Printf("Error retrieving chapter by ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	page, limit, err := utils.ReadPaginationParams(ctx)
	if err != nil {
		ch.logger.Printf("ERROR: readPaginationParams %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	comments, totalItems, err := ch.chapterCommentStore.GetCommentsByChapterID(chapterID, page, limit)
	if err != nil {
		ch.logger.Printf("ERROR: GetCommentsByChapterID %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	totalPages := (totalItems + limit - 1) / limit

	ctx.JSON(http.StatusOK, PaginatedCommentsResponse{
		Items:      comments,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}
