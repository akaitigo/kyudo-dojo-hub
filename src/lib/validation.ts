import { z } from "zod";

/** 稽古日誌の入力バリデーションスキーマ */
export const practiceFormSchema = z.object({
	date: z
		.string()
		.min(1, "日付を入力してください")
		.refine((val) => {
			const d = new Date(val);
			const now = new Date();
			now.setHours(23, 59, 59, 999);
			return d <= now;
		}, "未来日は入力できません")
		.refine((val) => {
			const d = new Date(val);
			const oneYearAgo = new Date();
			oneYearAgo.setFullYear(oneYearAgo.getFullYear() - 1);
			return d >= oneYearAgo;
		}, "過去1年以内の日付を入力してください"),
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
