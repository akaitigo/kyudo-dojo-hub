import { z } from "zod";
import { getLocalDateString, getOneYearAgoDateString } from "@/lib/date-utils";

/** 稽古日誌の入力バリデーションスキーマ */
export const practiceFormSchema = z.object({
	date: z
		.string()
		.min(1, "日付を入力してください")
		// タイムゾーン依存を避けるため、Date オブジェクトの絶対時刻比較ではなく
		// ローカル日付文字列 (YYYY-MM-DD) 同士の辞書順比較で判定する。
		.refine((val) => val <= getLocalDateString(), "未来日は入力できません")
		.refine((val) => val >= getOneYearAgoDateString(), "過去1年以内の日付を入力してください"),
	hitRate: z
		.number()
		.int("整数を入力してください")
		.min(0, "0以上を入力してください")
		.max(100, "100以下を入力してください"),
	arrowCount: z
		.number()
		.int("整数を入力してください")
		.min(1, "1以上を入力してください")
		.max(1000, "1000以下を入力してください"),
	notes: z.string().max(5000, "5,000文字以内で入力してください"),
	instructorComment: z.string().max(5000, "5,000文字以内で入力してください"),
});

export type PracticeFormValues = z.infer<typeof practiceFormSchema>;
