"""
射形分析 MediaPipe ワーカーサーバー

Go API バックエンドから HTTP 経由で分析リクエストを受け取り、
HassetsuAnalyzer で八節分析を実行して結果を返す。

起動:
    python server.py [--port 8081]

環境変数:
    MEDIAPIPE_WORKER_HOST: バインドアドレス (デフォルト: 127.0.0.1)
    MEDIAPIPE_WORKER_PORT: ポート番号 (デフォルト: 8081)
"""

from __future__ import annotations

import argparse
import logging
import os
import sys
from dataclasses import asdict

from flask import Flask, Response, jsonify, request

from analyzer import HassetsuAnalyzer

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)
logger = logging.getLogger(__name__)

app = Flask(__name__)
analyzer = HassetsuAnalyzer()


@app.route("/health", methods=["GET"])
def health() -> tuple[Response, int]:
    """ヘルスチェックエンドポイント"""
    resp: Response = jsonify({"status": "ok", "mode": analyzer.mode})
    return resp, 200


@app.route("/analyze", methods=["POST"])
def analyze() -> tuple[Response, int]:
    """
    射形分析エンドポイント

    Request body (JSON):
        video_id: str  - 動画ID
        duration: float - 動画長 (秒)

    Response (JSON):
        scores: dict[str, int]  - 各フェーズのスコア
        phases: list[dict]      - フェーズ分割結果
        overall_score: int      - 総合スコア
        feedback: str           - 分析コメント
    """
    data = request.get_json(silent=True)
    if data is None:
        resp: Response = jsonify({"error": "リクエストボディが不正です"})
        return resp, 400

    video_id = data.get("video_id", "")
    duration = data.get("duration", 0.0)

    if not video_id:
        resp = jsonify({"error": "video_id は必須です"})
        return resp, 400

    if not isinstance(duration, (int, float)) or duration < 0:
        resp = jsonify({"error": "duration は正の数値を指定してください"})
        return resp, 400

    logger.info(
        "Analyzing video_id=%s, duration=%.1f, mode=%s",
        video_id,
        duration,
        analyzer.mode,
    )

    result = analyzer.analyze(video_id, float(duration))

    # Convert dataclass to dict for JSON serialization
    phases_list = [
        {
            "phase": p.phase,
            "startTime": p.start_time,
            "endTime": p.end_time,
        }
        for p in result.phases
    ]

    response_data = {
        "scores": result.scores,
        "phases": phases_list,
        "overall_score": result.overall_score,
        "feedback": result.feedback,
    }

    logger.info(
        "Analysis complete: overall_score=%d, mode=%s",
        result.overall_score,
        analyzer.mode,
    )

    resp = jsonify(response_data)
    return resp, 200


def main() -> None:
    """エントリーポイント"""
    parser = argparse.ArgumentParser(description="射形分析 MediaPipe ワーカー")
    parser.add_argument(
        "--port",
        type=int,
        default=int(os.environ.get("MEDIAPIPE_WORKER_PORT", "8081")),
        help="ポート番号 (default: 8081)",
    )
    args = parser.parse_args()

    host = os.environ.get("MEDIAPIPE_WORKER_HOST", "127.0.0.1")
    logger.info(
        "Starting MediaPipe worker on %s:%d (mode: %s)",
        host,
        args.port,
        analyzer.mode,
    )
    app.run(host=host, port=args.port, debug=False)


if __name__ == "__main__":
    main()
