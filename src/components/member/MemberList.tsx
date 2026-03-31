import { useState } from "react";
import type { User } from "@/types/domain";
import { DAN_LABELS, DAN_RANKS, SHOGO_LABELS } from "@/types/domain";

interface MemberListProps {
	readonly members: readonly User[];
}

export function MemberList({ members }: MemberListProps) {
	const [sortBy, setSortBy] = useState<"name" | "dan" | "joinedAt">("dan");

	const sorted = [...members].sort((a, b) => {
		if (sortBy === "name") return a.name.localeCompare(b.name);
		if (sortBy === "joinedAt") return a.joinedAt.localeCompare(b.joinedAt);
		// Sort by dan rank
		const aIdx = a.dan ? DAN_RANKS.indexOf(a.dan) : -1;
		const bIdx = b.dan ? DAN_RANKS.indexOf(b.dan) : -1;
		return bIdx - aIdx;
	});

	const getDanDisplay = (user: User): string => {
		const parts: string[] = [];
		if (user.dan) parts.push(DAN_LABELS[user.dan]);
		if (user.shogo) parts.push(SHOGO_LABELS[user.shogo]);
		return parts.length > 0 ? parts.join(" / ") : "未取得";
	};

	if (members.length === 0) {
		return <p>会員がいません</p>;
	}

	return (
		<div>
			<div
				style={{
					marginBottom: "1rem",
					display: "flex",
					gap: "0.5rem",
					alignItems: "center",
				}}
			>
				<span style={{ fontWeight: "bold" }}>ソート:</span>
				{(["dan", "name", "joinedAt"] as const).map((key) => (
					<button
						key={key}
						type="button"
						onClick={() => setSortBy(key)}
						style={{
							padding: "0.25rem 0.75rem",
							backgroundColor: sortBy === key ? "#1a1a2e" : "#e0e0e0",
							color: sortBy === key ? "#fff" : "#333",
							border: "none",
							borderRadius: "4px",
							cursor: "pointer",
							fontSize: "0.85rem",
						}}
					>
						{key === "dan" ? "段位" : key === "name" ? "名前" : "入会日"}
					</button>
				))}
			</div>

			<table style={{ borderCollapse: "collapse", width: "100%" }}>
				<thead>
					<tr>
						{["名前", "段位・称号", "入会日", "連絡先"].map((h) => (
							<th
								key={h}
								style={{
									border: "1px solid #e0e0e0",
									padding: "0.5rem",
									backgroundColor: "#f5f5f5",
									textAlign: "left",
								}}
							>
								{h}
							</th>
						))}
					</tr>
				</thead>
				<tbody>
					{sorted.map((member) => (
						<tr key={member.id}>
							<td style={{ border: "1px solid #e0e0e0", padding: "0.5rem" }}>{member.name}</td>
							<td style={{ border: "1px solid #e0e0e0", padding: "0.5rem" }}>{getDanDisplay(member)}</td>
							<td style={{ border: "1px solid #e0e0e0", padding: "0.5rem" }}>{member.joinedAt}</td>
							<td style={{ border: "1px solid #e0e0e0", padding: "0.5rem" }}>{member.email}</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
}
