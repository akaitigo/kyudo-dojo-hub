import { describe, expect, it } from "vitest";
import { analyzeVideo, getUsers } from "@/lib/api";

// テスト環境では VITE_API_BASE_URL が未設定のため、ファサードはモック実装に
// フォールバックする。ここではファサードがモック経由でデータを返すこと
// （＝ mock/real 切り替えが api.ts に一元化されていること）を検証する。
describe("api ファサード（モックフォールバック）", () => {
	it("getUsers がモックのユーザー一覧を返す", async () => {
		const result = await getUsers();
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data.length).toBeGreaterThanOrEqual(10);
		}
	});

	it("analyzeVideo が videoId から決定的な分析結果を返す", async () => {
		const first = await analyzeVideo({
			videoId: "video-xyz",
			userId: "user-001",
		});
		const second = await analyzeVideo({
			videoId: "video-xyz",
			userId: "user-001",
		});
		expect(first.success && second.success).toBe(true);
		if (first.success && second.success) {
			// スコアは videoId から決定的に生成されるため 2 回の呼び出しで一致する
			expect(first.data.scores).toEqual(second.data.scores);
			expect(first.data.overallScore).toBe(second.data.overallScore);
			expect(first.data.phases.length).toBe(8);
			expect(first.data.overallScore).toBeGreaterThanOrEqual(0);
			expect(first.data.overallScore).toBeLessThanOrEqual(100);
		}
	});
});
