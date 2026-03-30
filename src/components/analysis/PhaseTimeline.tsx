import { HASSETSU_LABELS } from "@/types/domain";
import type { PhaseSegment } from "@/types/domain";

interface PhaseTimelineProps {
	readonly phases: readonly PhaseSegment[];
	readonly onPhaseClick: (startTime: number) => void;
	readonly currentTime?: number;
}

const PHASE_COLORS: Record<string, string> = {
	ashibumi: "#e74c3c",
	dozukuri: "#e67e22",
	yugamae: "#f1c40f",
	uchiokoshi: "#2ecc71",
	hikiwake: "#1abc9c",
	kai: "#3498db",
	hanare: "#9b59b6",
	zanshin: "#34495e",
};

export function PhaseTimeline({ phases, onPhaseClick, currentTime = 0 }: PhaseTimelineProps) {
	if (phases.length === 0) return null;

	const lastPhase = phases[phases.length - 1];
	const totalDuration = lastPhase?.endTime ?? 0;
	if (totalDuration === 0) return null;

	return (
		<div style={{ marginBottom: "1rem" }}>
			<h3 style={{ marginBottom: "0.5rem" }}>八節フェーズタイムライン</h3>
			<div
				style={{
					display: "flex",
					borderRadius: "4px",
					overflow: "hidden",
					height: "40px",
					marginBottom: "0.5rem",
				}}
			>
				{phases.map((phase) => {
					const widthPercent = ((phase.endTime - phase.startTime) / totalDuration) * 100;
					const isActive = currentTime >= phase.startTime && currentTime < phase.endTime;
					return (
						<button
							key={phase.phase}
							type="button"
							onClick={() => onPhaseClick(phase.startTime)}
							style={{
								width: `${String(widthPercent)}%`,
								backgroundColor: PHASE_COLORS[phase.phase] ?? "#999",
								border: "none",
								cursor: "pointer",
								opacity: isActive ? 1 : 0.7,
								display: "flex",
								alignItems: "center",
								justifyContent: "center",
								color: "#fff",
								fontSize: "0.7rem",
								fontWeight: isActive ? "bold" : "normal",
								outline: isActive ? "2px solid #fff" : "none",
								outlineOffset: "-2px",
								transition: "opacity 0.2s",
							}}
							title={`${HASSETSU_LABELS[phase.phase]} (${phase.startTime.toFixed(1)}s - ${phase.endTime.toFixed(1)}s)`}
						>
							{HASSETSU_LABELS[phase.phase]}
						</button>
					);
				})}
			</div>
			<div
				style={{
					display: "flex",
					flexWrap: "wrap",
					gap: "0.5rem",
					fontSize: "0.8rem",
				}}
			>
				{phases.map((phase) => (
					<span key={phase.phase} style={{ display: "flex", alignItems: "center", gap: "0.25rem" }}>
						<span
							style={{
								width: "12px",
								height: "12px",
								backgroundColor: PHASE_COLORS[phase.phase] ?? "#999",
								borderRadius: "2px",
								display: "inline-block",
							}}
						/>
						{HASSETSU_LABELS[phase.phase]}
					</span>
				))}
			</div>
		</div>
	);
}
