package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chillman2101/gits-catalogue/internal/handler"
	"github.com/chillman2101/gits-catalogue/internal/mocks"
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/query"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthorRouter(svc *mocks.MockAuthorService) *gin.Engine {
	r := gin.New()
	h := handler.NewAuthorHandler(svc)
	r.GET("/authors", h.GetAll)
	r.GET("/authors/:id", h.GetByID)
	r.POST("/authors", h.Create)
	r.PUT("/authors/:id", h.Update)
	r.DELETE("/authors/:id", h.Delete)
	return r
}

func TestAuthorHandler_GetAll(t *testing.T) {
	t.Run("200 success", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		authors := []model.Author{{ID: 1, Name: "Author One", Bio: "Bio"}}
		svc.On("GetAll", query.Params{Page: 1, Limit: 10, Sort: "id", Order: "ASC"}).
			Return(authors, int64(1), nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/authors", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var body map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &body)
		assert.True(t, body["success"].(bool))
		svc.AssertExpectations(t)
	})

	t.Run("500 service error", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("GetAll", query.Params{Page: 1, Limit: 10, Sort: "id", Order: "ASC"}).
			Return([]model.Author{}, int64(0), res_err.ErrAuthorNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/authors", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestAuthorHandler_GetByID(t *testing.T) {
	t.Run("200 found", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		author := &model.Author{ID: 1, Name: "Author One"}
		svc.On("GetByID", uint(1)).Return(author, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/authors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("404 not found", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("GetByID", uint(99)).Return(nil, res_err.ErrAuthorNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/authors/99", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestAuthorHandler_Create(t *testing.T) {
	t.Run("201 created", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Create", &model.Author{Name: "New Author", Bio: "Bio text"}).Return(nil)

		body := `{"name":"New Author","bio":"Bio text"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("400 validation error - name too short", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		body := `{"name":"A","bio":"Bio text"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "validation failed", resp["message"])
		svc.AssertNotCalled(t, "Create")
	})

	t.Run("400 validation error - missing name", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		body := `{"bio":"Bio text"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		svc.AssertNotCalled(t, "Create")
	})

	t.Run("409 conflict", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Create", &model.Author{Name: "Existing Author", Bio: ""}).Return(res_err.ErrAuthorConflict)

		body := `{"name":"Existing Author"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestAuthorHandler_Update(t *testing.T) {
	t.Run("200 updated", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Update", uint(1), &model.Author{Name: "Updated Name", Bio: ""}).Return(nil)

		body := `{"name":"Updated Name"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/authors/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("404 not found", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Update", uint(99), &model.Author{Name: "Updated Name", Bio: ""}).Return(res_err.ErrAuthorNotFound)

		body := `{"name":"Updated Name"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/authors/99", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestAuthorHandler_Delete(t *testing.T) {
	t.Run("200 deleted", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Delete", uint(1)).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/authors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("500 service error", func(t *testing.T) {
		svc := new(mocks.MockAuthorService)
		r := setupAuthorRouter(svc)

		svc.On("Delete", uint(99)).Return(res_err.ErrAuthorNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/authors/99", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}
