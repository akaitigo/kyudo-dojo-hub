import { lazy, StrictMode, Suspense } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import { AppLayout } from "@/components/layout/AppLayout";

// Lazy-load route pages to reduce initial bundle size
const HomePage = lazy(() => import("@/pages/HomePage").then((m) => ({ default: m.HomePage })));
const PracticesPage = lazy(() => import("@/pages/PracticesPage").then((m) => ({ default: m.PracticesPage })));
const VideoAnalysisPage = lazy(() =>
	import("@/pages/VideoAnalysisPage").then((m) => ({
		default: m.VideoAnalysisPage,
	})),
);
const ExamChecklistPage = lazy(() =>
	import("@/pages/ExamChecklistPage").then((m) => ({
		default: m.ExamChecklistPage,
	})),
);
const DashboardPage = lazy(() => import("@/pages/DashboardPage").then((m) => ({ default: m.DashboardPage })));

const rootElement = document.getElementById("root");
if (rootElement) {
	createRoot(rootElement).render(
		<StrictMode>
			<BrowserRouter>
				<Suspense fallback={<div style={{ padding: "2rem", textAlign: "center" }}>読み込み中...</div>}>
					<Routes>
						<Route element={<AppLayout />}>
							<Route index element={<HomePage />} />
							<Route path="practices" element={<PracticesPage />} />
							<Route path="video-analysis" element={<VideoAnalysisPage />} />
							<Route path="exam-checklist" element={<ExamChecklistPage />} />
							<Route path="dashboard" element={<DashboardPage />} />
						</Route>
					</Routes>
				</Suspense>
			</BrowserRouter>
		</StrictMode>,
	);
}
