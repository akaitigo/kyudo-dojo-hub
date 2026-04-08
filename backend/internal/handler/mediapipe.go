package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/model"
)

const mediaPipeWorkerURL = "http://localhost:8081/analyze"

// mediaPipeRequest is the payload sent to the Python MediaPipe worker.
type mediaPipeRequest struct {
	VideoID  string  `json:"video_id"`
	Duration float64 `json:"duration"`
}

// mediaPipeResponse is the payload returned by the Python MediaPipe worker.
type mediaPipeResponse struct {
	Scores       model.HassetsuScores `json:"scores"`
	Phases       []model.PhaseSegment `json:"phases"`
	OverallScore int                  `json:"overall_score"`
	Feedback     string               `json:"feedback"`
}

// callMediaPipeWorker sends an analysis request to the Python MediaPipe worker.
func callMediaPipeWorker(ctx context.Context, video model.Video, userID string) (model.Analysis, error) {
	reqBody := mediaPipeRequest{
		VideoID:  video.ID,
		Duration: video.Duration,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return model.Analysis{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, mediaPipeWorkerURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return model.Analysis{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Analysis{}, fmt.Errorf("failed to call MediaPipe worker: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return model.Analysis{}, fmt.Errorf("MediaPipe worker returned status %d, failed to read body: %w", resp.StatusCode, readErr)
		}
		return model.Analysis{}, fmt.Errorf("MediaPipe worker returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var mpResp mediaPipeResponse
	if err := json.NewDecoder(resp.Body).Decode(&mpResp); err != nil {
		return model.Analysis{}, fmt.Errorf("failed to decode MediaPipe response: %w", err)
	}

	now := time.Now()
	analysis := model.Analysis{
		ID:           fmt.Sprintf("analysis-%d-%s", now.UnixNano(), randomString(5)),
		VideoID:      video.ID,
		UserID:       userID,
		Scores:       mpResp.Scores,
		Phases:       mpResp.Phases,
		OverallScore: mpResp.OverallScore,
		Feedback:     mpResp.Feedback,
		CreatedAt:    now,
	}
	return analysis, nil
}

// randomString is defined in store package but we need a local copy here.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// generateFallbackAnalysis produces a deterministic simulated analysis result
// when the Python MediaPipe worker is unavailable. It uses the video duration
// to produce a reasonable eight-phase timeline and deterministic scores.
func generateFallbackAnalysis(video model.Video, userID string) model.Analysis {
	duration := video.Duration
	if duration <= 0 {
		duration = 35.0
	}

	// Typical phase duration ratios for a kyudo shot (total = 1.0)
	ratios := []float64{0.08, 0.10, 0.14, 0.11, 0.22, 0.17, 0.03, 0.15}
	phases := make([]model.PhaseSegment, len(model.AllPhases))
	cursor := 0.0
	for i, phase := range model.AllPhases {
		start := cursor
		end := cursor + duration*ratios[i]
		phases[i] = model.PhaseSegment{
			Phase:     phase,
			StartTime: math.Round(start*10) / 10,
			EndTime:   math.Round(end*10) / 10,
		}
		cursor = end
	}

	// Deterministic scores based on video file name hash
	hashVal := 0
	for _, c := range video.FileName {
		hashVal = (hashVal*31 + int(c)) & 0x7fffffff
	}
	baseScore := 60 + hashVal%25 // 60-84 range

	scores := model.HassetsuScores{
		Ashibumi:    clampScore(baseScore + (hashVal>>3)%10 - 5),
		Dozukuri:    clampScore(baseScore + (hashVal>>5)%10 - 5),
		Yugamae:     clampScore(baseScore + (hashVal>>7)%10 - 5),
		Uchiokoshi:  clampScore(baseScore + (hashVal>>9)%10 - 5),
		Hikiwake:    clampScore(baseScore + (hashVal>>11)%10 - 5),
		Kai:         clampScore(baseScore + (hashVal>>13)%10 - 5),
		Hanare:      clampScore(baseScore + (hashVal>>15)%10 - 5),
		Zanshin:     clampScore(baseScore + (hashVal>>17)%10 - 5),
	}

	overallScore := (scores.Ashibumi + scores.Dozukuri + scores.Yugamae +
		scores.Uchiokoshi + scores.Hikiwake + scores.Kai +
		scores.Hanare + scores.Zanshin) / 8

	feedback := generateFeedback(scores)

	now := time.Now()
	return model.Analysis{
		ID:           fmt.Sprintf("analysis-%d", now.UnixNano()),
		VideoID:      video.ID,
		UserID:       userID,
		Scores:       scores,
		Phases:       phases,
		OverallScore: overallScore,
		Feedback:     feedback,
		CreatedAt:    now,
	}
}

func clampScore(s int) int {
	if s < 0 {
		return 0
	}
	if s > 100 {
		return 100
	}
	return s
}

// generateFeedback produces Japanese feedback text based on scores.
func generateFeedback(scores model.HassetsuScores) string {
	type phaseInfo struct {
		name  string
		score int
	}
	all := []phaseInfo{
		{"足踏み", scores.Ashibumi},
		{"胴造り", scores.Dozukuri},
		{"弓構え", scores.Yugamae},
		{"打起し", scores.Uchiokoshi},
		{"引分け", scores.Hikiwake},
		{"会", scores.Kai},
		{"離れ", scores.Hanare},
		{"残心", scores.Zanshin},
	}

	// Find best and worst phases
	best := all[0]
	worst := all[0]
	for _, p := range all[1:] {
		if p.score > best.score {
			best = p
		}
		if p.score < worst.score {
			worst = p
		}
	}

	return fmt.Sprintf(
		"%sの安定感が良好（%d点）。%sに改善の余地あり（%d点）。全体的なバランスを意識して稽古を続けること。",
		best.name, best.score, worst.name, worst.score,
	)
}
