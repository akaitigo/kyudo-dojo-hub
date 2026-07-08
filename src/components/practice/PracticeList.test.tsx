// @vitest-environment jsdom
import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it } from "vitest";
import type { Practice, User } from "@/types/domain";
import { PracticeList } from "./PracticeList";

const user: User = {
	id: "user-001",
	name: "田中太郎",
	email: "tanaka@example.com",
	role: "practitioner",
	dan: "sandan",
	joinedAt: "2020-04-01",
	createdAt: "2020-04-01T00:00:00Z",
	updatedAt: "2020-04-01T00:00:00Z",
};

const practice: Practice = {
	id: "practice-001",
	userId: "user-001",
	date: "2026-04-01",
	hitRate: 70,
	arrowCount: 20,
	notes: "気づきメモ",
	instructorComment: "",
	createdAt: "2026-04-01T00:00:00Z",
	updatedAt: "2026-04-01T00:00:00Z",
};

describe("PracticeList（ファサード経由の users プロップで名前解決）", () => {
	afterEach(() => {
		cleanup();
	});

	it("props で渡された users からユーザー名と段位を表示する", () => {
		render(<PracticeList practices={[practice]} users={[user]} />);
		// モック直参照ではなく props の users から解決される
		expect(screen.getByText(/田中太郎 \(三段\)/)).toBeTruthy();
	});

	it("users に存在しない userId は id をフォールバック表示する", () => {
		const orphan: Practice = {
			...practice,
			id: "practice-002",
			userId: "unknown-user",
		};
		render(<PracticeList practices={[orphan]} users={[user]} />);
		expect(screen.getByText(/unknown-user/)).toBeTruthy();
	});
});
