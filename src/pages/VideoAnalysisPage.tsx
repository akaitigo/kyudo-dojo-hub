import { useCallback, useEffect, useRef, useState } from "react";
import { PhaseTimeline } from "@/components/analysis/PhaseTimeline";
import { ScoreChart } from "@/components/analysis/ScoreChart";
import { VideoUploader } from "@/components/video/VideoUploader";
import { analyzeVideo, createVideo } from "@/lib/api";
import type { Analysis } from "@/types/domain";

/** 現在のユーザーID */
const CURRENT_USER_ID = "user-001";

/** 許可されたMIMEタイプ */
const ALLOWED_MIME_TYPES = ["video/mp4", "video/quicktime", "video/webm"] as const;
type AllowedMimeType = (typeof ALLOWED_MIME_TYPES)[number];

function isAllowedMimeType(value: string): value is AllowedMimeType {
	return (ALLOWED_MIME_TYPES as readonly string[]).includes(value);
}

export function VideoAnalysisPage() {
	const [videoUrl, setVideoUrl] = useState<string | null>(null);
	const [analysis, setAnalysis] = useState<Analysis | null>(null);
	const [analyzing, setAnalyzing] = useState(false);
	const [errorMessage, setErrorMessage] = useState<string | null>(null);
	const [currentTime, setCurrentTime] = useState(0);
	const videoRef = useRef<HTMLVideoElement>(null);

	// Revoke the previous Object URL when videoUrl changes or on unmount
	const prevVideoUrlRef = useRef<string | null>(null);
	useEffect(() => {
		if (prevVideoUrlRef.current && prevVideoUrlRef.current !== videoUrl) {
			URL.revokeObjectURL(prevVideoUrlRef.current);
		}
		prevVideoUrlRef.current = videoUrl;
		return () => {
			if (prevVideoUrlRef.current) {
				URL.revokeObjectURL(prevVideoUrlRef.current);
			}
		};
	}, [videoUrl]);

	const handleUpload = useCallback(async (file: File, objectUrl: string, duration: number) => {
		setVideoUrl(objectUrl);
		setAnalyzing(true);
		setAnalysis(null);
		setErrorMessage(null);

		if (!isAllowedMimeType(file.type)) {
			setErrorMessage("対応していない動画形式です。mp4, mov, webm のいずれかをアップロードしてください。");
			setAnalyzing(false);
			return;
		}

		// 1. 動画メタデータをバックエンドに登録
		const videoResult = await createVideo({
			userId: CURRENT_USER_ID,
			fileName: file.name,
			fileSize: file.size,
			duration,
			mimeType: file.type,
			url: objectUrl,
		});

		if (!videoResult.success) {
			setErrorMessage(videoResult.error.message);
			setAnalyzing(false);
			return;
		}

		// 2. バックエンド経由で射形分析を実行
		const analysisResult = await analyzeVideo({
			videoId: videoResult.data.id,
			userId: CURRENT_USER_ID,
		});

		if (analysisResult.success) {
			setAnalysis(analysisResult.data);
		} else {
			setErrorMessage(analysisResult.error.message);
		}
		setAnalyzing(false);
	}, []);

	const handlePhaseClick = useCallback((startTime: number) => {
		const video = videoRef.current;
		if (video) {
			video.currentTime = startTime;
			void video.play();
		}
	}, []);

	const handleTimeUpdate = () => {
		const video = videoRef.current;
		if (video) {
			setCurrentTime(video.currentTime);
		}
	};

	return (
		<div>
			<h1>射形動画分析</h1>

			{!videoUrl && (
				<section style={{ marginBottom: "2rem" }}>
					<h2>動画アップロード</h2>
					<VideoUploader onUpload={handleUpload} />
				</section>
			)}

			{videoUrl && (
				<div>
					<section style={{ marginBottom: "1.5rem" }}>
						<h2>動画プレーヤー</h2>
						{/* biome-ignore lint/a11y/useMediaCaption: MVP段階ではキャプション不要 */}
						<video
							ref={videoRef}
							src={videoUrl}
							controls
							onTimeUpdate={handleTimeUpdate}
							style={{
								width: "100%",
								maxHeight: "400px",
								borderRadius: "8px",
								backgroundColor: "#000",
							}}
						/>
					</section>

					{analyzing && <p style={{ textAlign: "center", color: "#666", padding: "1rem" }}>射形を分析中...</p>}

					{errorMessage && (
						<div
							role="alert"
							style={{
								padding: "1rem",
								marginBottom: "1.5rem",
								backgroundColor: "#fef2f2",
								border: "1px solid #fca5a5",
								borderRadius: "8px",
								color: "#991b1b",
							}}
						>
							{errorMessage}
						</div>
					)}

					{analysis && (
						<>
							<section style={{ marginBottom: "1.5rem" }}>
								<PhaseTimeline phases={analysis.phases} onPhaseClick={handlePhaseClick} currentTime={currentTime} />
							</section>

							<section style={{ marginBottom: "1.5rem" }}>
								<ScoreChart scores={analysis.scores} />
							</section>

							<section
								style={{
									marginBottom: "1.5rem",
									padding: "1rem",
									backgroundColor: "#f8f9fa",
									borderRadius: "8px",
									border: "1px solid #e0e0e0",
								}}
							>
								<h3 style={{ margin: "0 0 0.5rem" }}>総合スコア: {analysis.overallScore}/100</h3>
								<p style={{ margin: 0, color: "#555" }}>{analysis.feedback}</p>
							</section>
						</>
					)}

					<button
						type="button"
						onClick={() => {
							setVideoUrl(null);
							setAnalysis(null);
							setErrorMessage(null);
						}}
						style={{
							padding: "0.5rem 1rem",
							backgroundColor: "#999",
							color: "#fff",
							border: "none",
							borderRadius: "4px",
							cursor: "pointer",
						}}
					>
						別の動画をアップロード
					</button>
				</div>
			)}
		</div>
	);
}
