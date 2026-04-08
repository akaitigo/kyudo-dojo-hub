/**
 * API クライアント層
 *
 * Go バックエンド (REST API) に HTTP リクエストを送信する。
 * mock-api.ts と同じインターフェースを維持し、差し替え可能にする。
 */
import type { Analysis, ApiResult, Dojo, ExamChecklist, Practice, Reservation, User, Video } from "@/types/domain";

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const API_BASE_URL = (import.meta.env["VITE_API_BASE_URL"] as string | undefined) ?? "";

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

async function apiFetch<T>(path: string, init?: RequestInit): Promise<ApiResult<T>> {
	try {
		const response = await fetch(`${API_BASE_URL}${path}`, {
			headers: { "Content-Type": "application/json" },
			...init,
		});

		const json: unknown = await response.json();

		// The Go API returns { success: true/false, data/error }
		if (typeof json === "object" && json !== null && "success" in json) {
			return json as ApiResult<T>;
		}

		return {
			success: false,
			error: { code: "UNKNOWN", message: "予期しないレスポンス形式です" },
		};
	} catch (err: unknown) {
		const message = err instanceof Error ? err.message : "通信エラーが発生しました";
		return {
			success: false,
			error: { code: "NETWORK_ERROR", message },
		};
	}
}

// ---------------------------------------------------------------------------
// Users API
// ---------------------------------------------------------------------------

export async function getUsers(): Promise<ApiResult<readonly User[]>> {
	return apiFetch<readonly User[]>("/api/users");
}

export async function getUser(id: string): Promise<ApiResult<User>> {
	return apiFetch<User>(`/api/users/${encodeURIComponent(id)}`);
}

export async function getUsersByDojo(dojoId: string): Promise<ApiResult<readonly User[]>> {
	return apiFetch<readonly User[]>(`/api/users?dojoId=${encodeURIComponent(dojoId)}`);
}

// ---------------------------------------------------------------------------
// Dojos API
// ---------------------------------------------------------------------------

export async function getDojos(): Promise<ApiResult<readonly Dojo[]>> {
	return apiFetch<readonly Dojo[]>("/api/dojos");
}

export async function getDojo(id: string): Promise<ApiResult<Dojo>> {
	return apiFetch<Dojo>(`/api/dojos/${encodeURIComponent(id)}`);
}

// ---------------------------------------------------------------------------
// Practices API
// ---------------------------------------------------------------------------

export async function getPractices(userId?: string): Promise<ApiResult<readonly Practice[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiFetch<readonly Practice[]>(`/api/practices${qs}`);
}

export async function getPractice(id: string): Promise<ApiResult<Practice>> {
	return apiFetch<Practice>(`/api/practices/${encodeURIComponent(id)}`);
}

export interface CreatePracticeInput {
	readonly userId: string;
	readonly dojoId?: string;
	readonly date: string;
	readonly hitRate: number;
	readonly arrowCount: number;
	readonly notes: string;
	readonly instructorComment: string;
}

export async function createPractice(input: CreatePracticeInput): Promise<ApiResult<Practice>> {
	return apiFetch<Practice>("/api/practices", {
		method: "POST",
		body: JSON.stringify(input),
	});
}

// ---------------------------------------------------------------------------
// Videos API
// ---------------------------------------------------------------------------

export async function getVideos(userId?: string): Promise<ApiResult<readonly Video[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiFetch<readonly Video[]>(`/api/videos${qs}`);
}

export async function getVideo(id: string): Promise<ApiResult<Video>> {
	return apiFetch<Video>(`/api/videos/${encodeURIComponent(id)}`);
}

export interface CreateVideoInput {
	readonly userId: string;
	readonly practiceId?: string;
	readonly fileName: string;
	readonly fileSize: number;
	readonly duration: number;
	readonly mimeType: "video/mp4" | "video/quicktime" | "video/webm";
	readonly url: string;
}

export async function createVideo(input: CreateVideoInput): Promise<ApiResult<Video>> {
	return apiFetch<Video>("/api/videos", {
		method: "POST",
		body: JSON.stringify(input),
	});
}

// ---------------------------------------------------------------------------
// Analyses API
// ---------------------------------------------------------------------------

export async function getAnalyses(userId?: string): Promise<ApiResult<readonly Analysis[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiFetch<readonly Analysis[]>(`/api/analyses${qs}`);
}

export async function getAnalysis(id: string): Promise<ApiResult<Analysis>> {
	return apiFetch<Analysis>(`/api/analyses/${encodeURIComponent(id)}`);
}

export async function getAnalysisByVideo(videoId: string): Promise<ApiResult<Analysis>> {
	return apiFetch<Analysis>(`/api/analyses/video/${encodeURIComponent(videoId)}`);
}

export interface AnalyzeVideoInput {
	readonly videoId: string;
	readonly userId: string;
}

export async function analyzeVideo(input: AnalyzeVideoInput): Promise<ApiResult<Analysis>> {
	return apiFetch<Analysis>("/api/analyses/analyze", {
		method: "POST",
		body: JSON.stringify(input),
	});
}

// ---------------------------------------------------------------------------
// Reservations API
// ---------------------------------------------------------------------------

export async function getReservations(dojoId?: string, date?: string): Promise<ApiResult<readonly Reservation[]>> {
	const params = new URLSearchParams();
	if (dojoId) params.set("dojoId", dojoId);
	if (date) params.set("date", date);
	const qs = params.toString() ? `?${params.toString()}` : "";
	return apiFetch<readonly Reservation[]>(`/api/reservations${qs}`);
}

export async function getReservation(id: string): Promise<ApiResult<Reservation>> {
	return apiFetch<Reservation>(`/api/reservations/${encodeURIComponent(id)}`);
}

export interface CreateReservationInput {
	readonly dojoId: string;
	readonly userId: string;
	readonly laneNumber: number;
	readonly date: string;
	readonly startTime: string;
	readonly endTime: string;
}

export async function createReservation(input: CreateReservationInput): Promise<ApiResult<Reservation>> {
	return apiFetch<Reservation>("/api/reservations", {
		method: "POST",
		body: JSON.stringify(input),
	});
}

export async function deleteReservation(id: string): Promise<ApiResult<{ readonly deleted: true }>> {
	return apiFetch<{ readonly deleted: true }>(`/api/reservations/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
}

// ---------------------------------------------------------------------------
// Exam Checklists API
// ---------------------------------------------------------------------------

export async function getExamChecklists(userId?: string): Promise<ApiResult<readonly ExamChecklist[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiFetch<readonly ExamChecklist[]>(`/api/exam-checklists${qs}`);
}

export async function getExamChecklist(id: string): Promise<ApiResult<ExamChecklist>> {
	return apiFetch<ExamChecklist>(`/api/exam-checklists/${encodeURIComponent(id)}`);
}

export async function toggleChecklistItem(checklistId: string, itemId: string): Promise<ApiResult<ExamChecklist>> {
	return apiFetch<ExamChecklist>(
		`/api/exam-checklists/${encodeURIComponent(checklistId)}/items/${encodeURIComponent(itemId)}/toggle`,
		{ method: "PATCH" },
	);
}

// ---------------------------------------------------------------------------
// Dashboard API
// ---------------------------------------------------------------------------

export interface DashboardSummary {
	readonly todayReservationCount: number;
	readonly totalMemberCount: number;
	readonly todayReservations: readonly Reservation[];
}

export async function getDashboardSummary(dojoId: string): Promise<ApiResult<DashboardSummary>> {
	return apiFetch<DashboardSummary>(`/api/dashboard/${encodeURIComponent(dojoId)}`);
}
