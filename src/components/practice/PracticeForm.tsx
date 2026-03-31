import { zodResolver } from "@hookform/resolvers/zod";
import { type FieldErrors, useForm } from "react-hook-form";
import { type PracticeFormValues, practiceFormSchema } from "@/lib/validation";

interface PracticeFormProps {
	readonly onSubmit: (values: PracticeFormValues) => void;
	readonly isSubmitting?: boolean;
}

export function PracticeForm({ onSubmit, isSubmitting = false }: PracticeFormProps) {
	const {
		register,
		handleSubmit,
		formState: { errors: rawErrors },
	} = useForm({
		resolver: zodResolver(practiceFormSchema),
		defaultValues: {
			date: new Date().toISOString().split("T")[0] ?? "",
			hitRate: 0,
			arrowCount: 1,
			notes: "",
			instructorComment: "",
		},
	});

	const errors = rawErrors as FieldErrors<PracticeFormValues>;

	const fieldStyle = {
		width: "100%",
		padding: "0.5rem",
		border: "1px solid #ccc",
		borderRadius: "4px",
		fontSize: "1rem",
	} as const;

	const errorStyle = {
		color: "#d32f2f",
		fontSize: "0.85rem",
		marginTop: "0.25rem",
	} as const;

	return (
		<form
			onSubmit={handleSubmit((values) => onSubmit(values as PracticeFormValues))}
			style={{
				display: "flex",
				flexDirection: "column",
				gap: "1rem",
				maxWidth: "600px",
			}}
		>
			<div>
				<label
					htmlFor="date"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					日付
				</label>
				<input id="date" type="date" {...register("date")} style={fieldStyle} />
				{errors.date?.message && <p style={errorStyle}>{errors.date.message}</p>}
			</div>

			<div style={{ display: "flex", gap: "1rem" }}>
				<div style={{ flex: 1 }}>
					<label
						htmlFor="hitRate"
						style={{
							display: "block",
							marginBottom: "0.25rem",
							fontWeight: "bold",
						}}
					>
						的中率 (%)
					</label>
					<input id="hitRate" type="number" min={0} max={100} {...register("hitRate")} style={fieldStyle} />
					{errors.hitRate?.message && <p style={errorStyle}>{errors.hitRate.message}</p>}
				</div>
				<div style={{ flex: 1 }}>
					<label
						htmlFor="arrowCount"
						style={{
							display: "block",
							marginBottom: "0.25rem",
							fontWeight: "bold",
						}}
					>
						矢数
					</label>
					<input id="arrowCount" type="number" min={1} max={1000} {...register("arrowCount")} style={fieldStyle} />
					{errors.arrowCount?.message && <p style={errorStyle}>{errors.arrowCount.message}</p>}
				</div>
			</div>

			<div>
				<label
					htmlFor="notes"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					気づき
				</label>
				<textarea id="notes" rows={4} {...register("notes")} style={fieldStyle} />
				{errors.notes?.message && <p style={errorStyle}>{errors.notes.message}</p>}
			</div>

			<div>
				<label
					htmlFor="instructorComment"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					師範コメント
				</label>
				<textarea id="instructorComment" rows={3} {...register("instructorComment")} style={fieldStyle} />
				{errors.instructorComment?.message && <p style={errorStyle}>{errors.instructorComment.message}</p>}
			</div>

			<button
				type="submit"
				disabled={isSubmitting}
				style={{
					padding: "0.75rem",
					backgroundColor: "#1a1a2e",
					color: "#fff",
					border: "none",
					borderRadius: "4px",
					fontSize: "1rem",
					cursor: isSubmitting ? "not-allowed" : "pointer",
					opacity: isSubmitting ? 0.7 : 1,
				}}
			>
				{isSubmitting ? "保存中..." : "稽古を記録"}
			</button>
		</form>
	);
}
