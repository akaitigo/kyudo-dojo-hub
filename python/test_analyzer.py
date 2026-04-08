"""analyzer モジュールのユニットテスト"""

from __future__ import annotations

import pytest

from analyzer import (
    HASSETSU_LABELS,
    HASSETSU_PHASES,
    AnalysisResult,
    HassetsuAnalyzer,
    PhaseSegment,
)


@pytest.fixture
def analyzer() -> HassetsuAnalyzer:
    """テスト用の HassetsuAnalyzer インスタンスを生成"""
    return HassetsuAnalyzer()


class TestHassetsuPhases:
    """八節フェーズ定義のテスト"""

    def test_phase_count(self) -> None:
        assert len(HASSETSU_PHASES) == 8

    def test_phase_order(self) -> None:
        expected = [
            "ashibumi",
            "dozukuri",
            "yugamae",
            "uchiokoshi",
            "hikiwake",
            "kai",
            "hanare",
            "zanshin",
        ]
        assert HASSETSU_PHASES == expected

    def test_all_phases_have_labels(self) -> None:
        for phase in HASSETSU_PHASES:
            assert phase in HASSETSU_LABELS
            assert len(HASSETSU_LABELS[phase]) > 0


class TestHassetsuAnalyzer:
    """HassetsuAnalyzer のテスト"""

    def test_mode_is_simulation_without_mediapipe(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        # MediaPipe がインストールされていない環境では simulation モード
        assert analyzer.mode in ("mediapipe", "simulation")

    def test_analyze_returns_analysis_result(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        assert isinstance(result, AnalysisResult)

    def test_analyze_has_all_phase_scores(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        for phase in HASSETSU_PHASES:
            assert phase in result.scores
            assert 0 <= result.scores[phase] <= 100

    def test_analyze_has_eight_phase_segments(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        assert len(result.phases) == 8
        for segment in result.phases:
            assert isinstance(segment, PhaseSegment)
            assert segment.start_time >= 0
            assert segment.end_time > segment.start_time

    def test_phases_are_contiguous(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        for i in range(len(result.phases) - 1):
            assert result.phases[i].end_time == pytest.approx(
                result.phases[i + 1].start_time, abs=0.2
            )

    def test_phases_cover_full_duration(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        duration = 45.0
        result = analyzer.analyze("video-001", duration)
        assert result.phases[0].start_time == 0.0
        # Last phase should end close to the full duration
        assert result.phases[-1].end_time == pytest.approx(duration, abs=1.0)

    def test_overall_score_is_average(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        expected_avg = sum(result.scores.values()) // len(result.scores)
        assert result.overall_score == expected_avg

    def test_feedback_is_nonempty(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        result = analyzer.analyze("video-001", 45.0)
        assert len(result.feedback) > 0

    def test_deterministic_results(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        """同じ入力に対して同じ結果を返すことを確認"""
        r1 = analyzer.analyze("video-001", 45.0)
        r2 = analyzer.analyze("video-001", 45.0)
        assert r1.scores == r2.scores
        assert r1.overall_score == r2.overall_score
        assert r1.feedback == r2.feedback

    def test_different_videos_get_different_scores(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        """異なる動画IDに対して異なるスコアを返すことを確認"""
        r1 = analyzer.analyze("video-001", 45.0)
        r2 = analyzer.analyze("video-002", 45.0)
        assert r1.scores != r2.scores

    def test_zero_duration_defaults(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        """duration=0 のときデフォルト値が使われることを確認"""
        result = analyzer.analyze("video-001", 0)
        assert len(result.phases) == 8
        assert result.phases[-1].end_time > 0

    def test_negative_duration_defaults(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        """負の duration のときデフォルト値が使われることを確認"""
        result = analyzer.analyze("video-001", -10)
        assert len(result.phases) == 8
        assert result.phases[-1].end_time > 0

    def test_close_does_not_raise(
        self, analyzer: HassetsuAnalyzer
    ) -> None:
        """close() がエラーなく完了することを確認"""
        analyzer.close()
