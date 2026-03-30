import { AppLayout } from "@/components/layout/AppLayout";
import { ExamChecklistPage } from "@/pages/ExamChecklistPage";
import { HomePage } from "@/pages/HomePage";
import { PracticesPage } from "@/pages/PracticesPage";
import { VideoAnalysisPage } from "@/pages/VideoAnalysisPage";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";

const rootElement = document.getElementById("root");
if (rootElement) {
	createRoot(rootElement).render(
		<StrictMode>
			<BrowserRouter>
				<Routes>
					<Route element={<AppLayout />}>
						<Route index element={<HomePage />} />
						<Route path="practices" element={<PracticesPage />} />
						<Route path="exam-checklist" element={<ExamChecklistPage />} />
						<Route path="video-analysis" element={<VideoAnalysisPage />} />
					</Route>
				</Routes>
			</BrowserRouter>
		</StrictMode>,
	);
}
