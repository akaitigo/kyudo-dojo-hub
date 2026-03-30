import { describe, expect, it } from "vitest";
import {
	createPractice,
	createReservation,
	deleteReservation,
	getAnalyses,
	getAnalysisByVideo,
	getDashboardSummary,
	getDojo,
	getDojos,
	getExamChecklist,
	getExamChecklists,
	getPractice,
	getPractices,
	getReservation,
	getReservations,
	getUser,
	getUsers,
	getUsersByDojo,
	getVideo,
	getVideos,
	toggleChecklistItem,
} from "./mock-api";

describe("モック API", () => {
	describe("Users", () => {
		it("全ユーザーを取得できる", async () => {
			const result = await getUsers();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThanOrEqual(10);
			}
		});

		it("IDでユーザーを取得できる", async () => {
			const result = await getUser("user-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.name).toBe("田中太郎");
			}
		});

		it("存在しないユーザーでエラーを返す", async () => {
			const result = await getUser("nonexistent");
			expect(result.success).toBe(false);
		});

		it("道場でユーザーをフィルタできる", async () => {
			const result = await getUsersByDojo("dojo-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThan(0);
				for (const user of result.data) {
					expect(user.dojoId).toBe("dojo-001");
				}
			}
		});
	});

	describe("Dojos", () => {
		it("全道場を取得できる", async () => {
			const result = await getDojos();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBe(2);
			}
		});

		it("IDで道場を取得できる", async () => {
			const result = await getDojo("dojo-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.name).toBe("東京弓道場");
			}
		});
	});

	describe("Practices", () => {
		it("全稽古記録を取得できる", async () => {
			const result = await getPractices();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThanOrEqual(10);
			}
		});

		it("ユーザーIDでフィルタできる", async () => {
			const result = await getPractices("user-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThan(0);
				for (const practice of result.data) {
					expect(practice.userId).toBe("user-001");
				}
			}
		});

		it("日付降順でソートされている", async () => {
			const result = await getPractices("user-001");
			expect(result.success).toBe(true);
			if (result.success && result.data.length >= 2) {
				const first = result.data[0]?.date ?? "";
				const second = result.data[1]?.date ?? "";
				expect(first >= second).toBe(true);
			}
		});

		it("IDで稽古記録を取得できる", async () => {
			const result = await getPractice("practice-001");
			expect(result.success).toBe(true);
		});

		it("新しい稽古記録を作成できる", async () => {
			const result = await createPractice({
				userId: "user-001",
				dojoId: "dojo-001",
				date: "2026-03-30",
				hitRate: 60,
				arrowCount: 36,
				notes: "テスト稽古",
				instructorComment: "",
			});
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.hitRate).toBe(60);
			}
		});

		it("的中率のバリデーション", async () => {
			const result = await createPractice({
				userId: "user-001",
				date: "2026-03-30",
				hitRate: 101,
				arrowCount: 36,
				notes: "",
				instructorComment: "",
			});
			expect(result.success).toBe(false);
		});

		it("矢数のバリデーション", async () => {
			const result = await createPractice({
				userId: "user-001",
				date: "2026-03-30",
				hitRate: 50,
				arrowCount: 0,
				notes: "",
				instructorComment: "",
			});
			expect(result.success).toBe(false);
		});
	});

	describe("Videos", () => {
		it("全動画を取得できる", async () => {
			const result = await getVideos();
			expect(result.success).toBe(true);
		});

		it("IDで動画を取得できる", async () => {
			const result = await getVideo("video-001");
			expect(result.success).toBe(true);
		});
	});

	describe("Analyses", () => {
		it("全分析結果を取得できる", async () => {
			const result = await getAnalyses();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThan(0);
			}
		});

		it("動画IDで分析結果を取得できる", async () => {
			const result = await getAnalysisByVideo("video-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(Object.keys(result.data.scores)).toHaveLength(8);
			}
		});

		it("八節スコアが8段階全てある", async () => {
			const result = await getAnalysisByVideo("video-001");
			expect(result.success).toBe(true);
			if (result.success) {
				const { scores } = result.data;
				expect(scores.ashibumi).toBeDefined();
				expect(scores.dozukuri).toBeDefined();
				expect(scores.yugamae).toBeDefined();
				expect(scores.uchiokoshi).toBeDefined();
				expect(scores.hikiwake).toBeDefined();
				expect(scores.kai).toBeDefined();
				expect(scores.hanare).toBeDefined();
				expect(scores.zanshin).toBeDefined();
			}
		});
	});

	describe("Reservations", () => {
		it("全予約を取得できる", async () => {
			const result = await getReservations();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThanOrEqual(10);
			}
		});

		it("道場と日付でフィルタできる", async () => {
			const result = await getReservations("dojo-001", "2026-03-30");
			expect(result.success).toBe(true);
			if (result.success) {
				for (const reservation of result.data) {
					expect(reservation.dojoId).toBe("dojo-001");
					expect(reservation.date).toBe("2026-03-30");
				}
			}
		});

		it("新しい予約を作成できる", async () => {
			const result = await createReservation({
				dojoId: "dojo-001",
				userId: "user-001",
				laneNumber: 6,
				date: "2026-04-01",
				startTime: "09:00",
				endTime: "10:00",
			});
			expect(result.success).toBe(true);
		});

		it("予約を削除できる", async () => {
			const result = await deleteReservation("res-001");
			expect(result.success).toBe(true);
		});

		it("存在しない予約の削除でエラー", async () => {
			const result = await deleteReservation("nonexistent");
			expect(result.success).toBe(false);
		});

		it("IDで予約を取得できる", async () => {
			const result = await getReservation("res-002");
			expect(result.success).toBe(true);
		});
	});

	describe("ExamChecklists", () => {
		it("全チェックリストを取得できる", async () => {
			const result = await getExamChecklists();
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.length).toBeGreaterThan(0);
			}
		});

		it("ユーザーIDでフィルタできる", async () => {
			const result = await getExamChecklists("user-001");
			expect(result.success).toBe(true);
			if (result.success) {
				for (const checklist of result.data) {
					expect(checklist.userId).toBe("user-001");
				}
			}
		});

		it("IDでチェックリストを取得できる", async () => {
			const result = await getExamChecklist("exam-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.targetDan).toBe("yondan");
			}
		});

		it("チェックリスト項目をトグルできる", async () => {
			const before = await getExamChecklist("exam-002");
			expect(before.success).toBe(true);
			if (!before.success) return;
			const itemBefore = before.data.items.find((i) => i.id === "item-011");

			const result = await toggleChecklistItem("exam-002", "item-011");
			expect(result.success).toBe(true);
			if (result.success && itemBefore) {
				const itemAfter = result.data.items.find((i) => i.id === "item-011");
				expect(itemAfter?.completed).toBe(!itemBefore.completed);
			}
		});
	});

	describe("Dashboard", () => {
		it("ダッシュボードサマリーを取得できる", async () => {
			const result = await getDashboardSummary("dojo-001");
			expect(result.success).toBe(true);
			if (result.success) {
				expect(result.data.totalMemberCount).toBeGreaterThan(0);
				expect(result.data.todayReservationCount).toBeDefined();
			}
		});
	});
});
