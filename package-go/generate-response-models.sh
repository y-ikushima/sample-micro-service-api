#!/bin/bash

set -e

echo "🚀 Response モデル生成開始..."

# app-serviceディレクトリに移動してモデル生成
cd response/app-service
../scripts/generate-models.sh

cd ../..
echo "�� Response モデル生成完了！" 