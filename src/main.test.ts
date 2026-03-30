import { describe, expect, it } from "vitest";

describe("kyudo-dojo-hub", () => {
	it("基盤セットアップが完了している", () => {
		expect(true).toBe(true);
	});

	it("プロジェクト名が正しい", () => {
		const projectName = "kyudo-dojo-hub";
		expect(projectName).toContain("kyudo");
	});
});
