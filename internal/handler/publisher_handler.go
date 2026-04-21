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

type PublisherHandler struct {
	service service.PublisherService
}

func NewPublisherHandler(service service.PublisherService) *PublisherHandler {
	return &PublisherHandler{service}
}

// GetAll godoc
// @Summary      List all publishers
// @Tags         publishers
// @Security     BearerAuth
// @Produce      json
// @Param        page    query     int     false  "Page number (default: 1)"
// @Param        limit   query     int     false  "Items per page (default: 10, max: 100)"
// @Param        search  query     string  false  "Search by name"
// @Param        sort    query     string  false  "Sort field (default: id)"
// @Param        order   query     string  false  "Sort order: ASC or DESC (default: ASC)"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /api/v1/publishers [get]
func (h *PublisherHandler) GetAll(c *gin.Context) {
	p := query.Parse(c)
	publishers, total, err := h.service.GetAll(p)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	totalPages := int(total) / p.Limit
	if int(total)%p.Limit != 0 {
		totalPages++
	}
	response.OKWithMeta(c, "publishers retrieved successfully", dto.ToPublisherResponseList(publishers), &response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary      Get publisher by ID
// @Tags         publishers
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Publisher ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /api/v1/publishers/{id} [get]
func (h *PublisherHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	publisher, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "publisher not found")
		return
	}
	response.OK(c, "publisher retrieved successfully", dto.ToPublisherResponse(*publisher))
}

// Create godoc
// @Summary      Create a publisher
// @Tags         publishers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        publisher  body      dto.CreatePublisherRequest  true  "Publisher data"
// @Success      201        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Router       /api/v1/publishers [post]
func (h *PublisherHandler) Create(c *gin.Context) {
	var req dto.CreatePublisherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	publisher := req.ToModel()
	if err := h.service.Create(&publisher); err != nil {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: err.Error()})
		return
	}
	response.Created(c, "publisher created successfully", dto.ToPublisherResponse(publisher))
}

// Update godoc
// @Summary      Update a publisher
// @Tags         publishers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id         path      int                         true  "Publisher ID"
// @Param        publisher  body      dto.UpdatePublisherRequest  true  "Publisher data"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Router       /api/v1/publishers/{id} [put]
func (h *PublisherHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdatePublisherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	publisher := req.ToModel()
	if err := h.service.Update(uint(id), &publisher); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "publisher updated successfully", dto.ToPublisherResponse(publisher))
}

// Delete godoc
// @Summary      Delete a publisher
// @Tags         publishers
// @Security     BearerAuth
// @Param        id   path      int  true  "Publisher ID"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/publishers/{id} [delete]
func (h *PublisherHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}
