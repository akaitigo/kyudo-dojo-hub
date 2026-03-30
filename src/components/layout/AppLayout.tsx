import { NavLink, Outlet } from "react-router";

const navItems = [
	{ to: "/", label: "ホーム" },
	{ to: "/practices", label: "稽古日誌" },
	{ to: "/video-analysis", label: "動画分析" },
	{ to: "/exam-checklist", label: "審査チェック" },
	{ to: "/dashboard", label: "道場管理" },
] as const;

export function AppLayout() {
	return (
		<div style={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
			<header
				style={{
					borderBottom: "1px solid #e0e0e0",
					padding: "0.75rem 1rem",
					backgroundColor: "#1a1a2e",
					color: "#fff",
				}}
			>
				<nav
					style={{
						display: "flex",
						alignItems: "center",
						gap: "1.5rem",
						maxWidth: "1200px",
						margin: "0 auto",
					}}
				>
					<span style={{ fontWeight: "bold", fontSize: "1.1rem" }}>kyudo-dojo-hub</span>
					{navItems.map((item) => (
						<NavLink
							key={item.to}
							to={item.to}
							style={({ isActive }) => ({
								color: isActive ? "#ffd700" : "#ccc",
								textDecoration: "none",
								fontSize: "0.9rem",
							})}
						>
							{item.label}
						</NavLink>
					))}
				</nav>
			</header>
			<main
				style={{
					flex: 1,
					padding: "1rem",
					maxWidth: "1200px",
					margin: "0 auto",
					width: "100%",
				}}
			>
				<Outlet />
			</main>
		</div>
	);
}
