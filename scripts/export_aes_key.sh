#!/bin/bash

# AES-256 の鍵を生成 (32バイト)
KEY=$(openssl rand -hex 32)
echo "Generated AES-256 key: $KEY"

# Base64URL エンコード関数
base64url_encode() {
  local input=$1
  echo -n "$input" | xxd -r -p | base64 | tr '+/' '-_' | tr -d '='
}

# 鍵を Base64URL エンコード
KEY_BASE64URL=$(base64url_encode $KEY)

# JWK を作成
JWK=$(jq -n --arg kty "oct" --arg k "$KEY_BASE64URL" --arg alg "A256GCM" --argjson key_ops '["encrypt", "decrypt"]' '{
  kty: $kty,
  k: $k,
  alg: $alg,
  ext: true,
  key_ops: $key_ops
}')

echo "JWK:"
echo $JWK | jq