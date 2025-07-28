package seed

import (
	"fmt"

	"sample-micro-service-api/package-go/database/model"

	"gorm.io/gorm"
)

// stringのポインタを作成するヘルパー関数
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// SeedSystems inserts test data for systems table using GORM
func SeedSystems(db *gorm.DB) error {
	systems := []model.System{
		{
			SystemName:        "住民基本台帳システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "juki-admin@chiyoda.tokyo.jp",
			Telephone:         stringPtr("03-1234-5678"),
			Remark:            stringPtr("住民基本台帳の管理を行うシステム"),
		},
		{
			SystemName:        "税務管理システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "zeimu-admin@chiyoda.tokyo.jp",
			Telephone:         stringPtr("03-1234-5679"),
			Remark:            stringPtr("税務関連業務の管理システム"),
		},
		{
			SystemName:        "健康管理システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "kenkou-admin@yokohama.lg.jp",
			Telephone:         stringPtr("045-1234-5678"),
			Remark:            stringPtr("市民の健康管理を支援するシステム"),
		},
		{
			SystemName:        "介護保険システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "kaigo-admin@yokohama.lg.jp",
			Telephone:         stringPtr("045-1234-5679"),
			Remark:            stringPtr("介護保険業務の管理システム"),
		},
		{
			SystemName:        "教育情報システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "kyoiku-admin@nagoya.lg.jp",
			Telephone:         stringPtr("052-1234-5678"),
			Remark:            stringPtr("教育関連情報の管理システム"),
		},
		{
			SystemName:        "共通基盤システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "platform-admin@gov-cloud.go.jp",
			Telephone:         stringPtr("03-0000-0000"),
			Remark:            stringPtr("自治体共通で使用する基盤システム"),
		},
		{
			SystemName:        "災害対応システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "saigai-admin@osaka.lg.jp",
			Telephone:         stringPtr("06-1234-5678"),
			Remark:            stringPtr("災害時の対応管理システム"),
		},
		{
			SystemName:        "図書館管理システム",
			LocalGovernmentID: stringPtr(""),
			MailAddress:       "library-admin@chiyoda.tokyo.jp",
			Telephone:         nil,  // NULL値を表現
			Remark:            stringPtr("図書館の蔵書・貸出管理システム"),
		},
	}

	fmt.Println("Seeding systems data...")
	for i, systemData := range systems {
		result := db.Create(&systemData)
		if result.Error != nil {
			return fmt.Errorf("failed to create system %d: %w", i+1, result.Error)
		}
		systemID := ""
		if systemData.ID != nil {
			systemID = *systemData.ID
		}
		fmt.Printf("Created system: %s (ID: %s)\n", systemData.SystemName, systemID)
	}

	fmt.Printf("Successfully seeded %d systems\n", len(systems))
	return nil
}
