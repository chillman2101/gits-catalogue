package handler

import (
	"net/http"
	"strconv"

	"github.com/chillman2101/gits-catalogue/internal/dto"
	"github.com/chillman2101/gits-catalogue/internal/query"
	"github.com/chillman2101/gits-catalogue/internal/response"
	"github.com/chillman2101/gits-catalogue/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	service service.AuthorService
}

func NewAuthorHandler(service service.AuthorService) *AuthorHandler {
	return &AuthorHandler{service}
}

// GetAll godoc
// @Summary      List all authors
// @Tags         authors
// @Security     BearerAuth
// @Produce      json
// @Param        page    query     int     false  "Page number (default: 1)"
// @Param        limit   query     int     false  "Items per page (default: 10, max: 100)"
// @Param        search  query     string  false  "Search by name"
// @Param        sort    query     string  false  "Sort field (default: id)"
// @Param        order   query     string  false  "Sort order: ASC or DESC (default: ASC)"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /api/v1/authors [get]
func (h *AuthorHandler) GetAll(c *gin.Context) {
	p := query.Parse(c)
	authors, total, err := h.service.GetAll(p)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	totalPages := int(total) / p.Limit
	if int(total)%p.Limit != 0 {
		totalPages++
	}
	response.OKWithMeta(c, "authors retrieved successfully", dto.ToAuthorResponseList(authors), &response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary      Get author by ID
// @Tags         authors
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Author ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /api/v1/authors/{id} [get]
func (h *AuthorHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	author, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "author not found")
		return
	}
	response.OK(c, "author retrieved successfully", dto.ToAuthorResponse(*author))
}

// Create godoc
// @Summary      Create an author
// @Tags         authors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        author  body      dto.CreateAuthorRequest  true  "Author data"
// @Success      201     {object}  response.Response
// @Failure      400     {object}  response.Response
// @Router       /api/v1/authors [post]
func (h *AuthorHandler) Create(c *gin.Context) {
	var req dto.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	author := req.ToModel()
	if err := h.service.Create(&author); err != nil {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: err.Error()})
		return
	}
	response.Created(c, "author created successfully", dto.ToAuthorResponse(author))
}

// Update godoc
// @Summary      Update an author
// @Tags         authors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      int                      true  "Author ID"
// @Param        author  body      dto.UpdateAuthorRequest  true  "Author data"
// @Success      200     {object}  response.Response
// @Failure      400     {object}  response.Response
// @Failure      404     {object}  response.Response
// @Router       /api/v1/authors/{id} [put]
func (h *AuthorHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	author := req.ToModel()
	if err := h.service.Update(uint(id), &author); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "author updated successfully", dto.ToAuthorResponse(author))
}

// Delete godoc
// @Summary      Delete an author
// @Tags         authors
// @Security     BearerAuth
// @Param        id   path      int  true  "Author ID"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/authors/{id} [delete]
func (h *AuthorHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}
