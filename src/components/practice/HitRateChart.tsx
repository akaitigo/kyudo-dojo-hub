import type { Practice } from "@/types/domain";
import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

interface HitRateChartProps {
	readonly practices: readonly Practice[];
}

export function HitRateChart({ practices }: HitRateChartProps) {
	const chartData = [...practices]
		.sort((a, b) => a.date.localeCompare(b.date))
		.map((p) => ({
			date: p.date,
			hitRate: p.hitRate,
		}));

	if (chartData.length === 0) {
		return <p>まだ稽古記録がありません</p>;
	}

	return (
		<div style={{ width: "100%", height: 300 }}>
			<ResponsiveContainer>
				<LineChart data={chartData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
					<CartesianGrid strokeDasharray="3 3" />
					<XAxis dataKey="date" fontSize={12} />
					<YAxis domain={[0, 100]} fontSize={12} />
					<Tooltip />
					<Line type="monotone" dataKey="hitRate" name="的中率 (%)" stroke="#1a1a2e" strokeWidth={2} />
				</LineChart>
			</ResponsiveContainer>
		</div>
	);
}
