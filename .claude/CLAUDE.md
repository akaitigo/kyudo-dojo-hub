# kyudo-dojo-hub アーキテクチャ概要

## システム構成

```
React SPA (Vite) --> Go API (gRPC/REST) --> PostgreSQL / Redis
                          |
                          v
                  Python MediaPipe Worker (非同期)
                          |
                          v
                  GCP Cloud Storage (動画)
```

## 主要な設計判断

- ADR-001: (未作成) MVP では TypeScript/React フロントエンドに集中し、Go API と Python ML はモックで開始
- ADR-002: (未作成) 射法八節の判定はルールベース分類から開始、段階的に ML 精度を向上

## 外部サービス連携

| サービス | 用途 | 認証方式 |
|----------|------|----------|
| GCP Cloud Storage | 動画保存 | サービスアカウント |
| GCP Cloud Run | API/Worker デプロイ | IAM |
| MediaPipe | 骨格推定 | ローカル実行（API不要） |

## データモデル概要

- User: 稽古者・道場管理者
- Dojo: 道場情報・的場
- Practice: 稽古記録（的中率・気づき・師範コメント）
- Video: 射形動画メタデータ
- Analysis: 骨格推定結果・八節フェーズ分割
- Reservation: 的場予約
- ExamChecklist: 段位審査チェックリスト
