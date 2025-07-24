package systems_service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"sample-micro-service-api/package-go/database"
	appservice "sample-micro-service-api/package-go/response/app-service"
)

// ServiceInterface はSystemsServiceのインターフェース
type ServiceInterface interface {
	GetSystems(ctx context.Context) ([]appservice.ModelSystem, error)
	SearchSystems(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error)
	SearchSystemsDynamic(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error) // 新しいメソッド追加
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

// SearchSystems - システム検索
func (s *Service) SearchSystems(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error) {
	params := database.SearchSystemsParams{
		Column1: systemName,      // systemName
		Column2: email,           // email
		Column3: localGovernmentId, // localGovernmentId
	}
	
	systems, err := s.dbClient.Queries.SearchSystems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search systems: %w", err)
	}

	// DBモデルをResponseモデルに変換
	var response []appservice.ModelSystem
	for _, system := range systems {
		response = append(response, s.convertToModelSystem(system))
	}

	return response, nil
}

// SearchSystemsDynamic - システム検索（動的SQL構築版サンプル）
func (s *Service) SearchSystemsDynamic(ctx context.Context, systemName, email, localGovernmentId string) ([]appservice.ModelSystem, error) {
	baseQuery := `
		SELECT id, "systemName", "localGovernmentId", "createdAt", "updatedAt", 
		       "mailAddress", telephone, remark
		FROM public.system
	`
	
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	// IF的な条件分岐でクエリ構築
	if systemName != "" {
		conditions = append(conditions, fmt.Sprintf(`"systemName" ILIKE $%d`, argIndex))
		args = append(args, "%"+systemName+"%")
		argIndex++
	}
	
	if email != "" {
		conditions = append(conditions, fmt.Sprintf(`"mailAddress" = $%d`, argIndex))
		args = append(args, email)
		argIndex++
	}
	
	if localGovernmentId != "" {
		conditions = append(conditions, fmt.Sprintf(`"localGovernmentId" = $%d`, argIndex))
		args = append(args, localGovernmentId)
		argIndex++
	}
	
	// WHERE句の構築
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	baseQuery += ` ORDER BY "createdAt" DESC`
	
	// 実行
	rows, err := s.dbClient.DB.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search systems: %w", err)
	}
	defer rows.Close()
	
	// 結果の処理
	var systems []database.System
	for rows.Next() {
		var system database.System
		err := rows.Scan(
			&system.ID,
			&system.SystemName,
			&system.LocalGovernmentId,
			&system.CreatedAt,
			&system.UpdatedAt,
			&system.MailAddress,
			&system.Telephone,
			&system.Remark,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		systems = append(systems, system)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
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