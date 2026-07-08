import { z } from "zod";
import { getLocalDateString } from "@/lib/date-utils";

/** 予約フォームのバリデーションスキーマ */
export const reservationFormSchema = z.object({
	date: z.string().min(1, "日付を入力してください"),
	startTime: z.string().min(1, "開始時刻を選択してください"),
	laneNumber: z.number().int().min(1, "的場番号を選択してください"),
});

export type ReservationFormValues = z.infer<typeof reservationFormSchema>;

/**
 * 予約フォームの初期値を返す。日付はテストで安定させられるよう注入可能。
 * 未指定時のみ現在日付（ローカル）を用いる。
 */
export function getReservationFormDefaults(today: string = getLocalDateString()): ReservationFormValues {
	return { date: today, startTime: "", laneNumber: 1 };
}

/** 営業時間内の1時間単位の時間帯を生成 */
export function generateTimeSlots(openTime: string, closeTime: string): readonly string[] {
	const openHour = Number.parseInt(openTime.split(":")[0] ?? "9", 10);
	const closeHour = Number.parseInt(closeTime.split(":")[0] ?? "21", 10);
	const slots: string[] = [];
	for (let h = openHour; h < closeHour; h++) {
		slots.push(`${String(h).padStart(2, "0")}:00`);
	}
	return slots;
}

/** 終了時刻を計算（1時間単位） */
export function getEndTime(startTime: string): string {
	const hour = Number.parseInt(startTime.split(":")[0] ?? "0", 10);
	return `${String(hour + 1).padStart(2, "0")}:00`;
}
