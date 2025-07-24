package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"sample-micro-service-api/package-go/database"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

type SystemsHandler struct {
	dbClient *database.Client
}

func NewSystemsHandler(dbClient *database.Client) *SystemsHandler {
	return &SystemsHandler{
		dbClient: dbClient,
	}
}

// GetSystems - システム一覧取得
func (h *SystemsHandler) GetSystems(c *gin.Context) {
	systems, err := h.dbClient.Queries.GetSystems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Detail: stringPtr("Failed to retrieve systems"),
		})
		return
	}

	// DBモデルをResponseモデルに変換
	var response []appservice.ModelSystem
	for _, system := range systems {
		response = append(response, appservice.ModelSystem{
			Id:                system.ID,
			SystemName:        system.SystemName,
			LocalGovernmentId: nullStringToPtr(system.LocalGovernmentId),
			CreatedAt:         system.CreatedAt,
			UpdatedAt:         system.UpdatedAt,
			MailAddress:       types.Email(system.MailAddress),
			Telephone:         nullStringToPtr(system.Telephone),
			Remark:            nullStringToPtr(system.Remark),
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetSystemById - システム詳細取得
func (h *SystemsHandler) GetSystemById(c *gin.Context) {
	idParam := c.Param("id")
	systemId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid system ID format"),
		})
		return
	}

	system, err := h.dbClient.Queries.GetSystem(c.Request.Context(), systemId)
	if err != nil {
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Detail: stringPtr("System not found"),
		})
		return
	}

	response := appservice.ModelSystem{
		Id:                system.ID,
		SystemName:        system.SystemName,
		LocalGovernmentId: nullStringToPtr(system.LocalGovernmentId),
		CreatedAt:         system.CreatedAt,
		UpdatedAt:         system.UpdatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         nullStringToPtr(system.Telephone),
		Remark:            nullStringToPtr(system.Remark),
	}

	c.JSON(http.StatusOK, response)
}

// CreateSystem - システム作成
func (h *SystemsHandler) CreateSystem(c *gin.Context) {
	var req appservice.CreateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	// DB用のパラメータを準備
	params := database.CreateSystemParams{
		SystemName:        req.SystemName,
		LocalGovernmentId: ptrToNullString(req.LocalGovernmentId),
		MailAddress:       string(req.MailAddress),
		Telephone:         ptrToNullString(req.Telephone),
		Remark:            ptrToNullString(req.Remark),
	}

	system, err := h.dbClient.Queries.CreateSystem(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, appservice.CommonError{
			Status: http.StatusInternalServerError,
			Detail: stringPtr("Failed to create system"),
		})
		return
	}

	response := appservice.ModelSystem{
		Id:                system.ID,
		SystemName:        system.SystemName,
		LocalGovernmentId: nullStringToPtr(system.LocalGovernmentId),
		CreatedAt:         system.CreatedAt,
		UpdatedAt:         system.UpdatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         nullStringToPtr(system.Telephone),
		Remark:            nullStringToPtr(system.Remark),
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateSystem - システム更新
func (h *SystemsHandler) UpdateSystem(c *gin.Context) {
	idParam := c.Param("id")
	systemId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid system ID format"),
		})
		return
	}

	var req appservice.UpdateSystemJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid request body"),
		})
		return
	}

	// DB用のパラメータを準備
	params := database.UpdateSystemParams{
		ID:                systemId,
		SystemName:        req.SystemName,
		LocalGovernmentId: ptrToNullString(req.LocalGovernmentId),
		MailAddress:       string(req.MailAddress),
		Telephone:         ptrToNullString(req.Telephone),
		Remark:            ptrToNullString(req.Remark),
	}

	system, err := h.dbClient.Queries.UpdateSystem(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusNotFound, appservice.CommonError{
			Status: http.StatusNotFound,
			Detail: stringPtr("System not found or failed to update"),
		})
		return
	}

	response := appservice.ModelSystem{
		Id:                system.ID,
		SystemName:        system.SystemName,
		LocalGovernmentId: nullStringToPtr(system.LocalGovernmentId),
		CreatedAt:         system.CreatedAt,
		UpdatedAt:         system.UpdatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         nullStringToPtr(system.Telephone),
		Remark:            nullStringToPtr(system.Remark),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSystem - システム削除
func (h *SystemsHandler) DeleteSystem(c *gin.Context) {
	idParam := c.Param("id")
	systemId, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, appservice.CommonError{
			Status: http.StatusBadRequest,
			Detail: stringPtr("Invalid system ID format"),
		})
		return
	}

	err = h.dbClient.Queries.DeleteSystem(c.Request.Context(), systemId)
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

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func ptrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
} 