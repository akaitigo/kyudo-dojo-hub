import { afterEach, describe, expect, it, vi } from "vitest";
import { apiGet, apiPost } from "./api-client";

/** fetch のモックレスポンスを生成する */
function mockResponse(options: { ok: boolean; status: number; json: () => unknown | Promise<unknown> }): Response {
	return {
		ok: options.ok,
		status: options.status,
		json: async () => options.json(),
	} as unknown as Response;
}

describe("api-client request", () => {
	afterEach(() => {
		vi.unstubAllGlobals();
		vi.restoreAllMocks();
	});

	it("成功レスポンスをそのまま返す", async () => {
		vi.stubGlobal(
			"fetch",
			vi.fn().mockResolvedValue(
				mockResponse({
					ok: true,
					status: 200,
					json: () => ({ success: true, data: [{ id: "u1" }] }),
				}),
			),
		);
		const result = await apiGet<Array<{ id: string }>>("/api/users");
		expect(result.success).toBe(true);
		if (result.success) {
			expect(result.data[0]?.id).toBe("u1");
		}
	});

	it("API エラーレスポンス（4xx）をエラー結果として返す", async () => {
		vi.stubGlobal(
			"fetch",
			vi.fn().mockResolvedValue(
				mockResponse({
					ok: false,
					status: 404,
					json: () => ({
						success: false,
						error: { code: "NOT_FOUND", message: "見つかりません" },
					}),
				}),
			),
		);
		const result = await apiGet("/api/users/x");
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.code).toBe("NOT_FOUND");
		}
	});

	it("HTTP エラーかつ非 JSON ボディを HTTP_ERROR として返す", async () => {
		vi.stubGlobal(
			"fetch",
			vi.fn().mockResolvedValue(
				mockResponse({
					ok: false,
					status: 500,
					json: () => {
						throw new Error("Unexpected token");
					},
				}),
			),
		);
		const result = await apiGet("/api/users");
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.code).toBe("HTTP_ERROR");
			expect(result.error.message).toContain("500");
		}
	});

	it("fetch 自体が throw した場合は NETWORK_ERROR を返す", async () => {
		vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("Failed to fetch")));
		const result = await apiPost("/api/practices", { hitRate: 50 });
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.code).toBe("NETWORK_ERROR");
			expect(result.error.message).toBe("Failed to fetch");
		}
	});

	it("res.ok だが想定外の形状は INVALID_RESPONSE を返す", async () => {
		vi.stubGlobal(
			"fetch",
			vi.fn().mockResolvedValue(
				mockResponse({
					ok: true,
					status: 200,
					json: () => ({ unexpected: true }),
				}),
			),
		);
		const result = await apiGet("/api/users");
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.code).toBe("INVALID_RESPONSE");
		}
	});
});
