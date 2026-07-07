import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { getLocalDateString, getOneYearAgoDateString } from "./date-utils";
import { practiceFormSchema } from "./validation";

describe("practiceFormSchema", () => {
	it("有効な入力を受け付ける", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
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
			date: getLocalDateString(),
			hitRate: 101,
			arrowCount: 40,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("的中率 -1 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
			hitRate: -1,
			arrowCount: 40,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("矢数 0 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
			hitRate: 50,
			arrowCount: 0,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("矢数 1001 を拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
			hitRate: 50,
			arrowCount: 1001,
			notes: "",
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("5001文字の気づきを拒否する", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
			hitRate: 50,
			arrowCount: 40,
			notes: "a".repeat(5001),
			instructorComment: "",
		});
		expect(result.success).toBe(false);
	});

	it("5000文字の気づきを受け付ける", () => {
		const result = practiceFormSchema.safeParse({
			date: getLocalDateString(),
			hitRate: 50,
			arrowCount: 40,
			notes: "a".repeat(5000),
			instructorComment: "",
		});
		expect(result.success).toBe(true);
	});
});

describe("practiceFormSchema 日付バリデーション（タイムゾーン非依存）", () => {
	const base = {
		hitRate: 50,
		arrowCount: 40,
		notes: "",
		instructorComment: "",
	} as const;

	beforeEach(() => {
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	// JST の 0:00〜8:59 に相当する早朝の瞬間（UTC 上の日付とローカル日付が
	// ずれ得る時間帯）でも、当日の日付が「未来日」と誤判定されないこと。
	it("早朝の時刻でも当日の日付を受け付ける", () => {
		vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
		const result = practiceFormSchema.safeParse({
			...base,
			date: getLocalDateString(),
		});
		expect(result.success).toBe(true);
	});

	it("翌日の日付は未来日として拒否する", () => {
		vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
		const tomorrow = getLocalDateString(new Date("2026-04-06T05:00:00Z"));
		const result = practiceFormSchema.safeParse({ ...base, date: tomorrow });
		expect(result.success).toBe(false);
	});

	it("ちょうど1年前の日付を受け付ける", () => {
		vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
		const result = practiceFormSchema.safeParse({
			...base,
			date: getOneYearAgoDateString(),
		});
		expect(result.success).toBe(true);
	});

	it("1年より前の日付を拒否する", () => {
		vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
		const twoYearsAgo = getLocalDateString(new Date("2024-04-05T05:00:00Z"));
		const result = practiceFormSchema.safeParse({ ...base, date: twoYearsAgo });
		expect(result.success).toBe(false);
	});
});
