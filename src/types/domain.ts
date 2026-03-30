/**
 * 弓道ドメインのコアデータモデル型定義
 *
 * 射法八節: 足踏み→胴造り→弓構え→打起し→引分け→会→離れ→残心
 * 段位体系: 初段〜十段 + 錬士・教士・範士
 */

// ---------------------------------------------------------------------------
// 射法八節（Shaho Hassetsu）
// ---------------------------------------------------------------------------

/** 射法八節の各フェーズ */
export const HASSETSU_PHASES = [
	"ashibumi",
	"dozukuri",
	"yugamae",
	"uchiokoshi",
	"hikiwake",
	"kai",
	"hanare",
	"zanshin",
] as const;

export type HassetsuPhase = (typeof HASSETSU_PHASES)[number];

/** 八節フェーズの日本語ラベル */
export const HASSETSU_LABELS: Record<HassetsuPhase, string> = {
	ashibumi: "足踏み",
	dozukuri: "胴造り",
	yugamae: "弓構え",
	uchiokoshi: "打起し",
	hikiwake: "引分け",
	kai: "会",
	hanare: "離れ",
	zanshin: "残心",
};

/** 八節スコア: 各段階 0〜100 の整数 */
export type HassetsuScores = Record<HassetsuPhase, number>;

/** フェーズ分割結果（動画タイムライン） */
export interface PhaseSegment {
	readonly phase: HassetsuPhase;
	/** 開始時刻（秒） */
	readonly startTime: number;
	/** 終了時刻（秒） */
	readonly endTime: number;
}

// ---------------------------------------------------------------------------
// 段位（Dan / Shogo）
// ---------------------------------------------------------------------------

export const DAN_RANKS = [
	"shodan",
	"nidan",
	"sandan",
	"yondan",
	"godan",
	"rokudan",
	"nanadan",
	"hachidan",
	"kudan",
	"judan",
] as const;

export type DanRank = (typeof DAN_RANKS)[number];

export const SHOGO_TITLES = ["renshi", "kyoshi", "hanshi"] as const;

export type ShogoTitle = (typeof SHOGO_TITLES)[number];

/** 段位の日本語ラベル */
export const DAN_LABELS: Record<DanRank, string> = {
	shodan: "初段",
	nidan: "二段",
	sandan: "三段",
	yondan: "四段",
	godan: "五段",
	rokudan: "六段",
	nanadan: "七段",
	hachidan: "八段",
	kudan: "九段",
	judan: "十段",
};

/** 称号の日本語ラベル */
export const SHOGO_LABELS: Record<ShogoTitle, string> = {
	renshi: "錬士",
	kyoshi: "教士",
	hanshi: "範士",
};

// ---------------------------------------------------------------------------
// ユーザー
// ---------------------------------------------------------------------------

export type UserRole = "practitioner" | "manager" | "admin";

export interface User {
	readonly id: string;
	readonly name: string;
	readonly email: string;
	readonly role: UserRole;
	readonly dan?: DanRank;
	readonly shogo?: ShogoTitle;
	readonly dojoId?: string;
	readonly joinedAt: string;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// 道場
// ---------------------------------------------------------------------------

export interface Dojo {
	readonly id: string;
	readonly name: string;
	readonly address: string;
	/** 的場の数 */
	readonly targetLanes: number;
	/** 営業開始時刻（HH:mm） */
	readonly openTime: string;
	/** 営業終了時刻（HH:mm） */
	readonly closeTime: string;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// 稽古記録
// ---------------------------------------------------------------------------

export interface Practice {
	readonly id: string;
	readonly userId: string;
	readonly dojoId?: string;
	readonly date: string;
	/** 的中率: 0〜100 の整数 */
	readonly hitRate: number;
	/** 矢数: 1〜1000 の整数 */
	readonly arrowCount: number;
	/** 気づき: 最大 5,000 文字 */
	readonly notes: string;
	/** 師範コメント: 最大 5,000 文字 */
	readonly instructorComment: string;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// 動画
// ---------------------------------------------------------------------------

export type VideoStatus = "uploading" | "processing" | "completed" | "failed";

export interface Video {
	readonly id: string;
	readonly userId: string;
	readonly practiceId?: string;
	readonly fileName: string;
	/** ファイルサイズ（バイト） */
	readonly fileSize: number;
	/** 動画長（秒） */
	readonly duration: number;
	readonly mimeType: "video/mp4" | "video/quicktime" | "video/webm";
	readonly status: VideoStatus;
	/** Object URL or Cloud Storage URL */
	readonly url: string;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// 分析結果
// ---------------------------------------------------------------------------

export interface Analysis {
	readonly id: string;
	readonly videoId: string;
	readonly userId: string;
	readonly scores: HassetsuScores;
	readonly phases: readonly PhaseSegment[];
	/** 総合スコア: 0〜100 */
	readonly overallScore: number;
	/** 分析コメント */
	readonly feedback: string;
	readonly createdAt: string;
}

// ---------------------------------------------------------------------------
// 的場予約
// ---------------------------------------------------------------------------

export interface Reservation {
	readonly id: string;
	readonly dojoId: string;
	readonly userId: string;
	/** 的場レーン番号（1始まり） */
	readonly laneNumber: number;
	/** 予約日 (YYYY-MM-DD) */
	readonly date: string;
	/** 開始時刻 (HH:mm) */
	readonly startTime: string;
	/** 終了時刻 (HH:mm) */
	readonly endTime: string;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// 段位審査チェックリスト
// ---------------------------------------------------------------------------

export interface ExamChecklistItem {
	readonly id: string;
	readonly category: string;
	readonly description: string;
	readonly completed: boolean;
}

export interface ExamChecklist {
	readonly id: string;
	readonly userId: string;
	readonly targetDan: DanRank;
	readonly items: readonly ExamChecklistItem[];
	/** 進捗率: 0〜100 */
	readonly progressRate: number;
	readonly createdAt: string;
	readonly updatedAt: string;
}

// ---------------------------------------------------------------------------
// API レスポンス型
// ---------------------------------------------------------------------------

export interface ApiResponse<T> {
	readonly data: T;
	readonly success: true;
}

export interface ApiError {
	readonly success: false;
	readonly error: {
		readonly code: string;
		readonly message: string;
	};
}

export type ApiResult<T> = ApiResponse<T> | ApiError;
