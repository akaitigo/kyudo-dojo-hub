# kyudo-dojo-hub

弓道の稽古記録・射形分析・道場運営を統合するプラットフォーム。

## コマンド
- ビルド: `make build`
- テスト: `make test`
- lint: `make lint`
- フォーマット: `make format`
- 全チェック: `make check`
- 品質ゲート: `make quality`

## ワークフロー
1. research.md を作成（調査結果の記録）
2. plan.md を作成（実装計画。人間承認まで実装禁止）
3. 承認後に実装開始。plan.md のtodoを進捗管理に使用

## 技術スタック
- TypeScript/React (Vite) — フロントエンド
- Go (gRPC + REST gateway) — API層
- Python (MediaPipe) — 骨格推定
- PostgreSQL, Redis, GCP Cloud Run

## ルール
- TypeScript: ~/.claude/rules/typescript.md 参照
- ADR: docs/adr/ 参照。新規決定はADRを書いてから実装
- テスト: 機能追加時は必ずテストを同時に書く
- lint設定の変更禁止（ADR必須）

## 構造
src/ — React SPA (Vite)
test/e2e/ — Playwright E2Eテスト
docs/ — ADR・品質チェックリスト

## 禁止事項
- any型(TS) → unknown + 型ガード
- console.log / TODO コメントのコミット
- .env・credentials のコミット
- lint設定の無効化

## 状態管理
- git log + GitHub Issues でセッション間の状態を管理
- セッション開始: `bash .claude/startup.sh`
