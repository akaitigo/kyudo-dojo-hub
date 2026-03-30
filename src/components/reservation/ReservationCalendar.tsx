import { MOCK_USERS } from "@/lib/mock-data";
import { generateTimeSlots } from "@/lib/reservation-validation";
import type { Dojo, Reservation } from "@/types/domain";

interface ReservationCalendarProps {
	readonly dojo: Dojo;
	readonly reservations: readonly Reservation[];
	readonly selectedDate: string;
	readonly onDateChange: (date: string) => void;
	readonly onDeleteReservation: (id: string) => void;
}

function getUserName(userId: string): string {
	return MOCK_USERS.find((u) => u.id === userId)?.name ?? userId;
}

export function ReservationCalendar({
	dojo,
	reservations,
	selectedDate,
	onDateChange,
	onDeleteReservation,
}: ReservationCalendarProps) {
	const timeSlots = generateTimeSlots(dojo.openTime, dojo.closeTime);
	const lanes = Array.from({ length: dojo.targetLanes }, (_, i) => i + 1);

	const getReservation = (lane: number, time: string) =>
		reservations.find((r) => r.laneNumber === lane && r.startTime === time && r.date === selectedDate);

	return (
		<div>
			<div
				style={{
					display: "flex",
					alignItems: "center",
					gap: "1rem",
					marginBottom: "1rem",
				}}
			>
				<label htmlFor="calendar-date" style={{ fontWeight: "bold" }}>
					日付:
				</label>
				<input
					id="calendar-date"
					type="date"
					value={selectedDate}
					onChange={(e) => onDateChange(e.target.value)}
					style={{
						padding: "0.5rem",
						border: "1px solid #ccc",
						borderRadius: "4px",
					}}
				/>
			</div>

			<div style={{ overflowX: "auto" }}>
				<table
					style={{
						borderCollapse: "collapse",
						width: "100%",
						minWidth: "600px",
					}}
				>
					<thead>
						<tr>
							<th
								style={{
									border: "1px solid #e0e0e0",
									padding: "0.5rem",
									backgroundColor: "#f5f5f5",
								}}
							>
								時間
							</th>
							{lanes.map((lane) => (
								<th
									key={lane}
									style={{
										border: "1px solid #e0e0e0",
										padding: "0.5rem",
										backgroundColor: "#f5f5f5",
									}}
								>
									的場 {lane}
								</th>
							))}
						</tr>
					</thead>
					<tbody>
						{timeSlots.map((time) => (
							<tr key={time}>
								<td
									style={{
										border: "1px solid #e0e0e0",
										padding: "0.5rem",
										fontWeight: "bold",
										textAlign: "center",
									}}
								>
									{time}
								</td>
								{lanes.map((lane) => {
									const reservation = getReservation(lane, time);
									return (
										<td
											key={lane}
											style={{
												border: "1px solid #e0e0e0",
												padding: "0.5rem",
												backgroundColor: reservation ? "#e3f2fd" : "#fff",
												textAlign: "center",
												fontSize: "0.85rem",
											}}
										>
											{reservation ? (
												<div>
													<div>{getUserName(reservation.userId)}</div>
													<button
														type="button"
														onClick={() => onDeleteReservation(reservation.id)}
														style={{
															marginTop: "0.25rem",
															padding: "0.15rem 0.5rem",
															fontSize: "0.75rem",
															backgroundColor: "#d32f2f",
															color: "#fff",
															border: "none",
															borderRadius: "3px",
															cursor: "pointer",
														}}
													>
														取消
													</button>
												</div>
											) : (
												<span style={{ color: "#ccc" }}>-</span>
											)}
										</td>
									);
								})}
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	);
}
