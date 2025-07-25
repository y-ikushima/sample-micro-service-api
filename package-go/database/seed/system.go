package seed

import (
	"fmt"

	"sample-micro-service-api/package-go/database/model"

	"gorm.io/gorm"
)

// SeedSystems inserts test data for systems table using GORM
func SeedSystems(db *gorm.DB) error {
	systems := []model.System{
		{
			SystemName:        "住民基本台帳システム",
			LocalGovernmentID: "",
			MailAddress:       "juki-admin@chiyoda.tokyo.jp",
			Telephone:         "03-1234-5678",
			Remark:            "住民基本台帳の管理を行うシステム",
		},
		{
			SystemName:        "税務管理システム",
			LocalGovernmentID: "",
			MailAddress:       "zeimu-admin@chiyoda.tokyo.jp",
			Telephone:         "03-1234-5679",
			Remark:            "税務関連業務の管理システム",
		},
		{
			SystemName:        "健康管理システム",
			LocalGovernmentID: "",
			MailAddress:       "kenkou-admin@yokohama.lg.jp",
			Telephone:         "045-1234-5678",
			Remark:            "市民の健康管理を支援するシステム",
		},
		{
			SystemName:        "介護保険システム",
			LocalGovernmentID: "",
			MailAddress:       "kaigo-admin@yokohama.lg.jp",
			Telephone:         "045-1234-5679",
			Remark:            "介護保険業務の管理システム",
		},
		{
			SystemName:        "教育情報システム",
			LocalGovernmentID: "",
			MailAddress:       "kyoiku-admin@nagoya.lg.jp",
			Telephone:         "052-1234-5678",
			Remark:            "教育関連情報の管理システム",
		},
		{
			SystemName:        "共通基盤システム",
			LocalGovernmentID: "",
			MailAddress:       "platform-admin@gov-cloud.go.jp",
			Telephone:         "03-0000-0000",
			Remark:            "自治体共通で使用する基盤システム",
		},
		{
			SystemName:        "災害対応システム",
			LocalGovernmentID: "",
			MailAddress:       "saigai-admin@osaka.lg.jp",
			Telephone:         "06-1234-5678",
			Remark:            "災害時の対応管理システム",
		},
		{
			SystemName:        "図書館管理システム",
			LocalGovernmentID: "",
			MailAddress:       "library-admin@chiyoda.tokyo.jp",
			Telephone:         "",  // 空文字列でNULL値を表現
			Remark:            "図書館の蔵書・貸出管理システム",
		},
	}

	fmt.Println("Seeding systems data...")
	for i, systemData := range systems {
		result := db.Create(&systemData)
		if result.Error != nil {
			return fmt.Errorf("failed to create system %d: %w", i+1, result.Error)
		}
		fmt.Printf("Created system: %s (ID: %s)\n", systemData.SystemName, systemData.ID)
	}

	fmt.Printf("Successfully seeded %d systems\n", len(systems))
	return nil
}
