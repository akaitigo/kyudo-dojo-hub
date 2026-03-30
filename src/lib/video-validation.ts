/** 動画ファイルバリデーション */

const ALLOWED_MIME_TYPES = new Set(["video/mp4", "video/quicktime", "video/webm"]);
const MAX_FILE_SIZE = 500 * 1024 * 1024; // 500MB
const MAX_DURATION_SECONDS = 300; // 5分

export interface VideoValidationResult {
	readonly valid: boolean;
	readonly error?: string;
}

export function validateVideoFile(file: File): VideoValidationResult {
	if (!ALLOWED_MIME_TYPES.has(file.type)) {
		return {
			valid: false,
			error: "mp4, mov, webm 形式のファイルのみアップロードできます",
		};
	}

	if (file.size > MAX_FILE_SIZE) {
		return {
			valid: false,
			error: "ファイルサイズは 500MB 以下にしてください",
		};
	}

	return { valid: true };
}

export function validateVideoDuration(duration: number): VideoValidationResult {
	if (duration > MAX_DURATION_SECONDS) {
		return {
			valid: false,
			error: "動画長は 5 分以下にしてください",
		};
	}

	return { valid: true };
}
