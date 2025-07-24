package seed

import (
	"context"
	"database/sql"
	"fmt"

	"sample-micro-service-api/package-go/database/internal/db"
)

// SeedSystems inserts test data for systems table
func SeedSystems(database *sql.DB) error {
	queries := db.New(database)
	ctx := context.Background()

	systems := []db.CreateSystemParams{
		{
			SystemName:        "住民基本台帳システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "juki-admin@chiyoda.tokyo.jp",
			Telephone:         sql.NullString{String: "03-1234-5678", Valid: true},
			Remark:           sql.NullString{String: "住民基本台帳の管理を行うシステム", Valid: true},
		},
		{
			SystemName:        "税務管理システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "zeimu-admin@chiyoda.tokyo.jp",
			Telephone:         sql.NullString{String: "03-1234-5679", Valid: true},
			Remark:           sql.NullString{String: "税務関連業務の管理システム", Valid: true},
		},
		{
			SystemName:        "健康管理システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "kenkou-admin@yokohama.lg.jp",
			Telephone:         sql.NullString{String: "045-1234-5678", Valid: true},
			Remark:           sql.NullString{String: "市民の健康管理を支援するシステム", Valid: true},
		},
		{
			SystemName:        "介護保険システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "kaigo-admin@yokohama.lg.jp",
			Telephone:         sql.NullString{String: "045-1234-5679", Valid: true},
			Remark:           sql.NullString{String: "介護保険業務の管理システム", Valid: true},
		},
		{
			SystemName:        "教育情報システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "kyoiku-admin@nagoya.lg.jp",
			Telephone:         sql.NullString{String: "052-1234-5678", Valid: true},
			Remark:           sql.NullString{String: "教育関連情報の管理システム", Valid: true},
		},
		{
			SystemName:        "共通基盤システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "platform-admin@gov-cloud.go.jp",
			Telephone:         sql.NullString{String: "03-0000-0000", Valid: true},
			Remark:           sql.NullString{String: "自治体共通で使用する基盤システム", Valid: true},
		},
		{
			SystemName:        "災害対応システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "saigai-admin@osaka.lg.jp",
			Telephone:         sql.NullString{String: "06-1234-5678", Valid: true},
			Remark:           sql.NullString{String: "災害時の対応管理システム", Valid: true},
		},
		{
			SystemName:        "図書館管理システム",
			LocalGovernmentId: sql.NullString{String: "", Valid: false}, // null
			MailAddress:       "library-admin@chiyoda.tokyo.jp",
			Telephone:         sql.NullString{String: "", Valid: false},
			Remark:           sql.NullString{String: "図書館の蔵書・貸出管理システム", Valid: true},
		},
	}

	fmt.Println("Seeding systems data...")
	for i, systemData := range systems {
		system, err := queries.CreateSystem(ctx, systemData)
		if err != nil {
			return fmt.Errorf("failed to create system %d: %w", i+1, err)
		}
		fmt.Printf("Created system: %s (ID: %s)\n", system.SystemName, system.ID.String())
	}

	fmt.Printf("Successfully seeded %d systems\n", len(systems))
	return nil
}
