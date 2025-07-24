#!/bin/bash

set -e

echo "🚀 OpenAPI モデル生成開始..."

# app-service のモデル生成
echo "📦 app-service のモデルを生成中..."
oapi-codegen -config ../config/oapi-codegen-app-service.yaml ../../../doc/api/app-service/api.yaml

echo "✅ app-service モデル生成完了"


echo "🎉 全モデル生成完了！" 