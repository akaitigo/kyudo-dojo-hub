import { MOCK_USERS } from "@/lib/mock-data";
import type { Practice } from "@/types/domain";
import { DAN_LABELS } from "@/types/domain";

interface PracticeListProps {
	readonly practices: readonly Practice[];
}

function getUserName(userId: string): string {
	const user = MOCK_USERS.find((u) => u.id === userId);
	if (!user) return userId;
	const danLabel = user.dan ? ` (${DAN_LABELS[user.dan]})` : "";
	return `${user.name}${danLabel}`;
}

export function PracticeList({ practices }: PracticeListProps) {
	if (practices.length === 0) {
		return <p>稽古記録がありません</p>;
	}

	return (
		<div style={{ display: "flex", flexDirection: "column", gap: "0.75rem" }}>
			{practices.map((practice) => (
				<div
					key={practice.id}
					style={{
						border: "1px solid #e0e0e0",
						borderRadius: "8px",
						padding: "1rem",
						backgroundColor: "#fafafa",
					}}
				>
					<div
						style={{
							display: "flex",
							justifyContent: "space-between",
							marginBottom: "0.5rem",
						}}
					>
						<strong>{practice.date}</strong>
						<span
							style={{
								color: practice.hitRate >= 60 ? "#2e7d32" : "#d32f2f",
								fontWeight: "bold",
							}}
						>
							的中率: {practice.hitRate}%
						</span>
					</div>
					<div
						style={{
							fontSize: "0.9rem",
							color: "#666",
							marginBottom: "0.25rem",
						}}
					>
						{getUserName(practice.userId)} / 矢数: {practice.arrowCount}本
					</div>
					{practice.notes && <p style={{ margin: "0.5rem 0 0", fontSize: "0.9rem" }}>{practice.notes}</p>}
					{practice.instructorComment && (
						<p
							style={{
								margin: "0.5rem 0 0",
								fontSize: "0.9rem",
								color: "#1a1a2e",
								fontStyle: "italic",
							}}
						>
							師範: {practice.instructorComment}
						</p>
					)}
				</div>
			))}
		</div>
	);
}
