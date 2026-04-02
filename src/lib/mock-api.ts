/**
 * モック API 層
 *
 * ADR-001: MVP ではフロントエンドに集中し、バックエンドはモック。
 * 本番 API と同じインターフェースを持ち、将来の差し替えを容易にする。
 * 意図的な遅延（50-200ms）を追加し、非同期UIの動作確認を可能にする。
 */
import { getLocalDateString } from "@/lib/date-utils";
import type { Analysis, ApiResult, Dojo, ExamChecklist, Practice, Reservation, User, Video } from "@/types/domain";
import {
	MOCK_ANALYSES,
	MOCK_DOJOS,
	MOCK_EXAM_CHECKLISTS,
	MOCK_PRACTICES,
	MOCK_RESERVATIONS,
	MOCK_USERS,
	MOCK_VIDEOS,
} from "./mock-data";

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

/** 意図的な遅延（50-200ms）を追加 */
function simulateLatency(): Promise<void> {
	const delay = Math.floor(Math.random() * 150) + 50;
	return new Promise((resolve) => {
		setTimeout(resolve, delay);
	});
}

function success<T>(data: T): ApiResult<T> {
	return { success: true, data };
}

function notFound(resource: string): ApiResult<never> {
	return {
		success: false,
		error: { code: "NOT_FOUND", message: `${resource} が見つかりません` },
	};
}

function validationError(message: string): ApiResult<never> {
	return {
		success: false,
		error: { code: "VALIDATION_ERROR", message },
	};
}

/** 簡易 UUID 生成 */
function generateId(): string {
	return `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

// ---------------------------------------------------------------------------
// In-memory state (mutable copies)
// ---------------------------------------------------------------------------

const users = [...MOCK_USERS];
let practices = [...MOCK_PRACTICES];
let videos = [...MOCK_VIDEOS];
const analyses = [...MOCK_ANALYSES];
let reservations = [...MOCK_RESERVATIONS];
let examChecklists = [...MOCK_EXAM_CHECKLISTS];
const dojos = [...MOCK_DOJOS];

// ---------------------------------------------------------------------------
// Users API
// ---------------------------------------------------------------------------

export async function getUsers(): Promise<ApiResult<readonly User[]>> {
	await simulateLatency();
	return success(users);
}

export async function getUser(id: string): Promise<ApiResult<User>> {
	await simulateLatency();
	const user = users.find((u) => u.id === id);
	return user ? success(user) : notFound("ユーザー");
}

export async function getUsersByDojo(dojoId: string): Promise<ApiResult<readonly User[]>> {
	await simulateLatency();
	return success(users.filter((u) => u.dojoId === dojoId));
}

// ---------------------------------------------------------------------------
// Dojos API
// ---------------------------------------------------------------------------

export async function getDojos(): Promise<ApiResult<readonly Dojo[]>> {
	await simulateLatency();
	return success(dojos);
}

export async function getDojo(id: string): Promise<ApiResult<Dojo>> {
	await simulateLatency();
	const dojo = dojos.find((d) => d.id === id);
	return dojo ? success(dojo) : notFound("道場");
}

// ---------------------------------------------------------------------------
// Practices API
// ---------------------------------------------------------------------------

export async function getPractices(userId?: string): Promise<ApiResult<readonly Practice[]>> {
	await simulateLatency();
	const filtered = userId ? practices.filter((p) => p.userId === userId) : practices;
	const sorted = [...filtered].sort((a, b) => b.date.localeCompare(a.date));
	return success(sorted);
}

export async function getPractice(id: string): Promise<ApiResult<Practice>> {
	await simulateLatency();
	const practice = practices.find((p) => p.id === id);
	return practice ? success(practice) : notFound("稽古記録");
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
	await simulateLatency();

	// Validation
	if (input.hitRate < 0 || input.hitRate > 100) {
		return validationError("的中率は 0〜100 の範囲で入力してください");
	}
	if (input.arrowCount < 1 || input.arrowCount > 1000) {
		return validationError("矢数は 1〜1000 の範囲で入力してください");
	}
	if (input.notes.length > 5000) {
		return validationError("気づきは 5,000 文字以内で入力してください");
	}
	if (input.instructorComment.length > 5000) {
		return validationError("師範コメントは 5,000 文字以内で入力してください");
	}

	const now = new Date().toISOString();
	const base = {
		id: `practice-${generateId()}`,
		userId: input.userId,
		date: input.date,
		hitRate: input.hitRate,
		arrowCount: input.arrowCount,
		notes: input.notes,
		instructorComment: input.instructorComment,
		createdAt: now,
		updatedAt: now,
	};
	const practice: Practice = input.dojoId !== undefined ? { ...base, dojoId: input.dojoId } : base;
	practices = [practice, ...practices];
	return success(practice);
}

// ---------------------------------------------------------------------------
// Videos API
// ---------------------------------------------------------------------------

export async function getVideos(userId?: string): Promise<ApiResult<readonly Video[]>> {
	await simulateLatency();
	const filtered = userId ? videos.filter((v) => v.userId === userId) : videos;
	return success(filtered);
}

export async function getVideo(id: string): Promise<ApiResult<Video>> {
	await simulateLatency();
	const video = videos.find((v) => v.id === id);
	return video ? success(video) : notFound("動画");
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
	await simulateLatency();

	if (input.fileSize > 500 * 1024 * 1024) {
		return validationError("ファイルサイズは 500MB 以下にしてください");
	}
	if (input.duration > 300) {
		return validationError("動画長は 5 分以下にしてください");
	}

	const now = new Date().toISOString();
	const video: Video = {
		id: `video-${generateId()}`,
		...input,
		status: "completed",
		createdAt: now,
		updatedAt: now,
	};
	videos = [video, ...videos];
	return success(video);
}

// ---------------------------------------------------------------------------
// Analyses API
// ---------------------------------------------------------------------------

export async function getAnalyses(userId?: string): Promise<ApiResult<readonly Analysis[]>> {
	await simulateLatency();
	const filtered = userId ? analyses.filter((a) => a.userId === userId) : analyses;
	return success(filtered);
}

export async function getAnalysis(id: string): Promise<ApiResult<Analysis>> {
	await simulateLatency();
	const analysis = analyses.find((a) => a.id === id);
	return analysis ? success(analysis) : notFound("分析結果");
}

export async function getAnalysisByVideo(videoId: string): Promise<ApiResult<Analysis>> {
	await simulateLatency();
	const analysis = analyses.find((a) => a.videoId === videoId);
	return analysis ? success(analysis) : notFound("分析結果");
}

// ---------------------------------------------------------------------------
// Reservations API
// ---------------------------------------------------------------------------

export async function getReservations(dojoId?: string, date?: string): Promise<ApiResult<readonly Reservation[]>> {
	await simulateLatency();
	let filtered = reservations;
	if (dojoId) {
		filtered = filtered.filter((r) => r.dojoId === dojoId);
	}
	if (date) {
		filtered = filtered.filter((r) => r.date === date);
	}
	return success(filtered);
}

export async function getReservation(id: string): Promise<ApiResult<Reservation>> {
	await simulateLatency();
	const reservation = reservations.find((r) => r.id === id);
	return reservation ? success(reservation) : notFound("予約");
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
	await simulateLatency();

	// 重複チェック（時間帯の重なり判定: start_a < end_b AND start_b < end_a）
	const conflict = reservations.find(
		(r) =>
			r.dojoId === input.dojoId &&
			r.laneNumber === input.laneNumber &&
			r.date === input.date &&
			r.startTime < input.endTime &&
			input.startTime < r.endTime,
	);
	if (conflict) {
		return validationError("同一的場・同一時間帯に既に予約があります");
	}

	const now = new Date().toISOString();
	const reservation: Reservation = {
		id: `res-${generateId()}`,
		...input,
		createdAt: now,
		updatedAt: now,
	};
	reservations = [reservation, ...reservations];
	return success(reservation);
}

export async function deleteReservation(id: string): Promise<ApiResult<{ readonly deleted: true }>> {
	await simulateLatency();
	const index = reservations.findIndex((r) => r.id === id);
	if (index === -1) {
		return notFound("予約");
	}
	reservations = [...reservations.slice(0, index), ...reservations.slice(index + 1)];
	return success({ deleted: true });
}

// ---------------------------------------------------------------------------
// Exam Checklists API
// ---------------------------------------------------------------------------

export async function getExamChecklists(userId?: string): Promise<ApiResult<readonly ExamChecklist[]>> {
	await simulateLatency();
	const filtered = userId ? examChecklists.filter((c) => c.userId === userId) : examChecklists;
	return success(filtered);
}

export async function getExamChecklist(id: string): Promise<ApiResult<ExamChecklist>> {
	await simulateLatency();
	const checklist = examChecklists.find((c) => c.id === id);
	return checklist ? success(checklist) : notFound("審査チェックリスト");
}

export async function toggleChecklistItem(checklistId: string, itemId: string): Promise<ApiResult<ExamChecklist>> {
	await simulateLatency();
	const checklistIndex = examChecklists.findIndex((c) => c.id === checklistId);
	if (checklistIndex === -1) {
		return notFound("審査チェックリスト");
	}

	const checklist = examChecklists[checklistIndex];
	if (!checklist) {
		return notFound("審査チェックリスト");
	}

	const updatedItems = checklist.items.map((item) =>
		item.id === itemId ? { ...item, completed: !item.completed } : item,
	);

	const completedCount = updatedItems.filter((item) => item.completed).length;
	const progressRate = Math.round((completedCount / updatedItems.length) * 100);

	const updatedChecklist: ExamChecklist = {
		...checklist,
		items: updatedItems,
		progressRate,
		updatedAt: new Date().toISOString(),
	};

	examChecklists = [
		...examChecklists.slice(0, checklistIndex),
		updatedChecklist,
		...examChecklists.slice(checklistIndex + 1),
	];
	return success(updatedChecklist);
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
	await simulateLatency();
	const today = getLocalDateString();
	const todayReservations = reservations.filter((r) => r.dojoId === dojoId && r.date === today);
	const memberCount = users.filter((u) => u.dojoId === dojoId).length;

	return success({
		todayReservationCount: todayReservations.length,
		totalMemberCount: memberCount,
		todayReservations,
	});
}
