package systems_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"sample-micro-service-api/package-go/database"
	"sample-micro-service-api/package-go/database/model"
	"sample-micro-service-api/package-go/logging"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

// ServiceInterface はSystemsServiceのインターフェース
type ServiceInterface interface {
	GetSystems(ctx context.Context) ([]appservice.ModelSystem, error)
	SearchSystems(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error)
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
	logging.Debug("Service: Getting all systems using GORM")
	
	var systems []model.System
	if err := s.dbClient.GormDB.WithContext(ctx).
		Model(&model.System{}).
		Order("created_at DESC").
		Find(&systems).Error; err != nil {
		logging.Error("Service: Failed to retrieve systems from database using GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve systems: %w", err)
	}

	// SystemをResponseモデルに変換
	response := make([]appservice.ModelSystem, 0, len(systems))
	for _, system := range systems {
		response = append(response, s.convertToModelSystem(system))
	}

	logging.Debug("Service: Successfully retrieved systems using GORM", zap.Int("count", len(response)))
	return response, nil
}

// SearchSystems - システム検索
func (s *Service) SearchSystems(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error) {
	logging.Debug("Service: Searching systems using GORM",
		zap.String("systemName", systemName),
		zap.String("email", email),
		zap.String("localGovernmentId", localGovernmentId),
	)
	
	var systems []model.System
	query := s.dbClient.GormDB.WithContext(ctx).Model(&model.System{})
	
	// 動的にWHERE条件を追加
	if systemName != "" {
		query = query.Where("system_name ILIKE ?", "%"+systemName+"%")
	}
	if email != "" {
		query = query.Where(&model.System{MailAddress: email})
	}
	if localGovernmentId != "" {
		query = query.Where(&model.System{LocalGovernmentID: localGovernmentId})
	}
	
	if err := query.Order("created_at DESC").Find(&systems).Error; err != nil {
		logging.Error("Service: Failed to search systems using GORM", 
			zap.Error(err),
			zap.String("systemName", systemName),
			zap.String("email", email),
			zap.String("localGovernmentId", localGovernmentId),
		)
		return nil, fmt.Errorf("failed to search systems: %w", err)
	}

	// SystemをResponseモデルに変換
	response := make([]appservice.ModelSystem, 0, len(systems))
	for _, system := range systems {
		response = append(response, s.convertToModelSystem(system))
	}

	logging.Debug("Service: Successfully searched systems using GORM", zap.Int("count", len(response)))
	return response, nil
}

// GetSystemById - システム詳細取得
func (s *Service) GetSystemById(ctx context.Context, id string) (*appservice.ModelSystem, error) {
	logging.Debug("Service: Getting system by ID using GORM", zap.String("id", id))
	
	var system model.System
	if err := s.dbClient.GormDB.WithContext(ctx).
		Model(&model.System{}).
		Where(&model.System{ID: id}).
		First(&system).Error; err != nil {
		logging.Warn("Service: System not found using GORM", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("system not found: %w", err)
	}

	response := s.convertToModelSystem(system)
	logging.Debug("Service: Successfully retrieved system using GORM", zap.String("id", id))
	return &response, nil
}

// CreateSystem - システム作成
func (s *Service) CreateSystem(ctx context.Context, req appservice.CreateSystemJSONBody) (*appservice.ModelSystem, error) {
	logging.Info("Service: Creating new system using GORM", zap.String("systemName", req.SystemName))
	
	system := model.System{
		SystemName:        req.SystemName,
		LocalGovernmentID: stringFromPtr(req.LocalGovernmentId),
		MailAddress:       string(req.MailAddress),
		Telephone:         stringFromPtr(req.Telephone),
		Remark:            stringFromPtr(req.Remark),
	}

	if err := s.dbClient.GormDB.WithContext(ctx).Create(&system).Error; err != nil {
		logging.Error("Service: Failed to create system using GORM", 
			zap.Error(err),
			zap.String("systemName", req.SystemName),
		)
		return nil, fmt.Errorf("failed to create system: %w", err)
	}

	response := s.convertToModelSystem(system)
	logging.Info("Service: Successfully created system using GORM", 
		zap.String("id", system.ID),
		zap.String("systemName", req.SystemName),
	)
	return &response, nil
}

// UpdateSystem - システム更新
func (s *Service) UpdateSystem(ctx context.Context, id string, req appservice.UpdateSystemJSONBody) (*appservice.ModelSystem, error) {
	logging.Info("Service: Updating system using GORM", 
		zap.String("id", id),
		zap.String("systemName", req.SystemName),
	)
	
	var system model.System
	if err := s.dbClient.GormDB.WithContext(ctx).
		Model(&model.System{}).
		Where(&model.System{ID: id}).
		First(&system).Error; err != nil {
		logging.Warn("Service: System not found for update using GORM", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("system not found: %w", err)
	}

	// システム情報を更新
	updates := model.System{
		SystemName:        req.SystemName,
		LocalGovernmentID: stringFromPtr(req.LocalGovernmentId),
		MailAddress:       string(req.MailAddress),
		Telephone:         stringFromPtr(req.Telephone),
		Remark:            stringFromPtr(req.Remark),
	}

	if err := s.dbClient.GormDB.WithContext(ctx).
		Model(&system).
		Updates(updates).Error; err != nil {
		logging.Error("Service: Failed to update system using GORM", 
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update system: %w", err)
	}

	// 更新後のデータを取得
	if err := s.dbClient.GormDB.WithContext(ctx).
		Model(&model.System{}).
		Where(&model.System{ID: id}).
		First(&system).Error; err != nil {
		logging.Error("Service: Failed to retrieve updated system", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve updated system: %w", err)
	}

	response := s.convertToModelSystem(system)
	logging.Info("Service: Successfully updated system using GORM", zap.String("id", id))
	return &response, nil
}

// DeleteSystem - システム削除
func (s *Service) DeleteSystem(ctx context.Context, id string) error {
	logging.Info("Service: Deleting system using GORM", zap.String("id", id))
	
	result := s.dbClient.GormDB.WithContext(ctx).
		Where(&model.System{ID: id}).
		Delete(&model.System{})
	
	if result.Error != nil {
		logging.Error("Service: Failed to delete system using GORM", 
			zap.String("id", id),
			zap.Error(result.Error),
		)
		return fmt.Errorf("failed to delete system: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		logging.Warn("Service: System not found for deletion", zap.String("id", id))
		return fmt.Errorf("system not found")
	}

	logging.Info("Service: Successfully deleted system using GORM", zap.String("id", id))
	return nil
}

// convertToModelSystem - SystemをAPIレスポンスモデルに変換
func (s *Service) convertToModelSystem(system model.System) appservice.ModelSystem {
	// string型のIDをuuid.UUIDに変換
	systemId, _ := uuid.Parse(system.ID)
	
	return appservice.ModelSystem{
		Id:                systemId,
		SystemName:        system.SystemName,
		LocalGovernmentId: stringToPtr(system.LocalGovernmentID),
		CreatedAt:         system.CreatedAt,
		UpdatedAt:         system.UpdatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         stringToPtr(system.Telephone),
		Remark:            stringToPtr(system.Remark),
	}
}

// ヘルパー関数
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
} 