// @vitest-environment jsdom
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { getPractices, getUsers } from "@/lib/api";
import { PracticesPage } from "./PracticesPage";

vi.mock("@/lib/api", () => ({
	getPractices: vi.fn(),
	getUsers: vi.fn(),
	createPractice: vi.fn(),
}));
// recharts はテスト環境で描画寸法を持たないためスタブ化する
vi.mock("@/components/practice/HitRateChart", () => ({
	HitRateChart: () => null,
}));

describe("PracticesPage エラーUI", () => {
	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	it("稽古記録の読み込み失敗時に role=alert を表示する", async () => {
		vi.mocked(getPractices).mockResolvedValue({
			success: false,
			error: { code: "NETWORK_ERROR", message: "通信エラーが発生しました" },
		});
		vi.mocked(getUsers).mockResolvedValue({ success: true, data: [] });

		render(<PracticesPage />);

		const alert = await screen.findByRole("alert");
		expect(alert.textContent).toContain("通信エラーが発生しました");
	});

	it("読み込み成功時はエラーを表示しない", async () => {
		vi.mocked(getPractices).mockResolvedValue({ success: true, data: [] });
		vi.mocked(getUsers).mockResolvedValue({ success: true, data: [] });

		render(<PracticesPage />);

		await waitFor(() => {
			expect(screen.getByText("稽古履歴")).toBeTruthy();
		});
		expect(screen.queryByRole("alert")).toBeNull();
	});
});
