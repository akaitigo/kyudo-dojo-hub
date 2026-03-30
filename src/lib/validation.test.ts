import { describe, expect, it } from "vitest";
import { practiceFormSchema } from "./validation";

describe("practiceFormSchema", () => {
	it("有効な入力を受け付ける", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 65,
			arrowCount: 40,
			notes: "テスト",
			instructorComment: "",
		});
		expect(result.success).toBe(true);
	});

	it("空の日付を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: "",
			hitRate: 65,
			arrowCount: 40,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("的中率 101 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 101,
			arrowCount: 40,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("的中率 -1 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: -1,
			arrowCount: 40,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("矢数 0 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 50,
			arrowCount: 0,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("矢数 1001 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 50,
			arrowCount: 1001,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("5001文字の気づきを拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 50,
			arrowCount: 40,
			notes: "a".repeat(5001),
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("5000文字の気づきを受け付ける", () => {
		const result = practiceFormSchema.safeParse({
			date: new Date().toISOString().split("T")[0],
			hitRate: 50,
			arrowCount: 40,
			notes: "a".repeat(5000),
			instructorComment: "",
		});
		expect(result.success).toBe(true);
	});
});
