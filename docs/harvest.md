# Harvest Report — kyudo-dojo-hub

> 生成日: 2026-03-30
> リポジトリ: akaitigo/kyudo-dojo-hub
> アイデア起源: idea #493

---

## 1. プロジェクト概要

弓道の稽古記録・射形分析・道場運営を統合するプラットフォーム。
MVP フェーズでは TypeScript/React フロントエンドに集中し、Go API・Python ML はモック層で代替。

**技術スタック**: TypeScript, React, Vite, Vitest, Playwright, Biome, oxlint

---

## 2. 主要メトリクス

| メトリクス | 値 |
|---|---|
| 非マージコミット数 | 7 |
| Issue 数（全体） | 5 |
| Issue クローズ率 | 5/5 (100%) |
| PR 数（全体） | 10 |
| PR マージ率 | 5/10 (50%) |
| PR 未マージ（dependabot） | 5 (deps更新、意図的に未マージ) |
| テストスイート | 6 ファイル / 60 テスト / 全パス |
| ADR 数 | 1 |
| CLAUDE.md 行数 | 43 (上限50以内) |
| プロジェクト期間 | 2026-03-30（単日完了） |

---

## 3. Issue / PR 一覧

### Issues（全5件 — 全クローズ）

| # | タイトル | ラベル |
|---|---|---|
| 1 | プロジェクト基盤セットアップ（CI/CD・リンター・テストフレームワーク） | model:haiku, mvp |
| 3 | コアデータモデル設計・APIモック層の構築 | model:opus, mvp |
| 5 | 稽古日誌・的中率記録・段位審査チェックリストUI | model:sonnet, mvp |
| 8 | 射形動画アップロード・八節フェーズ分割表示UI | model:sonnet, mvp |
| 10 | 道場管理ダッシュボード（的場予約・会員管理） | model:sonnet, mvp |

### PRs（全10件）

| # | タイトル | 状態 |
|---|---|---|
| 11 | feat: プロジェクト基盤セットアップ | MERGED |
| 12 | feat: コアデータモデル設計・APIモック層の構築 | MERGED |
| 13 | feat: 稽古日誌・的中率記録・段位審査チェックリストUI | MERGED |
| 14 | feat: 射形動画アップロード・八節フェーズ分割表示UI | MERGED |
| 15 | feat: 道場管理ダッシュボード（的場予約・会員管理） | MERGED |
| 2 | build(deps-dev): bump @vitejs/plugin-react | OPEN (dependabot) |
| 4 | build(deps-dev): bump vitest | OPEN (dependabot) |
| 6 | build(deps-dev): bump @biomejs/biome | OPEN (dependabot) |
| 7 | build(deps-dev): bump @vitest/coverage-v8 | OPEN (dependabot) |
| 9 | build(deps-dev): bump oxlint | OPEN (dependabot) |

---

## 4. ハーネス適用状況

### Layer-0: リポジトリ衛生

| 項目 | 状態 | 備考 |
|---|---|---|
| CLAUDE.md | ✅ 43行 | 50行以内。ポインタ中心の設計 |
| .claude/CLAUDE.md（アーキ概要） | ✅ | システム構成・データモデル・外部連携を記載 |
| ADR | ⚠️ 1件 | 001-frontend-first-mvp.md のみ。追加のADR未作成 |
| LICENSE | ✅ | MIT |
| README.md | ✅ | 存在確認済み |
| .gitignore | ✅ | |

### Layer-1: 決定論的ツール強制

| 項目 | 状態 | 備考 |
|---|---|---|
| settings.json (hooks) | ✅ | PreToolUse, PostToolUse, PreCompact, Stop フック完備 |
| PreToolUse: lint設定保護 | ✅ | biome.json, .oxlintrc.json, tsconfig.json, Makefile の編集ブロック |
| PreToolUse: 機密ファイル保護 | ✅ | .env, credentials, *.pem 等の編集ブロック |
| PreToolUse: 破壊的コマンドブロック | ✅ | rm -rf, DROP TABLE, --force, --no-verify |
| PostToolUse: 自動lint | ✅ | post-lint.sh で編集後自動チェック |
| PreCompact: CLAUDE.md バックアップ | ✅ | |
| Stop: 全チェック実行 | ✅ | make check && make quality |
| Stop: E2Eテスト | ✅ | Playwright テスト実行 |
| lefthook.yml | ✅ | pre-commit: lint, format, test, archgate |
| startup.sh | ✅ | 自動ツールインストール・ヘルスチェック |
| Makefile | ✅ | build, test, lint, format, typecheck, check, quality, e2e |
| CI/CD (GitHub Actions) | ✅ | ci.yml + dependabot-auto-merge.yml |
| Dependabot | ✅ | 5件の deps PR が自動生成済み |

### Layer-2: 計画と実行の分離

| 項目 | 状態 | 備考 |
|---|---|---|
| research.md → plan.md ワークフロー | ✅ | CLAUDE.md に明記 |
| Issue → PR の1:1対応 | ✅ | 全5 Issue が対応PRでクローズ |
| model ラベルによる難易度分類 | ✅ | haiku/sonnet/opus で分類 |
| quality ゲート（make quality） | ✅ | LICENSE, TODO, secrets, PRD, CLAUDE.md行数チェック |

### 総合スコア

| レイヤー | 充足率 | 判定 |
|---|---|---|
| Layer-0: リポジトリ衛生 | 5/6 (83%) | ⚠️ ADR不足 |
| Layer-1: ツール強制 | 12/12 (100%) | ✅ 完全適用 |
| Layer-2: 計画分離 | 4/4 (100%) | ✅ 完全適用 |
| **総合** | **21/22 (95%)** | **✅ 高水準** |

---

## 5. 良かった点

1. **単日でMVP完成**: 全5 Issue クローズ、60テスト全パス、v1.0.0リリースまで1日で完了
2. **ハーネス適用率95%**: Layer-1（ツール強制）とLayer-2（計画分離）は100%充足
3. **テスト品質**: 6ファイル60テストでドメインロジック・バリデーション・API モック層を網羅
4. **Issue-PR完全対応**: 全Issueが対応PRでクローズされ、トレーサビリティが高い
5. **Dependabot自動化**: CI/dependabot-auto-merge.yml で依存更新の自動マージパイプライン構築済み

---

## 6. 改善点

1. **ADR不足**: 001のみ。射法八節の判定方式（ルールベース vs ML）など、設計判断がADRとして残っていない
2. **dependabot PR未処理**: 5件のdeps更新PRが未マージ。メジャーバージョンアップのため慎重対応は正しいが、対応方針を決めるべき
3. **E2Eテスト未確認**: playwright.config.ts は存在するが、E2Eテストの実際のカバレッジは未検証
4. **Go API / Python ML のモック境界**: モック→実装への移行計画がドキュメント化されていない

---

## 7. テンプレート改善提案

idea-launch テンプレートおよびハーネス基盤への改善提案:

| # | カテゴリ | 提案 | 根拠 | 優先度 |
|---|---|---|---|---|
| 1 | ADR | idea-launch 時に ADR-001 を自動生成するテンプレートを追加 | 毎回ADR不足が指摘される。初期設計判断（技術選定・MVP範囲）をADRとして自動記録すべき | HIGH |
| 2 | Dependabot | dependabot PRのトリアージ基準をCLAUDE.mdに追記するテンプレート化 | メジャーバージョンアップの対応方針が不明確で放置されやすい | MEDIUM |
| 3 | モック境界 | モック→実装移行のチェックリストをPRD/Issue生成時に自動付与 | マルチサービス構成でモック境界の管理が属人化しやすい | MEDIUM |
| 4 | E2Eテスト | Stop フックのE2E結果をサマリとしてログ出力する仕組み | E2Eが実行されているか・何件パスしたかが不透明 | LOW |
| 5 | Harvest自動化 | ship フェーズ完了後に harvest レポート生成を自動トリガー | 手動実行だと忘れがち。idea-ship の最終ステップに組み込むべき | HIGH |
| 6 | メトリクス基準値 | テスト数/ADR数/CLAUDE.md行数の推奨基準値をharness文書に明記 | 「足りているか」の判断基準が暗黙知になっている | MEDIUM |

---

## 8. 次のアクション

- [ ] ADR-002: 射法八節判定方式（ルールベース vs ML）を作成
- [ ] Dependabot 5件のPRトリアージ（マージ or クローズ判断）
- [ ] E2Eテストカバレッジの確認と拡充
- [ ] idea-launch テンプレートへのADR自動生成機能追加（テンプレート改善提案 #1）
- [ ] idea-ship への harvest 自動トリガー追加（テンプレート改善提案 #5）
