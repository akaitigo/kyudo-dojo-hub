// =============================================================================
// Playwright E2E 設定
//
// Vite 開発サーバー (strictPort: 5173) を webServer として自動起動し、
// スモークテストを実行する。CI では GitHub Actions 上で同じ設定を使う。
// =============================================================================

import { defineConfig, devices } from "@playwright/test";

// Vite の dev server ポート。vite.config.ts の server.port と一致させること。
const PORT = 5173;
const BASE_URL = `http://localhost:${PORT}`;

export default defineConfig({
	// テストディレクトリ
	testDir: "./test/e2e",

	// テスト実行の並列数
	fullyParallel: true,

	// CI では未コミットの .only を検出したら失敗させる
	forbidOnly: !!process.env.CI,

	// CI では retry しない、ローカルでは1回リトライ
	retries: process.env.CI ? 0 : 1,

	// CI では並列ワーカー数を制限
	workers: process.env.CI ? 1 : undefined,

	// レポーター（CI では GitHub アノテーション + アーティファクト用 HTML レポート）
	reporter: process.env.CI ? [["github"], ["html", { open: "never" }]] : "html",

	// 共通設定
	use: {
		// dev server のURL
		baseURL: BASE_URL,

		// テスト失敗時にスクリーンショットを取得
		screenshot: "only-on-failure",

		// テスト失敗時にトレースを取得
		trace: "on-first-retry",
	},

	// テスト対象ブラウザ
	projects: [
		{
			name: "chromium",
			use: { ...devices["Desktop Chrome"] },
		},
		{
			name: "firefox",
			use: { ...devices["Desktop Firefox"] },
		},
		// モバイルビューポート
		{
			name: "mobile-chrome",
			use: { ...devices["Pixel 5"] },
		},
	],

	// dev server の自動起動
	webServer: {
		command: "npm run dev",
		url: BASE_URL,
		reuseExistingServer: !process.env.CI,
		timeout: 120_000,
	},
});
