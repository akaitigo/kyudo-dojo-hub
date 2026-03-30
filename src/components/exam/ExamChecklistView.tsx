import type { ExamChecklist } from "@/types/domain";
import { DAN_LABELS } from "@/types/domain";

interface ExamChecklistViewProps {
	readonly checklist: ExamChecklist;
	readonly onToggleItem: (checklistId: string, itemId: string) => void;
}

export function ExamChecklistView({ checklist, onToggleItem }: ExamChecklistViewProps) {
	const categories = [...new Set(checklist.items.map((item) => item.category))];

	return (
		<div
			style={{
				border: "1px solid #e0e0e0",
				borderRadius: "8px",
				padding: "1.5rem",
			}}
		>
			<h3 style={{ margin: "0 0 0.5rem" }}>{DAN_LABELS[checklist.targetDan]} 審査チェックリスト</h3>

			<div
				style={{
					marginBottom: "1rem",
					backgroundColor: "#f0f0f0",
					borderRadius: "4px",
					overflow: "hidden",
					height: "24px",
				}}
			>
				<div
					style={{
						width: `${String(checklist.progressRate)}%`,
						backgroundColor: checklist.progressRate >= 80 ? "#2e7d32" : "#ff9800",
						height: "100%",
						display: "flex",
						alignItems: "center",
						justifyContent: "center",
						color: "#fff",
						fontSize: "0.8rem",
						fontWeight: "bold",
						transition: "width 0.3s ease",
					}}
				>
					{checklist.progressRate}%
				</div>
			</div>

			{categories.map((category) => (
				<div key={category} style={{ marginBottom: "1rem" }}>
					<h4 style={{ margin: "0 0 0.5rem", color: "#555" }}>{category}</h4>
					{checklist.items
						.filter((item) => item.category === category)
						.map((item) => (
							<label
								key={item.id}
								style={{
									display: "flex",
									alignItems: "flex-start",
									gap: "0.5rem",
									marginBottom: "0.5rem",
									cursor: "pointer",
								}}
							>
								<input
									type="checkbox"
									checked={item.completed}
									onChange={() => onToggleItem(checklist.id, item.id)}
									style={{ marginTop: "0.2rem" }}
								/>
								<span
									style={{
										textDecoration: item.completed ? "line-through" : "none",
										color: item.completed ? "#999" : "#333",
									}}
								>
									{item.description}
								</span>
							</label>
						))}
				</div>
			))}
		</div>
	);
}
