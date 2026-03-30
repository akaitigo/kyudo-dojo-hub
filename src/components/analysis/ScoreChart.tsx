import { HASSETSU_LABELS, HASSETSU_PHASES } from "@/types/domain";
import type { HassetsuScores } from "@/types/domain";
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

interface ScoreChartProps {
	readonly scores: HassetsuScores;
}

export function ScoreChart({ scores }: ScoreChartProps) {
	const chartData = HASSETSU_PHASES.map((phase) => ({
		name: HASSETSU_LABELS[phase],
		score: scores[phase],
	}));

	return (
		<div>
			<h3 style={{ marginBottom: "0.5rem" }}>八節スコア</h3>
			<div style={{ width: "100%", height: 300 }}>
				<ResponsiveContainer>
					<BarChart data={chartData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
						<CartesianGrid strokeDasharray="3 3" />
						<XAxis dataKey="name" fontSize={11} angle={-30} textAnchor="end" height={60} />
						<YAxis domain={[0, 100]} fontSize={12} />
						<Tooltip />
						<Bar dataKey="score" name="スコア" fill="#1a1a2e" radius={[4, 4, 0, 0]} />
					</BarChart>
				</ResponsiveContainer>
			</div>
		</div>
	);
}
