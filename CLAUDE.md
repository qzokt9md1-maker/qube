# Qube - SNS Platform

## Overview
X + Threads合体版の新SNS。完全時系列タイムライン、未読管理が差別化。

## Tech Stack
- Backend: Go (Chi + pgx + go-redis + gorilla/websocket)
- Web: Next.js 16 (TypeScript + Tailwind CSS)
- Mobile: Flutter (Riverpod)
- DB: PostgreSQL (Neon, Singapore)
- Cache: Redis (Upstash, Tokyo)
- API: GraphQL (operationName routing)

## Production URLs
- Web: https://qube-kr6l.vercel.app
- API: https://qube-etp2.onrender.com
- GitHub: https://github.com/qzokt9md1-maker/qube

## Infrastructure
- Vercel (Web, Hobby Free)
- Render (Backend, Free tier - sleeps after 15min idle)
- Neon (PostgreSQL, Free tier, Singapore)
- Upstash (Redis, Free tier, Tokyo, TLS required)

## Design
- Color: モノクロ基調 + シャンパンゴールド #c9a96e
- 文章が主役。色控えめ。アクションボタンは押した時だけ色が変わる
- 紫NG。派手なグラデーションNG

## Monorepo Structure
```
qube/
├── backend/        # Go API server
│   ├── cmd/server/ # main.go (DI wiring)
│   ├── internal/   # config, db, handler, middleware, model, repository, service, ws
│   └── migrations/ # PostgreSQL DDL
├── mobile/         # Flutter app
├── web/            # Next.js app
├── graph/          # GraphQL schema
└── docker-compose.yml  # Local dev (PostgreSQL + Redis + Meilisearch)
```

## Key Commands
```bash
# Local dev
docker compose up -d
cd backend && go run ./cmd/server
cd web && npm run dev

# Build
cd backend && go build ./cmd/server
cd web && npm run build
```

## Environment Variables (Backend)
DATABASE_URL, REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_TLS, JWT_SECRET, PORT, ENV, CORS_ORIGINS, BASE_URL

## Status (2026-03-30)
- 本番デプロイ完了・稼働中
- 全機能動作確認済み: 認証, 投稿, フォロー, いいね, リポスト, DM, 通知, 画像アップロード, WebSocket
- Flutter mobileはテーマ更新済み、App Store未リリース
