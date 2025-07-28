package systems_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"sample-micro-service-api/package-go/database"
	"sample-micro-service-api/package-go/database/model"
	"sample-micro-service-api/package-go/database/query"
	"sample-micro-service-api/package-go/logging"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

// SystemWithLocalGovernment はJOIN結果を格納する構造体
type SystemWithLocalGovernment struct {
	model.System
	LocalGovernmentName *string `json:"local_government_name"`
}

// SystemUpdateFields は更新用の型安全な構造体
type SystemUpdateFields struct {
	SystemName        string  `gorm:"column:systemName"`
	LocalGovernmentID *string `gorm:"column:localGovernmentId"`
	MailAddress       string  `gorm:"column:mailAddress"`
	Telephone         *string `gorm:"column:telephone"`
	Remark            *string `gorm:"column:remark"`
}

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

// GetSystems - システム一覧取得（型安全）
func (s *Service) GetSystems(ctx context.Context) ([]appservice.ModelSystem, error) {
	logging.Debug("Service: Getting all systems using type-safe GORM Gen")
	
	q := query.Use(s.dbClient.GormDB)
	
	// 型安全なフィールド指定でソート
	systems, err := q.System.WithContext(ctx).
		Select(q.System.ALL).
		Order(q.System.CreatedAt.Desc()).
		Find()
	
	if err != nil {
		logging.Error("Service: Failed to retrieve systems from database using type-safe GORM Gen", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve systems: %w", err)
	}

	// 型安全な変換
	response := make([]appservice.ModelSystem, 0, len(systems))
	for _, system := range systems {
		response = append(response, s.convertToModelSystem(*system))
	}

	logging.Debug("Service: Successfully retrieved systems using type-safe GORM Gen", zap.Int("count", len(response)))
	return response, nil
}

// SearchSystems - gorm-genの型安全性の限界と対策を示すシステム検索
func (s *Service) SearchSystems(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error) {
	logging.Debug("Service: Searching systems with type safety considerations",
		zap.String("systemName", systemName),
		zap.String("email", email),
		zap.String("localGovernmentId", localGovernmentId),
	)
	
	q := query.Use(s.dbClient.GormDB)
	
	// 🚨 GORM-Genの型安全性の問題: これはコンパイルは通るが実行時エラー
	// q.MLocalGovernment.NonExistentField.As("non_existent") // コンパイルOK、実行時NG
	
	var results []SystemWithLocalGovernment
	
	// 💡 解決策1: モデル構造体を直接参照して安全性を向上
	// 実際に存在するフィールドのみを使用
	localGovModel := model.MLocalGovernment{} // コンパイル時に構造体が検証される
	_ = localGovModel.PrefectureName          // フィールドの存在を暗黙的に検証
	
	queryBuilder := q.System.WithContext(ctx).
		Select(
			q.System.ALL,
			// ✅ 構造体で検証済みのフィールドを使用
			q.MLocalGovernment.PrefectureName.As("local_government_name"),
		).
		LeftJoin(q.MLocalGovernment, q.MLocalGovernment.ID.EqCol(q.System.LocalGovernmentID))
	
	// 条件指定
	if systemName != "" {
		queryBuilder = queryBuilder.Where(q.System.SystemName.Like("%" + systemName + "%"))
	}
	if email != "" {
		queryBuilder = queryBuilder.Where(q.System.MailAddress.Eq(email))
	}
	if localGovernmentId != "" {
		queryBuilder = queryBuilder.Where(q.System.LocalGovernmentID.Eq(localGovernmentId))
	}
	
	err := queryBuilder.Order(q.System.CreatedAt.Desc()).Scan(&results)
	if err != nil {
		logging.Error("Service: Failed to search systems", 
			zap.Error(err),
			zap.String("systemName", systemName),
			zap.String("email", email),
			zap.String("localGovernmentId", localGovernmentId),
		)
		return nil, fmt.Errorf("failed to search systems: %w", err)
	}

	// 結果変換
	response := make([]appservice.ModelSystem, 0, len(results))
	for _, result := range results {
		response = append(response, s.convertToModelSystem(result.System))
	}

	logging.Debug("Service: Successfully searched systems", zap.Int("count", len(response)))
	return response, nil
}

// 💡 GORM-Genの型安全性を向上させる実用的なパターン

// validateRequiredFields - コンパイル時にフィールド存在を検証
func validateRequiredFields() {
	// この関数は実行されないが、コンパイル時にフィールドの存在をチェック
	systemModel := model.System{}
	localGovModel := model.MLocalGovernment{}
	
	// 使用予定のフィールドにアクセスしてコンパイル時検証
	_ = systemModel.ID                    // ✅ 存在確認
	_ = systemModel.SystemName            // ✅ 存在確認
	_ = systemModel.LocalGovernmentID     // ✅ 存在確認
	_ = systemModel.MailAddress           // ✅ 存在確認
	_ = systemModel.CreatedAt            // ✅ 存在確認
	_ = localGovModel.PrefectureName     // ✅ 存在確認
	
	// 🚨 実証: 存在しないフィールドにアクセスするとコンパイルエラー
	// _ = systemModel.NonExistentField     // ❌ これを有効にするとコンパイルエラー！
	// _ = localGovModel.FakeField          // ❌ これを有効にするとコンパイルエラー！
}

// init関数で初期化時に検証実行
func init() {
	validateRequiredFields()
}

// GetSystemById - 型安全なシステム詳細取得
func (s *Service) GetSystemById(ctx context.Context, id string) (*appservice.ModelSystem, error) {
	logging.Debug("Service: Getting system by ID using type-safe GORM Gen", zap.String("id", id))
	
	q := query.Use(s.dbClient.GormDB)
	
	// 型安全なID検索
	system, err := q.System.WithContext(ctx).
		Select(q.System.ALL).
		Where(q.System.ID.Eq(id)).
		First()
	
	if err != nil {
		logging.Warn("Service: System not found using type-safe GORM Gen", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("system not found: %w", err)
	}

	response := s.convertToModelSystem(*system)
	logging.Debug("Service: Successfully retrieved system using type-safe GORM Gen", zap.String("id", id))
	return &response, nil
}

// CreateSystem - 型安全なシステム作成
func (s *Service) CreateSystem(ctx context.Context, req appservice.CreateSystemJSONBody) (*appservice.ModelSystem, error) {
	logging.Info("Service: Creating new system using type-safe GORM Gen", zap.String("systemName", req.SystemName))
	
	// 型安全なモデル作成
	system := &model.System{
		SystemName:        req.SystemName,
		LocalGovernmentID: req.LocalGovernmentId,
		MailAddress:       string(req.MailAddress),
		Telephone:         req.Telephone,
		Remark:            req.Remark,
	}

	q := query.Use(s.dbClient.GormDB)
	
	// 型安全な作成
	err := q.System.WithContext(ctx).Create(system)
	if err != nil {
		logging.Error("Service: Failed to create system using type-safe GORM Gen", 
			zap.Error(err),
			zap.String("systemName", req.SystemName),
		)
		return nil, fmt.Errorf("failed to create system: %w", err)
	}

	response := s.convertToModelSystem(*system)
	systemID := ""
	if system.ID != nil {
		systemID = *system.ID
	}
	logging.Info("Service: Successfully created system using type-safe GORM Gen", 
		zap.String("id", systemID),
		zap.String("systemName", req.SystemName),
	)
	return &response, nil
}

// UpdateSystem - 型安全なシステム更新
func (s *Service) UpdateSystem(ctx context.Context, id string, req appservice.UpdateSystemJSONBody) (*appservice.ModelSystem, error) {
	logging.Info("Service: Updating system using type-safe GORM Gen", 
		zap.String("id", id),
		zap.String("systemName", req.SystemName),
	)
	
	q := query.Use(s.dbClient.GormDB)
	
	// 型安全な存在確認
	_, err := q.System.WithContext(ctx).
		Select(q.System.ID).
		Where(q.System.ID.Eq(id)).
		First()
	if err != nil {
		logging.Warn("Service: System not found for update using type-safe GORM Gen", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("system not found: %w", err)
	}

	// 型安全な更新フィールド構築
	updateFields := SystemUpdateFields{
		SystemName:        req.SystemName,
		LocalGovernmentID: req.LocalGovernmentId,
		MailAddress:       string(req.MailAddress),
		Telephone:         req.Telephone,
		Remark:            req.Remark,
	}

	// 型安全な更新（構造体ベース）
	_, err = q.System.WithContext(ctx).
		Select(
			q.System.SystemName,
			q.System.LocalGovernmentID,
			q.System.MailAddress,
			q.System.Telephone,
			q.System.Remark,
		).
		Where(q.System.ID.Eq(id)).
		Updates(updateFields)
	
	if err != nil {
		logging.Error("Service: Failed to update system using type-safe GORM Gen", 
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update system: %w", err)
	}

	// 型安全な更新後データ取得
	updatedSystem, err := q.System.WithContext(ctx).
		Select(q.System.ALL).
		Where(q.System.ID.Eq(id)).
		First()
	if err != nil {
		logging.Error("Service: Failed to retrieve updated system", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve updated system: %w", err)
	}

	response := s.convertToModelSystem(*updatedSystem)
	logging.Info("Service: Successfully updated system using type-safe GORM Gen", zap.String("id", id))
	return &response, nil
}

// DeleteSystem - 型安全なシステム削除
func (s *Service) DeleteSystem(ctx context.Context, id string) error {
	logging.Info("Service: Deleting system using type-safe GORM Gen", zap.String("id", id))
	
	q := query.Use(s.dbClient.GormDB)
	
	// 型安全な削除
	info, err := q.System.WithContext(ctx).
		Where(q.System.ID.Eq(id)).
		Delete()
	
	if err != nil {
		logging.Error("Service: Failed to delete system using type-safe GORM Gen", 
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete system: %w", err)
	}

	if info.RowsAffected == 0 {
		logging.Warn("Service: System not found for deletion", zap.String("id", id))
		return fmt.Errorf("system not found")
	}

	logging.Info("Service: Successfully deleted system using type-safe GORMGen", zap.String("id", id))
	return nil
}

// convertToModelSystem - 型安全なSystemからAPIレスポンスモデルへの変換
func (s *Service) convertToModelSystem(system model.System) appservice.ModelSystem {
	// 型安全なポインタ値取得
	var systemId uuid.UUID
	if system.ID != nil {
		if parsed, err := uuid.Parse(*system.ID); err == nil {
			systemId = parsed
		}
	}

	var createdAt, updatedAt time.Time
	if system.CreatedAt != nil {
		createdAt = *system.CreatedAt
	}
	if system.UpdatedAt != nil {
		updatedAt = *system.UpdatedAt
	}
	
	return appservice.ModelSystem{
		Id:                systemId,
		SystemName:        system.SystemName,
		LocalGovernmentId: system.LocalGovernmentID,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		MailAddress:       types.Email(system.MailAddress),
		Telephone:         system.Telephone,
		Remark:            system.Remark,
	}
} 