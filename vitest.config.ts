import { resolve } from "node:path";
import { defineConfig } from "vitest/config";

export default defineConfig({
	resolve: {
		alias: {
			"@": resolve(__dirname, "src"),
		},
	},
	test: {
		exclude: ["node_modules", "dist", "test/e2e/**"],
		environment: "node",
	},
});
