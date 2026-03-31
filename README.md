# kyudo-dojo-hub

弓道の稽古記録・射形分析・道場運営を統合するプラットフォーム。射法八節の各段階をスマートフォンの動画から骨格推定（pose estimation）で解析し、師範のお手本と比較できる。

## 技術スタック

- **フロントエンド**: TypeScript / React (Vite)
- **API層**: Go (gRPC + REST gateway)
- **ML/データ処理**: Python (MediaPipe)
- **データベース**: PostgreSQL / Redis
- **インフラ**: GCP Cloud Run + Cloud Storage

## セットアップ

```bash
# 依存インストール
npm install

# 開発サーバー起動
npm run dev

# テスト
make test

# 全チェック（lint + test + build）
make check
```

## 主な機能（MVP）

1. 射形動画アップロードと骨格推定による八節フェーズ自動分割・スコアリング
2. 稽古日誌（的中率、気づき、師範コメント）と段位審査チェックリスト管理
3. 道場向け的場予約カレンダーと会員管理ダッシュボード

## アーキテクチャ

```
React SPA (Vite) --> Go API (gRPC/REST) --> PostgreSQL / Redis
                          |
                          v
                  Python MediaPipe Worker (非同期)
                          |
                          v
                  GCP Cloud Storage (動画)
```

- フロントエンド: React SPA がユーザー操作と可視化を担当
- API層: Go gRPC サーバーが REST gateway 経由でフロントエンドと通信
- ML Worker: MediaPipe による骨格推定を非同期で処理
- ストレージ: 動画は GCP Cloud Storage、メタデータは PostgreSQL に保存

## ライセンス

MIT
