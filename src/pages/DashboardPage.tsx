import { useCallback, useEffect, useState } from "react";
import { MemberList } from "@/components/member/MemberList";
import { ReservationCalendar } from "@/components/reservation/ReservationCalendar";
import { ReservationForm } from "@/components/reservation/ReservationForm";
import type { DashboardSummary } from "@/lib/mock-api";
import {
	createReservation,
	deleteReservation,
	getDashboardSummary,
	getDojo,
	getReservations,
	getUsersByDojo,
} from "@/lib/mock-api";
import type { ReservationFormValues } from "@/lib/reservation-validation";
import { getEndTime } from "@/lib/reservation-validation";
import type { Dojo, Reservation, User } from "@/types/domain";

/** 現在のモック道場ID */
const CURRENT_DOJO_ID = "dojo-001";
const CURRENT_USER_ID = "user-002"; // 管理者ユーザー

export function DashboardPage() {
	const [dojo, setDojo] = useState<Dojo | null>(null);
	const [reservations, setReservations] = useState<readonly Reservation[]>([]);
	const [members, setMembers] = useState<readonly User[]>([]);
	const [summary, setSummary] = useState<DashboardSummary | null>(null);
	const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split("T")[0] ?? "");
	const [isSubmitting, setIsSubmitting] = useState(false);
	const [activeTab, setActiveTab] = useState<"calendar" | "members">("calendar");

	const loadData = useCallback(async () => {
		const [dojoResult, resResult, membersResult, summaryResult] = await Promise.all([
			getDojo(CURRENT_DOJO_ID),
			getReservations(CURRENT_DOJO_ID),
			getUsersByDojo(CURRENT_DOJO_ID),
			getDashboardSummary(CURRENT_DOJO_ID),
		]);

		if (dojoResult.success) setDojo(dojoResult.data);
		if (resResult.success) setReservations(resResult.data);
		if (membersResult.success) setMembers(membersResult.data);
		if (summaryResult.success) setSummary(summaryResult.data);
	}, []);

	useEffect(() => {
		void loadData();
	}, [loadData]);

	const handleCreateReservation = async (values: ReservationFormValues) => {
		setIsSubmitting(true);
		await createReservation({
			dojoId: CURRENT_DOJO_ID,
			userId: CURRENT_USER_ID,
			laneNumber: values.laneNumber,
			date: values.date,
			startTime: values.startTime,
			endTime: getEndTime(values.startTime),
		});
		setIsSubmitting(false);
		await loadData();
	};

	const handleDeleteReservation = async (id: string) => {
		await deleteReservation(id);
		await loadData();
	};

	if (!dojo) return <p>読み込み中...</p>;

	const tabStyle = (tab: "calendar" | "members") =>
		({
			padding: "0.5rem 1.5rem",
			backgroundColor: activeTab === tab ? "#1a1a2e" : "#e0e0e0",
			color: activeTab === tab ? "#fff" : "#333",
			border: "none",
			borderRadius: "4px 4px 0 0",
			cursor: "pointer",
			fontSize: "1rem",
		}) as const;

	return (
		<div>
			<h1>道場管理ダッシュボード</h1>
			<p style={{ color: "#666", marginBottom: "1.5rem" }}>{dojo.name}</p>

			{summary && (
				<div
					style={{
						display: "flex",
						gap: "1.5rem",
						marginBottom: "2rem",
						flexWrap: "wrap",
					}}
				>
					<div
						style={{
							padding: "1.5rem",
							backgroundColor: "#e3f2fd",
							borderRadius: "8px",
							flex: "1 1 200px",
							textAlign: "center",
						}}
					>
						<div style={{ fontSize: "2rem", fontWeight: "bold" }}>{summary.todayReservationCount}</div>
						<div style={{ color: "#555" }}>本日の予約</div>
					</div>
					<div
						style={{
							padding: "1.5rem",
							backgroundColor: "#e8f5e9",
							borderRadius: "8px",
							flex: "1 1 200px",
							textAlign: "center",
						}}
					>
						<div style={{ fontSize: "2rem", fontWeight: "bold" }}>{summary.totalMemberCount}</div>
						<div style={{ color: "#555" }}>会員数</div>
					</div>
				</div>
			)}

			<div style={{ display: "flex", gap: "0.25rem", marginBottom: "0" }}>
				<button type="button" onClick={() => setActiveTab("calendar")} style={tabStyle("calendar")}>
					予約カレンダー
				</button>
				<button type="button" onClick={() => setActiveTab("members")} style={tabStyle("members")}>
					会員管理
				</button>
			</div>

			<div
				style={{
					border: "1px solid #e0e0e0",
					borderRadius: "0 8px 8px 8px",
					padding: "1.5rem",
				}}
			>
				{activeTab === "calendar" && (
					<>
						<section style={{ marginBottom: "1.5rem" }}>
							<h2>新規予約</h2>
							<ReservationForm dojo={dojo} onSubmit={handleCreateReservation} isSubmitting={isSubmitting} />
						</section>

						<section>
							<h2>予約カレンダー</h2>
							<ReservationCalendar
								dojo={dojo}
								reservations={reservations}
								selectedDate={selectedDate}
								onDateChange={setSelectedDate}
								onDeleteReservation={handleDeleteReservation}
							/>
						</section>
					</>
				)}

				{activeTab === "members" && (
					<section>
						<h2>会員一覧 ({members.length}名)</h2>
						<MemberList members={members} />
					</section>
				)}
			</div>
		</div>
	);
}
