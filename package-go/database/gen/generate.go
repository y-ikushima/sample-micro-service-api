package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {


	databaseURL := "postgres://postgres:password@localhost:5432/sample_micro_service"
	if databaseURL == "" {
		log.Fatal("POSTGRES_URL environment variable is required")
	}

	// Initialize GORM connection
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize generator
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../query",               // クエリ用の別ディレクトリ
		OutFile:           "query.go",               // クエリファイル名  
		ModelPkgPath:      "../model", // モデルパッケージパス
		WithUnitTest:      false,                    // ユニットテスト生成するか
		FieldNullable:     true,                     // nullable フィールドの生成
		FieldCoverable:    true,                     // カバレッジのためのフィールド生成
		FieldSignable:     false,                    // 符号付きフィールドの生成
		FieldWithIndexTag: false,                    // インデックスタグの生成
		FieldWithTypeTag:  true,                     // 型タグの生成
	})

	// Use the existing database connection
	g.UseDB(db)

	// Generate all tables with basic query methods
	g.ApplyBasic(g.GenerateAllTable()...)

	// Generate specific models with custom query methods if needed
	// system := g.GenerateModel("system")
	// project := g.GenerateModel("project")
	// user := g.GenerateModel("gcasuser")
	
	// Apply additional query methods
	// g.ApplyBasic(system, project, user)

	// Execute the generator
	g.Execute()

	fmt.Println("✅ GORM models and queries generated successfully!")
} 