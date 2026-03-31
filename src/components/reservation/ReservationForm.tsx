import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { generateTimeSlots, type ReservationFormValues, reservationFormSchema } from "@/lib/reservation-validation";
import type { Dojo } from "@/types/domain";

interface ReservationFormProps {
	readonly dojo: Dojo;
	readonly onSubmit: (values: ReservationFormValues) => void;
	readonly isSubmitting?: boolean;
}

export function ReservationForm({ dojo, onSubmit, isSubmitting = false }: ReservationFormProps) {
	const {
		register,
		handleSubmit,
		formState: { errors: rawErrors },
	} = useForm({
		resolver: zodResolver(reservationFormSchema),
		defaultValues: {
			date: new Date().toISOString().split("T")[0] ?? "",
			startTime: "",
			laneNumber: 1,
		},
	});

	const timeSlots = generateTimeSlots(dojo.openTime, dojo.closeTime);
	const lanes = Array.from({ length: dojo.targetLanes }, (_, i) => i + 1);

	const fieldStyle = {
		padding: "0.5rem",
		border: "1px solid #ccc",
		borderRadius: "4px",
		fontSize: "1rem",
	} as const;

	const errorStyle = { color: "#d32f2f", fontSize: "0.85rem" } as const;

	return (
		<form
			onSubmit={handleSubmit((values) => onSubmit(values as ReservationFormValues))}
			style={{
				display: "flex",
				flexWrap: "wrap",
				gap: "1rem",
				alignItems: "flex-end",
			}}
		>
			<div>
				<label
					htmlFor="res-date"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					日付
				</label>
				<input id="res-date" type="date" {...register("date")} style={fieldStyle} />
				{rawErrors.date?.message && <p style={errorStyle}>{rawErrors.date.message}</p>}
			</div>

			<div>
				<label
					htmlFor="res-time"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					開始時刻
				</label>
				<select id="res-time" {...register("startTime")} style={fieldStyle}>
					<option value="">選択</option>
					{timeSlots.map((t) => (
						<option key={t} value={t}>
							{t}
						</option>
					))}
				</select>
				{rawErrors.startTime?.message && <p style={errorStyle}>{rawErrors.startTime.message}</p>}
			</div>

			<div>
				<label
					htmlFor="res-lane"
					style={{
						display: "block",
						marginBottom: "0.25rem",
						fontWeight: "bold",
					}}
				>
					的場番号
				</label>
				<select id="res-lane" {...register("laneNumber")} style={fieldStyle}>
					{lanes.map((l) => (
						<option key={l} value={l}>
							{l}
						</option>
					))}
				</select>
				{rawErrors.laneNumber?.message && <p style={errorStyle}>{rawErrors.laneNumber.message}</p>}
			</div>

			<button
				type="submit"
				disabled={isSubmitting}
				style={{
					padding: "0.5rem 1.5rem",
					backgroundColor: "#1a1a2e",
					color: "#fff",
					border: "none",
					borderRadius: "4px",
					cursor: isSubmitting ? "not-allowed" : "pointer",
					opacity: isSubmitting ? 0.7 : 1,
					height: "fit-content",
				}}
			>
				{isSubmitting ? "予約中..." : "予約する"}
			</button>
		</form>
	);
}
