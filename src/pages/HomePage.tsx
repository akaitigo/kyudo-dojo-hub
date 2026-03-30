import { Link } from "react-router";

export function HomePage() {
	return (
		<div>
			<h1>kyudo-dojo-hub</h1>
			<p>弓道の稽古記録・射形分析・道場運営プラットフォーム</p>
			<nav
				style={{
					marginTop: "2rem",
					display: "flex",
					flexDirection: "column",
					gap: "1rem",
				}}
			>
				<Link to="/practices" style={{ fontSize: "1.1rem" }}>
					稽古日誌
				</Link>
				<Link to="/video-analysis" style={{ fontSize: "1.1rem" }}>
					射形動画分析
				</Link>
				<Link to="/exam-checklist" style={{ fontSize: "1.1rem" }}>
					段位審査チェックリスト
				</Link>
			</nav>
		</div>
	);
}
