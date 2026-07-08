// @vitest-environment jsdom
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { getExamChecklists } from "@/lib/api";
import { ExamChecklistPage } from "./ExamChecklistPage";

vi.mock("@/lib/api", () => ({
	getExamChecklists: vi.fn(),
	toggleChecklistItem: vi.fn(),
}));

describe("ExamChecklistPage エラーUI", () => {
	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	it("チェックリストの読み込み失敗時に role=alert を表示する", async () => {
		vi.mocked(getExamChecklists).mockResolvedValue({
			success: false,
			error: { code: "HTTP_ERROR", message: "サーバーエラーが発生しました" },
		});

		render(<ExamChecklistPage />);

		const alert = await screen.findByRole("alert");
		expect(alert.textContent).toContain("サーバーエラーが発生しました");
	});

	it("空データ（成功）時はエラーではなく空メッセージを表示する", async () => {
		vi.mocked(getExamChecklists).mockResolvedValue({ success: true, data: [] });

		render(<ExamChecklistPage />);

		await waitFor(() => {
			expect(screen.getByText("チェックリストがありません")).toBeTruthy();
		});
		expect(screen.queryByRole("alert")).toBeNull();
	});
});
