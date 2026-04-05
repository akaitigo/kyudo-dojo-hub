/**
 * Unified API facade.
 *
 * When VITE_API_BASE_URL is set, delegates to the real Go API.
 * Otherwise, falls back to the in-memory mock API.
 */
import { apiDelete, apiGet, apiPatch, apiPost, isRealApiEnabled } from "@/lib/api-client";
import type {
	DashboardSummary,
	CreatePracticeInput as MockCreatePracticeInput,
	CreateReservationInput as MockCreateReservationInput,
} from "@/lib/mock-api";
import * as mockApi from "@/lib/mock-api";
import type { Analysis, ApiResult, Dojo, ExamChecklist, Practice, Reservation, User } from "@/types/domain";

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

export async function getUsers(): Promise<ApiResult<readonly User[]>> {
	if (isRealApiEnabled()) {
		return apiGet<readonly User[]>("/api/v1/users");
	}
	return mockApi.getUsers();
}

export async function getUser(id: string): Promise<ApiResult<User>> {
	if (isRealApiEnabled()) {
		return apiGet<User>(`/api/v1/users/${id}`);
	}
	return mockApi.getUser(id);
}

export async function getUsersByDojo(dojoId: string): Promise<ApiResult<readonly User[]>> {
	if (isRealApiEnabled()) {
		return apiGet<readonly User[]>(`/api/v1/dojos/${dojoId}/users`);
	}
	return mockApi.getUsersByDojo(dojoId);
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

export async function getDojos(): Promise<ApiResult<readonly Dojo[]>> {
	if (isRealApiEnabled()) {
		return apiGet<readonly Dojo[]>("/api/v1/dojos");
	}
	return mockApi.getDojos();
}

export async function getDojo(id: string): Promise<ApiResult<Dojo>> {
	if (isRealApiEnabled()) {
		return apiGet<Dojo>(`/api/v1/dojos/${id}`);
	}
	return mockApi.getDojo(id);
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

export async function getPractices(userId?: string): Promise<ApiResult<readonly Practice[]>> {
	if (isRealApiEnabled()) {
		const query = userId ? `?userId=${userId}` : "";
		return apiGet<readonly Practice[]>(`/api/v1/practices${query}`);
	}
	return mockApi.getPractices(userId);
}

export async function getPractice(id: string): Promise<ApiResult<Practice>> {
	if (isRealApiEnabled()) {
		return apiGet<Practice>(`/api/v1/practices/${id}`);
	}
	return mockApi.getPractice(id);
}

export async function createPractice(input: MockCreatePracticeInput): Promise<ApiResult<Practice>> {
	if (isRealApiEnabled()) {
		return apiPost<Practice>("/api/v1/practices", input);
	}
	return mockApi.createPractice(input);
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

export async function getAnalyses(userId?: string): Promise<ApiResult<readonly Analysis[]>> {
	if (isRealApiEnabled()) {
		const query = userId ? `?userId=${userId}` : "";
		return apiGet<readonly Analysis[]>(`/api/v1/analyses${query}`);
	}
	return mockApi.getAnalyses(userId);
}

export async function getAnalysis(id: string): Promise<ApiResult<Analysis>> {
	if (isRealApiEnabled()) {
		return apiGet<Analysis>(`/api/v1/analyses/${id}`);
	}
	return mockApi.getAnalysis(id);
}

export async function getAnalysisByVideo(videoId: string): Promise<ApiResult<Analysis>> {
	if (isRealApiEnabled()) {
		return apiGet<Analysis>(`/api/v1/analyses/by-video/${videoId}`);
	}
	return mockApi.getAnalysisByVideo(videoId);
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

export async function getReservations(dojoId?: string, date?: string): Promise<ApiResult<readonly Reservation[]>> {
	if (isRealApiEnabled()) {
		const params = new URLSearchParams();
		if (dojoId) params.set("dojoId", dojoId);
		if (date) params.set("date", date);
		const query = params.toString() ? `?${params.toString()}` : "";
		return apiGet<readonly Reservation[]>(`/api/v1/reservations${query}`);
	}
	return mockApi.getReservations(dojoId, date);
}

export async function getReservation(id: string): Promise<ApiResult<Reservation>> {
	if (isRealApiEnabled()) {
		return apiGet<Reservation>(`/api/v1/reservations/${id}`);
	}
	return mockApi.getReservation(id);
}

export async function createReservation(input: MockCreateReservationInput): Promise<ApiResult<Reservation>> {
	if (isRealApiEnabled()) {
		return apiPost<Reservation>("/api/v1/reservations", input);
	}
	return mockApi.createReservation(input);
}

export async function deleteReservation(id: string): Promise<ApiResult<{ readonly deleted: true }>> {
	if (isRealApiEnabled()) {
		return apiDelete<{ readonly deleted: true }>(`/api/v1/reservations/${id}`);
	}
	return mockApi.deleteReservation(id);
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

export async function getExamChecklists(userId?: string): Promise<ApiResult<readonly ExamChecklist[]>> {
	if (isRealApiEnabled()) {
		const query = userId ? `?userId=${userId}` : "";
		return apiGet<readonly ExamChecklist[]>(`/api/v1/exam-checklists${query}`);
	}
	return mockApi.getExamChecklists(userId);
}

export async function getExamChecklist(id: string): Promise<ApiResult<ExamChecklist>> {
	if (isRealApiEnabled()) {
		return apiGet<ExamChecklist>(`/api/v1/exam-checklists/${id}`);
	}
	return mockApi.getExamChecklist(id);
}

export async function toggleChecklistItem(checklistId: string, itemId: string): Promise<ApiResult<ExamChecklist>> {
	if (isRealApiEnabled()) {
		return apiPatch<ExamChecklist>(`/api/v1/exam-checklists/${checklistId}/items/${itemId}/toggle`);
	}
	return mockApi.toggleChecklistItem(checklistId, itemId);
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

export type { DashboardSummary } from "@/lib/mock-api";

export async function getDashboardSummary(dojoId: string): Promise<ApiResult<DashboardSummary>> {
	if (isRealApiEnabled()) {
		return apiGet<DashboardSummary>(`/api/v1/dojos/${dojoId}/dashboard`);
	}
	return mockApi.getDashboardSummary(dojoId);
}
