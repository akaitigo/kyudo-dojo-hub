import { useCallback, useEffect, useState } from "react";
import { ExamChecklistView } from "@/components/exam/ExamChecklistView";
import { getExamChecklists, toggleChecklistItem } from "@/lib/api";
import type { ExamChecklist } from "@/types/domain";

/** 現在のモックユーザーID */
const CURRENT_USER_ID = "user-001";

const alertStyle = {
	padding: "1rem",
	marginBottom: "1.5rem",
	backgroundColor: "#fef2f2",
	border: "1px solid #fca5a5",
	borderRadius: "8px",
	color: "#991b1b",
} as const;

export function ExamChecklistPage() {
	const [checklists, setChecklists] = useState<readonly ExamChecklist[]>([]);
	const [error, setError] = useState<string | null>(null);

	const loadChecklists = useCallback(async () => {
		setError(null);
		const result = await getExamChecklists(CURRENT_USER_ID);
		if (result.success) {
			setChecklists(result.data);
			return;
		}
		setError(result.error.message);
	}, []);

	useEffect(() => {
		void loadChecklists();
	}, [loadChecklists]);

	const handleToggle = async (checklistId: string, itemId: string) => {
		const result = await toggleChecklistItem(checklistId, itemId);
		if (result.success) {
			setChecklists((prev) => prev.map((c) => (c.id === result.data.id ? result.data : c)));
			return;
		}
		setError(result.error.message);
	};

	return (
		<div>
			<h1>段位審査チェックリスト</h1>
			{error && (
				<div role="alert" style={alertStyle}>
					チェックリストの読み込みに失敗しました: {error}
				</div>
			)}
			{!error && checklists.length === 0 ? (
				<p>チェックリストがありません</p>
			) : (
				<div style={{ display: "flex", flexDirection: "column", gap: "1.5rem" }}>
					{checklists.map((checklist) => (
						<ExamChecklistView key={checklist.id} checklist={checklist} onToggleItem={handleToggle} />
					))}
				</div>
			)}
		</div>
	);
}
