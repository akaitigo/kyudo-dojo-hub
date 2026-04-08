"""server モジュールのユニットテスト"""

from __future__ import annotations

import os
from unittest.mock import patch

import pytest

from server import app


@pytest.fixture
def client():
    """テスト用 Flask クライアント"""
    app.config["TESTING"] = True
    with app.test_client() as c:
        yield c


class TestHealthEndpoint:
    """ヘルスチェックエンドポイントのテスト"""

    def test_health_returns_ok(self, client) -> None:
        resp = client.get("/health")
        assert resp.status_code == 200
        data = resp.get_json()
        assert data["status"] == "ok"


class TestAnalyzeEndpoint:
    """分析エンドポイントのテスト"""

    def test_analyze_success(self, client) -> None:
        resp = client.post(
            "/analyze",
            json={"video_id": "test-001", "duration": 30.0},
        )
        assert resp.status_code == 200
        data = resp.get_json()
        assert "scores" in data
        assert "phases" in data
        assert "overall_score" in data
        assert "feedback" in data

    def test_analyze_missing_body(self, client) -> None:
        resp = client.post("/analyze", data="not json", content_type="text/plain")
        assert resp.status_code == 400

    def test_analyze_missing_video_id(self, client) -> None:
        resp = client.post("/analyze", json={"duration": 30.0})
        assert resp.status_code == 400

    def test_analyze_invalid_duration(self, client) -> None:
        resp = client.post(
            "/analyze",
            json={"video_id": "test-001", "duration": -5},
        )
        assert resp.status_code == 400


class TestHostBinding:
    """ホストバインド設定のテスト"""

    def test_default_host_is_localhost(self) -> None:
        """MEDIAPIPE_WORKER_HOST 未設定時は 127.0.0.1"""
        with patch.dict(os.environ, {}, clear=True):
            host = os.environ.get("MEDIAPIPE_WORKER_HOST", "127.0.0.1")
            assert host == "127.0.0.1"

    def test_custom_host_from_env(self) -> None:
        """MEDIAPIPE_WORKER_HOST 設定時はその値を使用"""
        with patch.dict(os.environ, {"MEDIAPIPE_WORKER_HOST": "0.0.0.0"}):
            host = os.environ.get("MEDIAPIPE_WORKER_HOST", "127.0.0.1")
            assert host == "0.0.0.0"
