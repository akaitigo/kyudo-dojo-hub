import { useCallback, useEffect, useState } from "react";
import { HitRateChart } from "@/components/practice/HitRateChart";
import { PracticeForm } from "@/components/practice/PracticeForm";
import { PracticeList } from "@/components/practice/PracticeList";
import { createPractice, getPractices } from "@/lib/api";
import type { PracticeFormValues } from "@/lib/validation";
import type { Practice } from "@/types/domain";

/** 現在のモックユーザーID */
const CURRENT_USER_ID = "user-001";

export function PracticesPage() {
	const [practices, setPractices] = useState<readonly Practice[]>([]);
	const [isSubmitting, setIsSubmitting] = useState(false);
	const [showForm, setShowForm] = useState(false);

	const loadPractices = useCallback(async () => {
		const result = await getPractices(CURRENT_USER_ID);
		if (result.success) {
			setPractices(result.data);
		}
	}, []);

	useEffect(() => {
		void loadPractices();
	}, [loadPractices]);

	const handleSubmit = async (values: PracticeFormValues) => {
		setIsSubmitting(true);
		const result = await createPractice({
			userId: CURRENT_USER_ID,
			dojoId: "dojo-001",
			...values,
		});
		setIsSubmitting(false);

		if (result.success) {
			setShowForm(false);
			await loadPractices();
		}
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

			{showForm && (
				<div style={{ marginBottom: "2rem" }}>
					<h2>稽古を記録</h2>
					<PracticeForm onSubmit={handleSubmit} isSubmitting={isSubmitting} />
				</div>
			)}

			<section style={{ marginBottom: "2rem" }}>
				<h2>的中率の推移</h2>
				<HitRateChart practices={practices} />
			</section>

			<section>
				<h2>稽古履歴</h2>
				<PracticeList practices={practices} />
			</section>
		</div>
	);
}
