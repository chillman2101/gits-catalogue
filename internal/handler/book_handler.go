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

type BookHandler struct {
	service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{service}
}

// GetAll godoc
// @Summary      List all books
// @Tags         books
// @Security     BearerAuth
// @Produce      json
// @Param        page    query     int     false  "Page number (default: 1)"
// @Param        limit   query     int     false  "Items per page (default: 10, max: 100)"
// @Param        search  query     string  false  "Search by title or ISBN"
// @Param        sort    query     string  false  "Sort field (default: id)"
// @Param        order   query     string  false  "Sort order: ASC or DESC (default: ASC)"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /api/v1/books [get]
func (h *BookHandler) GetAll(c *gin.Context) {
	p := query.Parse(c)
	books, total, err := h.service.GetAll(p)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	totalPages := int(total) / p.Limit
	if int(total)%p.Limit != 0 {
		totalPages++
	}
	response.OKWithMeta(c, "books retrieved successfully", dto.ToBookResponseList(books), &response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary      Get book by ID
// @Tags         books
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /api/v1/books/{id} [get]
func (h *BookHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "book not found")
		return
	}
	response.OK(c, "book retrieved successfully", dto.ToBookResponse(*book))
}

// Create godoc
// @Summary      Create a book
// @Tags         books
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        book  body      dto.CreateBookRequest  true  "Book data"
// @Success      201   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /api/v1/books [post]
func (h *BookHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	book := req.ToModel()
	if err := h.service.Create(&book); err != nil {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: err.Error()})
		return
	}
	response.Created(c, "book created successfully", dto.ToBookResponse(book))
}

// Update godoc
// @Summary      Update a book
// @Tags         books
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Book ID"
// @Param        book  body      dto.UpdateBookRequest  true  "Book data"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /api/v1/books/{id} [put]
func (h *BookHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	book := req.ToModel()
	if err := h.service.Update(uint(id), &book); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "book updated successfully", dto.ToBookResponse(book))
}

// Delete godoc
// @Summary      Delete a book
// @Tags         books
// @Security     BearerAuth
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/books/{id} [delete]
func (h *BookHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}
