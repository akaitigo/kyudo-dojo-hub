import { useCallback, useEffect, useRef, useState } from "react";
import { validateVideoDuration, validateVideoFile } from "@/lib/video-validation";

interface VideoUploaderProps {
	readonly onUpload: (file: File, objectUrl: string, duration: number) => void | Promise<void>;
}

export function VideoUploader({ onUpload }: VideoUploaderProps) {
	const [isDragging, setIsDragging] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [progress, setProgress] = useState<number | null>(null);
	const fileInputRef = useRef<HTMLInputElement>(null);
	const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
	const completionTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
	const resetTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
	const loadIdRef = useRef(0);

	useEffect(() => {
		return () => {
			if (intervalRef.current !== null) {
				clearInterval(intervalRef.current);
			}
			if (completionTimerRef.current !== null) {
				clearTimeout(completionTimerRef.current);
			}
			if (resetTimerRef.current !== null) {
				clearTimeout(resetTimerRef.current);
			}
		};
	}, []);

	const startUpload = useCallback(
		(file: File, duration: number) => {
			// Clear any previous timers
			if (intervalRef.current !== null) {
				clearInterval(intervalRef.current);
			}
			if (completionTimerRef.current !== null) {
				clearTimeout(completionTimerRef.current);
			}
			if (resetTimerRef.current !== null) {
				clearTimeout(resetTimerRef.current);
			}

			// Simulate upload progress
			setProgress(0);
			intervalRef.current = setInterval(() => {
				setProgress((prev) => {
					if (prev === null || prev >= 100) {
						if (intervalRef.current !== null) {
							clearInterval(intervalRef.current);
							intervalRef.current = null;
						}
						return 100;
					}
					return prev + 10;
				});
			}, 100);

			completionTimerRef.current = setTimeout(() => {
				if (intervalRef.current !== null) {
					clearInterval(intervalRef.current);
					intervalRef.current = null;
				}
				setProgress(100);
				const objectUrl = URL.createObjectURL(file);
				onUpload(file, objectUrl, duration);
				resetTimerRef.current = setTimeout(() => setProgress(null), 500);
			}, 1100);
		},
		[onUpload],
	);

	const processFile = useCallback(
		(file: File) => {
			setError(null);
			const validation = validateVideoFile(file);
			if (!validation.valid) {
				setError(validation.error ?? "不正なファイルです");
				return;
			}

			// Invalidate any in-flight metadata probe from a previous file selection
			const currentLoadId = ++loadIdRef.current;

			// Check video duration via metadata
			const video = document.createElement("video");
			video.preload = "metadata";
			const tempUrl = URL.createObjectURL(file);
			video.onloadedmetadata = () => {
				URL.revokeObjectURL(tempUrl);
				// Stale callback — a newer file was selected
				if (currentLoadId !== loadIdRef.current) return;
				const durationCheck = validateVideoDuration(video.duration);
				if (!durationCheck.valid) {
					setError(durationCheck.error ?? "不正な動画です");
					return;
				}
				startUpload(file, video.duration);
			};
			video.onerror = () => {
				URL.revokeObjectURL(tempUrl);
				if (currentLoadId !== loadIdRef.current) return;
				setError("動画メタデータの読み込みに失敗しました");
			};
			video.src = tempUrl;
		},
		[startUpload],
	);

	const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
		const file = e.target.files?.[0];
		if (file) {
			processFile(file);
		}
	};

	const handleDrop = (e: React.DragEvent) => {
		e.preventDefault();
		setIsDragging(false);
		const file = e.dataTransfer.files[0];
		if (file) {
			processFile(file);
		}
	};

	const handleDragOver = (e: React.DragEvent) => {
		e.preventDefault();
		setIsDragging(true);
	};

	const handleDragLeave = () => setIsDragging(false);

	return (
		<div>
			<button
				type="button"
				onDrop={handleDrop}
				onDragOver={handleDragOver}
				onDragLeave={handleDragLeave}
				onClick={() => fileInputRef.current?.click()}
				style={{
					border: `2px dashed ${isDragging ? "#1a1a2e" : "#ccc"}`,
					borderRadius: "8px",
					padding: "2rem",
					textAlign: "center",
					cursor: "pointer",
					backgroundColor: isDragging ? "#f0f0ff" : "#fafafa",
					transition: "all 0.2s ease",
					width: "100%",
					font: "inherit",
				}}
			>
				<p style={{ fontSize: "1.1rem", marginBottom: "0.5rem" }}>
					動画ファイルをドラッグ＆ドロップ、またはクリックして選択
				</p>
				<p style={{ color: "#999", fontSize: "0.85rem" }}>mp4, mov, webm / 最大 500MB / 最大 5分</p>
				<input
					ref={fileInputRef}
					type="file"
					accept="video/mp4,video/quicktime,video/webm"
					onChange={handleFileSelect}
					style={{ display: "none" }}
				/>
			</button>

			{progress !== null && (
				<div
					style={{
						marginTop: "1rem",
						backgroundColor: "#f0f0f0",
						borderRadius: "4px",
						overflow: "hidden",
						height: "24px",
					}}
				>
					<div
						style={{
							width: `${String(progress)}%`,
							backgroundColor: "#1a1a2e",
							height: "100%",
							display: "flex",
							alignItems: "center",
							justifyContent: "center",
							color: "#fff",
							fontSize: "0.8rem",
							transition: "width 0.1s ease",
						}}
					>
						{progress}%
					</div>
				</div>
			)}

			{error && <p style={{ color: "#d32f2f", marginTop: "0.5rem" }}>{error}</p>}
		</div>
	);
}
