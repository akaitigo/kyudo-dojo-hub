// @vitest-environment jsdom
import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { getLocalDateString } from "@/lib/date-utils";
import { PracticeForm } from "./PracticeForm";

describe("PracticeForm 初期日付（テスト安定化）", () => {
	afterEach(() => {
		cleanup();
		vi.useRealTimers();
	});

	it("注入した defaultDate を日付入力の初期値に使う", () => {
		render(<PracticeForm onSubmit={vi.fn()} defaultDate="2026-01-15" />);
		const input = screen.getByLabelText("日付") as HTMLInputElement;
		expect(input.value).toBe("2026-01-15");
	});

	it("defaultDate 未指定時は現在のローカル日付を使う（new Date に依存せず固定可能）", () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date("2026-04-05T05:00:00Z"));
		render(<PracticeForm onSubmit={vi.fn()} />);
		const input = screen.getByLabelText("日付") as HTMLInputElement;
		expect(input.value).toBe(getLocalDateString());
	});
});
