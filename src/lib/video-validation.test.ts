import { describe, expect, it } from "vitest";
import { validateVideoDuration, validateVideoFile } from "./video-validation";

describe("video-validation", () => {
	describe("validateVideoFile", () => {
		it("mp4 ファイルを受け付ける", () => {
			const file = new File(["dummy"], "test.mp4", { type: "video/mp4" });
			expect(validateVideoFile(file).valid).toBe(true);
		});

		it("quicktime (mov) ファイルを受け付ける", () => {
			const file = new File(["dummy"], "test.mov", { type: "video/quicktime" });
			expect(validateVideoFile(file).valid).toBe(true);
		});

		it("webm ファイルを受け付ける", () => {
			const file = new File(["dummy"], "test.webm", { type: "video/webm" });
			expect(validateVideoFile(file).valid).toBe(true);
		});

		it("avi ファイルを拒否する", () => {
			const file = new File(["dummy"], "test.avi", { type: "video/x-msvideo" });
			const result = validateVideoFile(file);
			expect(result.valid).toBe(false);
			expect(result.error).toContain("mp4");
		});

		it("画像ファイルを拒否する", () => {
			const file = new File(["dummy"], "test.jpg", { type: "image/jpeg" });
			const result = validateVideoFile(file);
			expect(result.valid).toBe(false);
		});
	});

	describe("validateVideoDuration", () => {
		it("5分以下を受け付ける", () => {
			expect(validateVideoDuration(300).valid).toBe(true);
			expect(validateVideoDuration(60).valid).toBe(true);
			expect(validateVideoDuration(0).valid).toBe(true);
		});

		it("5分超を拒否する", () => {
			const result = validateVideoDuration(301);
			expect(result.valid).toBe(false);
			expect(result.error).toContain("5 分");
		});
	});
});
