package systems_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	systems_service "sample-micro-service-api/apps/backend/app-service/internal/service/systems"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

type Handler struct {
	systemsService systems_service.ServiceInterface
}

func NewHandler(systemsService systems_service.ServiceInterface) *Handler {
	return &Handler{
		systemsService: systemsService,
	}
}

// GetSystems - システム一覧取得
func (h *Handler) GetSystems(c *gin.Context) {
	systems, err := h.systemsService.GetSystems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Detail: stringPtr("Failed to retrieve systems"),
		})
		return
	}

	c.JSON(http.StatusOK, systems)
}

// GetSystemById - システム詳細取得
func (h *Handler) GetSystemById(c *gin.Context) {
	idParam := c.Param("id")
	
	system, err := h.systemsService.GetSystemById(c.Request.Context(), idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Detail: stringPtr("System not found"),
		})
		return
	}

	c.JSON(http.StatusOK, system)
}

// CreateSystem - システム作成
func (h *Handler) CreateSystem(c *gin.Context) {
	var req appservice.CreateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	system, err := h.systemsService.CreateSystem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Detail: stringPtr("Failed to create system"),
		})
		return
	}

	c.JSON(http.StatusCreated, system)
}

// UpdateSystem - システム更新
func (h *Handler) UpdateSystem(c *gin.Context) {
	idParam := c.Param("id")
	
	var req appservice.UpdateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	system, err := h.systemsService.UpdateSystem(c.Request.Context(), idParam, req)
	if err != nil {
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Detail: stringPtr("System not found or failed to update"),
		})
		return
	}

	c.JSON(http.StatusOK, system)
}

// DeleteSystem - システム削除
func (h *Handler) DeleteSystem(c *gin.Context) {
	idParam := c.Param("id")
	
	err := h.systemsService.DeleteSystem(c.Request.Context(), idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Detail: stringPtr("System not found"),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
} 