import { useCallback, useEffect, useState } from "react";
import { HitRateChart } from "@/components/practice/HitRateChart";
import { PracticeForm } from "@/components/practice/PracticeForm";
import { PracticeList } from "@/components/practice/PracticeList";
import { createPractice, getPractices, getUsers } from "@/lib/api";
import type { PracticeFormValues } from "@/lib/validation";
import type { Practice, User } from "@/types/domain";

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

export function PracticesPage() {
	const [practices, setPractices] = useState<readonly Practice[]>([]);
	const [users, setUsers] = useState<readonly User[]>([]);
	const [isSubmitting, setIsSubmitting] = useState(false);
	const [showForm, setShowForm] = useState(false);
	const [loadError, setLoadError] = useState<string | null>(null);
	const [submitError, setSubmitError] = useState<string | null>(null);

	const loadData = useCallback(async () => {
		setLoadError(null);
		const [practicesResult, usersResult] = await Promise.all([getPractices(CURRENT_USER_ID), getUsers()]);
		if (!practicesResult.success) {
			setLoadError(practicesResult.error.message);
			return;
		}
		setPractices(practicesResult.data);
		if (usersResult.success) {
			setUsers(usersResult.data);
		}
	}, []);

	useEffect(() => {
		void loadData();
	}, [loadData]);

	const handleSubmit = async (values: PracticeFormValues) => {
		setIsSubmitting(true);
		setSubmitError(null);
		const result = await createPractice({
			userId: CURRENT_USER_ID,
			dojoId: "dojo-001",
			...values,
		});
		setIsSubmitting(false);

		if (result.success) {
			setShowForm(false);
			await loadData();
			return;
		}
		setSubmitError(result.error.message);
	};

	return (
		<div>
			<div
				style={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center",
					marginBottom: "1.5rem",
				}}
			>
				<h1>稽古日誌</h1>
				<button
					type="button"
					onClick={() => setShowForm(!showForm)}
					style={{
						padding: "0.5rem 1rem",
						backgroundColor: showForm ? "#999" : "#1a1a2e",
						color: "#fff",
						border: "none",
						borderRadius: "4px",
						cursor: "pointer",
					}}
				>
					{showForm ? "閉じる" : "新規記録"}
				</button>
			</div>

			{loadError && (
				<div role="alert" style={alertStyle}>
					稽古記録の読み込みに失敗しました: {loadError}
				</div>
			)}

			{showForm && (
				<div style={{ marginBottom: "2rem" }}>
					<h2>稽古を記録</h2>
					{submitError && (
						<div role="alert" style={alertStyle}>
							{submitError}
						</div>
					)}
					<PracticeForm onSubmit={handleSubmit} isSubmitting={isSubmitting} />
				</div>
			)}

			<section style={{ marginBottom: "2rem" }}>
				<h2>的中率の推移</h2>
				<HitRateChart practices={practices} />
			</section>

			<section>
				<h2>稽古履歴</h2>
				<PracticeList practices={practices} users={users} />
			</section>
		</div>
	);
}
