/**
 * 実 API 実装層
 *
 * Go バックエンド (REST API) に HTTP リクエストを送信する。低レベルの
 * HTTP・エラーハンドリングは api-client.ts に委譲し、mock-api.ts と同じ
 * インターフェースを提供する。呼び出し側の切り替えは api.ts で一元管理する。
 */
import { apiDelete, apiGet, apiPatch, apiPost } from "@/lib/api-client";
import type { Analysis, ApiResult, Dojo, ExamChecklist, Practice, Reservation, User, Video } from "@/types/domain";

// ---------------------------------------------------------------------------
// Users API
// ---------------------------------------------------------------------------

export async function getUsers(): Promise<ApiResult<readonly User[]>> {
	return apiGet<readonly User[]>("/api/users");
}

export async function getUser(id: string): Promise<ApiResult<User>> {
	return apiGet<User>(`/api/users/${encodeURIComponent(id)}`);
}

export async function getUsersByDojo(dojoId: string): Promise<ApiResult<readonly User[]>> {
	return apiGet<readonly User[]>(`/api/users?dojoId=${encodeURIComponent(dojoId)}`);
}

// ---------------------------------------------------------------------------
// Dojos API
// ---------------------------------------------------------------------------

export async function getDojos(): Promise<ApiResult<readonly Dojo[]>> {
	return apiGet<readonly Dojo[]>("/api/dojos");
}

export async function getDojo(id: string): Promise<ApiResult<Dojo>> {
	return apiGet<Dojo>(`/api/dojos/${encodeURIComponent(id)}`);
}

// ---------------------------------------------------------------------------
// Practices API
// ---------------------------------------------------------------------------

export async function getPractices(userId?: string): Promise<ApiResult<readonly Practice[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiGet<readonly Practice[]>(`/api/practices${qs}`);
}

export async function getPractice(id: string): Promise<ApiResult<Practice>> {
	return apiGet<Practice>(`/api/practices/${encodeURIComponent(id)}`);
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
	return apiPost<Practice>("/api/practices", input);
}

// ---------------------------------------------------------------------------
// Videos API
// ---------------------------------------------------------------------------

export async function getVideos(userId?: string): Promise<ApiResult<readonly Video[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiGet<readonly Video[]>(`/api/videos${qs}`);
}

export async function getVideo(id: string): Promise<ApiResult<Video>> {
	return apiGet<Video>(`/api/videos/${encodeURIComponent(id)}`);
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
	return apiPost<Video>("/api/videos", input);
}

// ---------------------------------------------------------------------------
// Analyses API
// ---------------------------------------------------------------------------

export async function getAnalyses(userId?: string): Promise<ApiResult<readonly Analysis[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiGet<readonly Analysis[]>(`/api/analyses${qs}`);
}

export async function getAnalysis(id: string): Promise<ApiResult<Analysis>> {
	return apiGet<Analysis>(`/api/analyses/${encodeURIComponent(id)}`);
}

export async function getAnalysisByVideo(videoId: string): Promise<ApiResult<Analysis>> {
	return apiGet<Analysis>(`/api/analyses/video/${encodeURIComponent(videoId)}`);
}

export interface AnalyzeVideoInput {
	readonly videoId: string;
	readonly userId: string;
}

export async function analyzeVideo(input: AnalyzeVideoInput): Promise<ApiResult<Analysis>> {
	return apiPost<Analysis>("/api/analyses/analyze", input);
}

// ---------------------------------------------------------------------------
// Reservations API
// ---------------------------------------------------------------------------

export async function getReservations(dojoId?: string, date?: string): Promise<ApiResult<readonly Reservation[]>> {
	const params = new URLSearchParams();
	if (dojoId) params.set("dojoId", dojoId);
	if (date) params.set("date", date);
	const qs = params.toString() ? `?${params.toString()}` : "";
	return apiGet<readonly Reservation[]>(`/api/reservations${qs}`);
}

export async function getReservation(id: string): Promise<ApiResult<Reservation>> {
	return apiGet<Reservation>(`/api/reservations/${encodeURIComponent(id)}`);
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
	return apiPost<Reservation>("/api/reservations", input);
}

export async function deleteReservation(id: string): Promise<ApiResult<{ readonly deleted: true }>> {
	return apiDelete<{ readonly deleted: true }>(`/api/reservations/${encodeURIComponent(id)}`);
}

// ---------------------------------------------------------------------------
// Exam Checklists API
// ---------------------------------------------------------------------------

export async function getExamChecklists(userId?: string): Promise<ApiResult<readonly ExamChecklist[]>> {
	const qs = userId ? `?userId=${encodeURIComponent(userId)}` : "";
	return apiGet<readonly ExamChecklist[]>(`/api/exam-checklists${qs}`);
}

export async function getExamChecklist(id: string): Promise<ApiResult<ExamChecklist>> {
	return apiGet<ExamChecklist>(`/api/exam-checklists/${encodeURIComponent(id)}`);
}

export async function toggleChecklistItem(checklistId: string, itemId: string): Promise<ApiResult<ExamChecklist>> {
	return apiPatch<ExamChecklist>(
		`/api/exam-checklists/${encodeURIComponent(checklistId)}/items/${encodeURIComponent(itemId)}/toggle`,
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
	return apiGet<DashboardSummary>(`/api/dashboard/${encodeURIComponent(dojoId)}`);
}
