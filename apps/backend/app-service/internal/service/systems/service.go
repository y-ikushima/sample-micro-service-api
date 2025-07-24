package systems_service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"sample-micro-service-api/package-go/database"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

// ServiceInterface はSystemsServiceのインターフェース
type ServiceInterface interface {
	GetSystems(ctx context.Context) ([]appservice.ModelSystem, error)
	GetSystemById(ctx context.Context, id string) (*appservice.ModelSystem, error)
	CreateSystem(ctx context.Context, req appservice.CreateSystemJSONBody) (*appservice.ModelSystem, error)
	UpdateSystem(ctx context.Context, id string, req appservice.UpdateSystemJSONBody) (*appservice.ModelSystem, error)
	DeleteSystem(ctx context.Context, id string) error
}

// Service はシステム関連のビジネスロジックを処理する
type Service struct {
	dbClient *database.Client
}

// NewService はServiceの新しいインスタンスを作成
func NewService(dbClient *database.Client) ServiceInterface {
	return &Service{
		dbClient: dbClient,
	}
}

// GetSystems - システム一覧取得
func (s *Service) GetSystems(ctx context.Context) ([]appservice.ModelSystem, error) {
	systems, err := s.dbClient.Queries.GetSystems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve systems: %w", err)
	}

	// DBモデルをResponseモデルに変換
	var response []appservice.ModelSystem
	for _, system := range systems {
		response = append(response, s.convertToModelSystem(system))
	}

	return response, nil
}

// GetSystemById - システム詳細取得
func (s *Service) GetSystemById(ctx context.Context, id string) (*appservice.ModelSystem, error) {
	systemId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid system ID format: %w", err)
	}

	system, err := s.dbClient.Queries.GetSystem(ctx, systemId)
	if err != nil {
		return nil, fmt.Errorf("system not found: %w", err)
	}

	response := s.convertToModelSystem(system)
	return &response, nil
}

// CreateSystem - システム作成
func (s *Service) CreateSystem(ctx context.Context, req appservice.CreateSystemJSONBody) (*appservice.ModelSystem, error) {
	// DB用のパラメータを準備
	params := database.CreateSystemParams{
		SystemName:        req.SystemName,
		LocalGovernmentId: ptrToNullString(req.LocalGovernmentId),
		MailAddress:       string(req.MailAddress),
		Telephone:         ptrToNullString(req.Telephone),
		Remark:            ptrToNullString(req.Remark),
	}

	system, err := s.dbClient.Queries.CreateSystem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create system: %w", err)
	}

	response := s.convertToModelSystem(system)
	return &response, nil
}

// UpdateSystem - システム更新
func (s *Service) UpdateSystem(ctx context.Context, id string, req appservice.UpdateSystemJSONBody) (*appservice.ModelSystem, error) {
	systemId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid system ID format: %w", err)
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

	system, err := s.dbClient.Queries.UpdateSystem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("system not found or failed to update: %w", err)
	}

	response := s.convertToModelSystem(system)
	return &response, nil
}

// DeleteSystem - システム削除
func (s *Service) DeleteSystem(ctx context.Context, id string) error {
	systemId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid system ID format: %w", err)
	}

	err = s.dbClient.Queries.DeleteSystem(ctx, systemId)
	if err != nil {
		return fmt.Errorf("system not found: %w", err)
	}

	return nil
}

// convertToModelSystem - DBモデルをAPIレスポンスモデルに変換
func (s *Service) convertToModelSystem(system database.System) appservice.ModelSystem {
	return appservice.ModelSystem{
		Id:                system.ID,
		SystemName:        system.SystemName,
		LocalGovernmentId: nullStringToPtr(system.LocalGovernmentId),
		CreatedAt:         system.CreatedAt,
		UpdatedAt:         system.UpdatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         nullStringToPtr(system.Telephone),
		Remark:            nullStringToPtr(system.Remark),
	}
}

// ヘルパー関数
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