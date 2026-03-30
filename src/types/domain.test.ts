import { describe, expect, it } from "vitest";
import { DAN_LABELS, DAN_RANKS, HASSETSU_LABELS, HASSETSU_PHASES, SHOGO_LABELS, SHOGO_TITLES } from "./domain";

describe("ドメイン型定義", () => {
	describe("射法八節", () => {
		it("8つのフェーズが定義されている", () => {
			expect(HASSETSU_PHASES).toHaveLength(8);
		});

		it("全フェーズに日本語ラベルがある", () => {
			for (const phase of HASSETSU_PHASES) {
				expect(HASSETSU_LABELS[phase]).toBeDefined();
				expect(HASSETSU_LABELS[phase].length).toBeGreaterThan(0);
			}
		});

		it("正しい順序で定義されている", () => {
			expect(HASSETSU_PHASES[0]).toBe("ashibumi");
			expect(HASSETSU_PHASES[7]).toBe("zanshin");
		});
	});

	describe("段位", () => {
		it("10段位が定義されている", () => {
			expect(DAN_RANKS).toHaveLength(10);
		});

		it("全段位に日本語ラベルがある", () => {
			for (const rank of DAN_RANKS) {
				expect(DAN_LABELS[rank]).toBeDefined();
			}
		});

		it("初段から十段の順序", () => {
			expect(DAN_RANKS[0]).toBe("shodan");
			expect(DAN_RANKS[9]).toBe("judan");
		});
	});

	describe("称号", () => {
		it("3つの称号が定義されている", () => {
			expect(SHOGO_TITLES).toHaveLength(3);
		});

		it("全称号に日本語ラベルがある", () => {
			for (const title of SHOGO_TITLES) {
				expect(SHOGO_LABELS[title]).toBeDefined();
			}
		});
	});
});
