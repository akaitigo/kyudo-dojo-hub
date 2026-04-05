import { useCallback, useEffect, useState } from "react";
import { ExamChecklistView } from "@/components/exam/ExamChecklistView";
import { getExamChecklists, toggleChecklistItem } from "@/lib/api";
import type { ExamChecklist } from "@/types/domain";

/** 現在のモックユーザーID */
const CURRENT_USER_ID = "user-001";

export function ExamChecklistPage() {
	const [checklists, setChecklists] = useState<readonly ExamChecklist[]>([]);

	const loadChecklists = useCallback(async () => {
		const result = await getExamChecklists(CURRENT_USER_ID);
		if (result.success) {
			setChecklists(result.data);
		}
	}, []);

	useEffect(() => {
		void loadChecklists();
	}, [loadChecklists]);

	const handleToggle = async (checklistId: string, itemId: string) => {
		const result = await toggleChecklistItem(checklistId, itemId);
		if (result.success) {
			setChecklists((prev) => prev.map((c) => (c.id === result.data.id ? result.data : c)));
		}
	};

	return (
		<div>
			<h1>段位審査チェックリスト</h1>
			{checklists.length === 0 ? (
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
