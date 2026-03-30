import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

function App() {
	return (
		<main>
			<h1>kyudo-dojo-hub</h1>
			<p>弓道の稽古記録・射形分析・道場運営プラットフォーム</p>
		</main>
	);
}

const rootElement = document.getElementById("root");
if (rootElement) {
	createRoot(rootElement).render(
		<StrictMode>
			<App />
		</StrictMode>,
	);
}
