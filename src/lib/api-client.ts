/**
 * Real API client — used when VITE_API_BASE_URL is set.
 *
 * Communicates with the Go REST API backend.
 */
import type { ApiError, ApiResult } from "@/types/domain";

const BASE_URL = import.meta.env["VITE_API_BASE_URL"] as string | undefined;

/**
 * Whether the real API backend is configured.
 * When false, the application should fall back to mock-api.
 */
export function isRealApiEnabled(): boolean {
	return typeof BASE_URL === "string" && BASE_URL.length > 0;
}

/** Go API の成功レスポンス構造を検証する型ガード */
function isApiSuccess<T>(value: unknown): value is { success: true; data: T } {
	return (
		typeof value === "object" &&
		value !== null &&
		"success" in value &&
		(value as Record<string, unknown>)["success"] === true &&
		"data" in value
	);
}

/** Go API のエラーレスポンス構造を検証する型ガード */
function isApiError(value: unknown): value is ApiError {
	if (typeof value !== "object" || value === null || !("success" in value)) {
		return false;
	}
	const obj = value as Record<string, unknown>;
	if (obj["success"] !== false || typeof obj["error"] !== "object" || obj["error"] === null) {
		return false;
	}
	const err = obj["error"] as Record<string, unknown>;
	return typeof err["code"] === "string" && typeof err["message"] === "string";
}

/** エラー結果を生成する */
function errorResult(code: string, message: string): ApiError {
	return { success: false, error: { code, message } };
}

async function request<T>(method: string, path: string, body?: unknown): Promise<ApiResult<T>> {
	const url = `${BASE_URL}${path}`;
	const headers: Record<string, string> = {
		Accept: "application/json",
	};
	if (body !== undefined) {
		headers["Content-Type"] = "application/json";
	}

	let res: Response;
	try {
		res = await fetch(url, {
			method,
			headers,
			body: body !== undefined ? JSON.stringify(body) : null,
		});
	} catch (err: unknown) {
		// ネットワーク断・DNS 失敗・CORS 拒否など fetch 自体が throw するケース。
		const message = err instanceof Error ? err.message : "通信エラーが発生しました";
		return errorResult("NETWORK_ERROR", message);
	}

	// レスポンスボディの JSON パース。HTTP エラー時にボディが空／非 JSON の
	// ことがあるため、パース失敗も個別にハンドリングする。
	let json: unknown;
	try {
		json = await res.json();
	} catch {
		if (!res.ok) {
			return errorResult("HTTP_ERROR", `サーバーエラーが発生しました (HTTP ${res.status})`);
		}
		return errorResult("INVALID_RESPONSE", "レスポンスの解析に失敗しました");
	}

	// Go API は成功・失敗どちらも { success, data/error } 形式で返す。
	if (isApiError(json)) {
		return json;
	}
	if (res.ok && isApiSuccess<T>(json)) {
		return json;
	}

	// res.ok でない、または想定外の形状。ステータスコードを含めて返す。
	if (!res.ok) {
		return errorResult("HTTP_ERROR", `サーバーエラーが発生しました (HTTP ${res.status})`);
	}
	return errorResult("INVALID_RESPONSE", "予期しないレスポンス形式です");
}

function get<T>(path: string): Promise<ApiResult<T>> {
	return request<T>("GET", path);
}

function post<T>(path: string, body: unknown): Promise<ApiResult<T>> {
	return request<T>("POST", path, body);
}

function patch<T>(path: string, body?: unknown): Promise<ApiResult<T>> {
	return request<T>("PATCH", path, body);
}

function del<T>(path: string): Promise<ApiResult<T>> {
	return request<T>("DELETE", path);
}

// ---------------------------------------------------------------------------
// Exports matching mock-api interface
// ---------------------------------------------------------------------------

export { del as apiDelete, get as apiGet, patch as apiPatch, post as apiPost };
