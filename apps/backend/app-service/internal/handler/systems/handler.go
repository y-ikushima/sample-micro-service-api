package systems_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	systems_service "sample-micro-service-api/apps/backend/app-service/internal/service/systems"
	"sample-micro-service-api/package-go/logging"
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
	// クエリパラメータを取得
	systemName := c.Query("systemName")
	email := c.Query("email")
	localGovernmentId := c.Query("localGovernmentId")

	var systems []appservice.ModelSystem
	var err error

	// 検索パラメータが指定されている場合は検索を実行、そうでなければ全件取得
	if systemName != "" || email != "" || localGovernmentId != "" {
		logging.Info("Searching systems with parameters",
			zap.String("systemName", systemName),
			zap.String("email", email),
			zap.String("localGovernmentId", localGovernmentId),
		)
		systems, err = h.systemsService.SearchSystems(c.Request.Context(), systemName, email, localGovernmentId)
	} else {
		logging.Debug("Getting all systems")
		systems, err = h.systemsService.GetSystems(c.Request.Context())
	}

	if err != nil {
		logging.Error("Failed to retrieve systems",
			zap.Error(err),
			zap.String("systemName", systemName),
			zap.String("email", email),
			zap.String("localGovernmentId", localGovernmentId),
		)
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: stringPtr("Failed to retrieve systems"),
		})
		return
	}

	logging.Info("Successfully retrieved systems", zap.Int("count", len(systems)))
	c.JSON(http.StatusOK, systems)
}

// GetSystemById - システム詳細取得
func (h *Handler) GetSystemById(c *gin.Context) {
	idParam := c.Param("id")
	
	logging.Debug("Getting system by ID", zap.String("id", idParam))
	
	system, err := h.systemsService.GetSystemById(c.Request.Context(), idParam)
	if err != nil {
		logging.Warn("System not found",
			zap.String("id", idParam),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Title:  "Not Found",
			Detail: stringPtr("System not found"),
		})
		return
	}

	logging.Info("Successfully retrieved system", zap.String("id", idParam))
	c.JSON(http.StatusOK, system)
}

// CreateSystem - システム作成
func (h *Handler) CreateSystem(c *gin.Context) {
	var req appservice.CreateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Warn("Invalid request body for system creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	logging.Info("Creating new system", zap.String("systemName", req.SystemName))

	system, err := h.systemsService.CreateSystem(c.Request.Context(), req)
	if err != nil {
		logging.Error("Failed to create system",
			zap.Error(err),
			zap.String("systemName", req.SystemName),
		)
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: stringPtr("Failed to create system"),
		})
		return
	}

	logging.Info("Successfully created system", 
		zap.String("id", system.Id.String()),
		zap.String("systemName", req.SystemName),
	)
	c.JSON(http.StatusCreated, system)
}

// UpdateSystem - システム更新
func (h *Handler) UpdateSystem(c *gin.Context) {
	idParam := c.Param("id")
	
	var req appservice.UpdateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Warn("Invalid request body for system update", 
			zap.String("id", idParam),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	logging.Info("Updating system", 
		zap.String("id", idParam),
		zap.String("systemName", req.SystemName),
	)

	system, err := h.systemsService.UpdateSystem(c.Request.Context(), idParam, req)
	if err != nil {
		logging.Error("Failed to update system",
			zap.String("id", idParam),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Title:  "Not Found",
			Detail: stringPtr("System not found or failed to update"),
		})
		return
	}

	logging.Info("Successfully updated system", zap.String("id", idParam))
	c.JSON(http.StatusOK, system)
}

// DeleteSystem - システム削除
func (h *Handler) DeleteSystem(c *gin.Context) {
	idParam := c.Param("id")
	
	logging.Info("Deleting system", zap.String("id", idParam))
	
	err := h.systemsService.DeleteSystem(c.Request.Context(), idParam)
	if err != nil {
		logging.Error("Failed to delete system",
			zap.String("id", idParam),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Title:  "Not Found",
			Detail: stringPtr("System not found"),
		})
		return
	}

	logging.Info("Successfully deleted system", zap.String("id", idParam))
	c.Status(http.StatusNoContent)
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
} 