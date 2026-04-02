import { useCallback, useEffect, useRef, useState } from "react";
import { PhaseTimeline } from "@/components/analysis/PhaseTimeline";
import { ScoreChart } from "@/components/analysis/ScoreChart";
import { VideoUploader } from "@/components/video/VideoUploader";
import { MOCK_ANALYSES } from "@/lib/mock-data";
import type { Analysis } from "@/types/domain";

export function VideoAnalysisPage() {
	const [videoUrl, setVideoUrl] = useState<string | null>(null);
	const [analysis, setAnalysis] = useState<Analysis | null>(null);
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

	const handleUpload = useCallback((file: File, objectUrl: string) => {
		setVideoUrl(objectUrl);
		// MVP: ファイル名ベースでモック分析結果を選択（異なる動画に異なる結果）
		if (MOCK_ANALYSES.length > 0) {
			let hash = 0;
			for (let i = 0; i < file.name.length; i++) {
				hash = (hash * 31 + file.name.charCodeAt(i)) | 0;
			}
			const index = Math.abs(hash) % MOCK_ANALYSES.length;
			const mockAnalysis = MOCK_ANALYSES[index];
			if (mockAnalysis) {
				setAnalysis(mockAnalysis);
			}
		}
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
