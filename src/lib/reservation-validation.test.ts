import { afterEach, describe, expect, it, vi } from "vitest";
import { getLocalDateString } from "./date-utils";
import {
	generateTimeSlots,
	getEndTime,
	getReservationFormDefaults,
	reservationFormSchema,
} from "./reservation-validation";

describe("reservation-validation", () => {
	describe("reservationFormSchema", () => {
		it("有効な入力を受け付ける", () => {
			const result = reservationFormSchema.safeParse({
				date: "2026-04-01",
				startTime: "10:00",
				laneNumber: 1,
			});
			expect(result.success).toBe(true);
		});

		it("空の日付を拒否する", () => {
			const result = reservationFormSchema.safeParse({
				date: "",
				startTime: "10:00",
				laneNumber: 1,
			});
			expect(result.success).toBe(false);
		});

		it("空の開始時刻を拒否する", () => {
			const result = reservationFormSchema.safeParse({
				date: "2026-04-01",
				startTime: "",
				laneNumber: 1,
			});
			expect(result.success).toBe(false);
		});
	});

	describe("generateTimeSlots", () => {
		it("9:00-21:00で12枠を生成する", () => {
			const slots = generateTimeSlots("09:00", "21:00");
			expect(slots).toHaveLength(12);
			expect(slots[0]).toBe("09:00");
			expect(slots[11]).toBe("20:00");
		});

		it("10:00-20:00で10枠を生成する", () => {
			const slots = generateTimeSlots("10:00", "20:00");
			expect(slots).toHaveLength(10);
		});
	});

	describe("getEndTime", () => {
		it("1時間後の時刻を返す", () => {
			expect(getEndTime("09:00")).toBe("10:00");
			expect(getEndTime("14:00")).toBe("15:00");
			expect(getEndTime("20:00")).toBe("21:00");
		});
	});

	describe("getReservationFormDefaults", () => {
		afterEach(() => {
			vi.useRealTimers();
		});

		it("注入した日付を初期値に使う（テスト安定化）", () => {
			const defaults = getReservationFormDefaults("2026-04-01");
			expect(defaults).toEqual({
				date: "2026-04-01",
				startTime: "",
				laneNumber: 1,
			});
		});

		it("未指定時は現在のローカル日付を使う", () => {
			vi.useFakeTimers();
			vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
			expect(getReservationFormDefaults().date).toBe(getLocalDateString());
		});
	});
});
