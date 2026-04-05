/**
 * Real API client — used when VITE_API_BASE_URL is set.
 *
 * Communicates with the Go REST API backend.
 */
import type { ApiResult } from "@/types/domain";

const BASE_URL = import.meta.env["VITE_API_BASE_URL"] as string | undefined;

/**
 * Whether the real API backend is configured.
 * When false, the application should fall back to mock-api.
 */
export function isRealApiEnabled(): boolean {
	return typeof BASE_URL === "string" && BASE_URL.length > 0;
}

async function request<T>(method: string, path: string, body?: unknown): Promise<ApiResult<T>> {
	const url = `${BASE_URL}${path}`;
	const headers: Record<string, string> = {
		Accept: "application/json",
	};
	if (body !== undefined) {
		headers["Content-Type"] = "application/json";
	}

	const res = await fetch(url, {
		method,
		headers,
		body: body !== undefined ? JSON.stringify(body) : null,
	});

	const json = (await res.json()) as ApiResult<T>;
	return json;
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
