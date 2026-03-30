import { describe, expect, it } from "vitest";
import { generateTimeSlots, getEndTime, reservationFormSchema } from "./reservation-validation";

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
});
