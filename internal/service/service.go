// Package service provides business logic for kyudo-dojo-hub.
package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ryusei/kyudo-dojo-hub/internal/model"
	"github.com/ryusei/kyudo-dojo-hub/internal/repository"
)

// ValidationError represents a client-side input error.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError represents a missing resource.
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s が見つかりません", e.Resource)
}

// ConflictError represents a scheduling conflict.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// Service provides business-level operations.
type Service struct {
	repo *repository.Repo
}

// New creates a new Service.
func New(repo *repository.Repo) *Service {
	return &Service{repo: repo}
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// ListUsers returns all users.
func (s *Service) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.ListUsers(ctx)
}

// GetUser returns a user by ID.
func (s *Service) GetUser(ctx context.Context, id string) (*model.User, error) {
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "ユーザー"}
		}
		return nil, err
	}
	return u, nil
}

// ListUsersByDojo returns users belonging to a specific dojo.
func (s *Service) ListUsersByDojo(ctx context.Context, dojoID string) ([]model.User, error) {
	return s.repo.ListUsersByDojo(ctx, dojoID)
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

// ListDojos returns all dojos.
func (s *Service) ListDojos(ctx context.Context) ([]model.Dojo, error) {
	return s.repo.ListDojos(ctx)
}

// GetDojo returns a dojo by ID.
func (s *Service) GetDojo(ctx context.Context, id string) (*model.Dojo, error) {
	d, err := s.repo.GetDojo(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "道場"}
		}
		return nil, err
	}
	return d, nil
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

// CreatePracticeInput represents input for creating a practice.
type CreatePracticeInput struct {
	UserID            string  `json:"userId"`
	DojoID            *string `json:"dojoId,omitempty"`
	Date              string  `json:"date"`
	HitRate           int     `json:"hitRate"`
	ArrowCount        int     `json:"arrowCount"`
	Notes             string  `json:"notes"`
	InstructorComment string  `json:"instructorComment"`
}

// ListPractices returns practices, optionally filtered by userID.
func (s *Service) ListPractices(ctx context.Context, userID *string) ([]model.Practice, error) {
	return s.repo.ListPractices(ctx, userID)
}

// GetPractice returns a practice by ID.
func (s *Service) GetPractice(ctx context.Context, id string) (*model.Practice, error) {
	p, err := s.repo.GetPractice(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "稽古記録"}
		}
		return nil, err
	}
	return p, nil
}

// CreatePractice validates and creates a new practice record.
func (s *Service) CreatePractice(ctx context.Context, input CreatePracticeInput) (*model.Practice, error) {
	if input.HitRate < 0 || input.HitRate > 100 {
		return nil, &ValidationError{Message: "的中率は 0〜100 の範囲で入力してください"}
	}
	if input.ArrowCount < 1 || input.ArrowCount > 1000 {
		return nil, &ValidationError{Message: "矢数は 1〜1000 の範囲で入力してください"}
	}
	if len(input.Notes) > 5000 {
		return nil, &ValidationError{Message: "気づきは 5,000 文字以内で入力してください"}
	}
	if len(input.InstructorComment) > 5000 {
		return nil, &ValidationError{Message: "師範コメントは 5,000 文字以内で入力してください"}
	}

	now := repository.Now()
	p := &model.Practice{
		ID:                uuid.New().String(),
		UserID:            input.UserID,
		DojoID:            input.DojoID,
		Date:              input.Date,
		HitRate:           input.HitRate,
		ArrowCount:        input.ArrowCount,
		Notes:             input.Notes,
		InstructorComment: input.InstructorComment,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.CreatePractice(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

// CreateVideoInput represents input for creating a video record.
type CreateVideoInput struct {
	UserID     string  `json:"userId"`
	PracticeID *string `json:"practiceId,omitempty"`
	FileName   string  `json:"fileName"`
	FileSize   int64   `json:"fileSize"`
	Duration   float64 `json:"duration"`
	MimeType   string  `json:"mimeType"`
	URL        string  `json:"url"`
}

// ListVideos returns videos, optionally filtered by userID.
func (s *Service) ListVideos(ctx context.Context, userID *string) ([]model.Video, error) {
	return s.repo.ListVideos(ctx, userID)
}

// GetVideo returns a video by ID.
func (s *Service) GetVideo(ctx context.Context, id string) (*model.Video, error) {
	v, err := s.repo.GetVideo(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "動画"}
		}
		return nil, err
	}
	return v, nil
}

// CreateVideo validates and creates a new video record.
func (s *Service) CreateVideo(ctx context.Context, input CreateVideoInput) (*model.Video, error) {
	if input.FileSize > 500*1024*1024 {
		return nil, &ValidationError{Message: "ファイルサイズは 500MB 以下にしてください"}
	}
	if input.Duration > 300 {
		return nil, &ValidationError{Message: "動画長は 5 分以下にしてください"}
	}

	now := repository.Now()
	v := &model.Video{
		ID:         uuid.New().String(),
		UserID:     input.UserID,
		PracticeID: input.PracticeID,
		FileName:   input.FileName,
		FileSize:   input.FileSize,
		Duration:   input.Duration,
		MimeType:   input.MimeType,
		Status:     model.VideoCompleted,
		URL:        input.URL,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateVideo(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

// ListAnalyses returns analyses, optionally filtered by userID.
func (s *Service) ListAnalyses(ctx context.Context, userID *string) ([]model.Analysis, error) {
	return s.repo.ListAnalyses(ctx, userID)
}

// GetAnalysis returns an analysis by ID.
func (s *Service) GetAnalysis(ctx context.Context, id string) (*model.Analysis, error) {
	a, err := s.repo.GetAnalysis(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "分析結果"}
		}
		return nil, err
	}
	return a, nil
}

// GetAnalysisByVideo returns an analysis for a given video.
func (s *Service) GetAnalysisByVideo(ctx context.Context, videoID string) (*model.Analysis, error) {
	a, err := s.repo.GetAnalysisByVideo(ctx, videoID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "分析結果"}
		}
		return nil, err
	}
	return a, nil
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

// CreateReservationInput represents input for creating a reservation.
type CreateReservationInput struct {
	DojoID     string `json:"dojoId"`
	UserID     string `json:"userId"`
	LaneNumber int    `json:"laneNumber"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

var hhmmPattern = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

// ListReservations returns reservations with optional filters.
func (s *Service) ListReservations(ctx context.Context, dojoID *string, date *string) ([]model.Reservation, error) {
	return s.repo.ListReservations(ctx, dojoID, date)
}

// GetReservation returns a reservation by ID.
func (s *Service) GetReservation(ctx context.Context, id string) (*model.Reservation, error) {
	r, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "予約"}
		}
		return nil, err
	}
	return r, nil
}

// CreateReservation validates and creates a new reservation.
func (s *Service) CreateReservation(ctx context.Context, input CreateReservationInput) (*model.Reservation, error) {
	if !hhmmPattern.MatchString(input.StartTime) || !hhmmPattern.MatchString(input.EndTime) {
		return nil, &ValidationError{Message: "時間は HH:mm 形式で入力してください"}
	}
	if input.EndTime <= input.StartTime {
		return nil, &ValidationError{Message: "終了時刻は開始時刻より後に設定してください"}
	}

	now := repository.Now()
	res := &model.Reservation{
		ID:         uuid.New().String(),
		DojoID:     input.DojoID,
		UserID:     input.UserID,
		LaneNumber: input.LaneNumber,
		Date:       input.Date,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateReservation(ctx, res); err != nil {
		if errors.Is(err, repository.ErrReservationConflict) {
			return nil, &ConflictError{Message: "同一的場・同一時間帯に既に予約があります"}
		}
		return nil, err
	}
	return res, nil
}

// DeleteReservation removes a reservation by ID.
func (s *Service) DeleteReservation(ctx context.Context, id string) error {
	if err := s.repo.DeleteReservation(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &NotFoundError{Resource: "予約"}
		}
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

// ListExamChecklists returns exam checklists, optionally filtered by userID.
func (s *Service) ListExamChecklists(ctx context.Context, userID *string) ([]model.ExamChecklist, error) {
	return s.repo.ListExamChecklists(ctx, userID)
}

// GetExamChecklist returns an exam checklist by ID.
func (s *Service) GetExamChecklist(ctx context.Context, id string) (*model.ExamChecklist, error) {
	c, err := s.repo.GetExamChecklist(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "審査チェックリスト"}
		}
		return nil, err
	}
	return c, nil
}

// ToggleChecklistItem toggles a specific item in a checklist.
func (s *Service) ToggleChecklistItem(ctx context.Context, checklistID string, itemID string) (*model.ExamChecklist, error) {
	c, err := s.repo.GetExamChecklist(ctx, checklistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &NotFoundError{Resource: "審査チェックリスト"}
		}
		return nil, err
	}

	found := false
	completedCount := 0
	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items[i].Completed = !item.Completed
			found = true
		}
		if c.Items[i].Completed {
			completedCount++
		}
	}
	if !found {
		return nil, &NotFoundError{Resource: "チェックリスト項目"}
	}

	if len(c.Items) > 0 {
		c.ProgressRate = (completedCount * 100) / len(c.Items)
	}
	c.UpdatedAt = repository.Now()

	if err := s.repo.UpdateExamChecklist(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// GetDashboardSummary returns the dashboard summary for a dojo.
func (s *Service) GetDashboardSummary(ctx context.Context, dojoID string) (*model.DashboardSummary, error) {
	today := time.Now().Format("2006-01-02")

	todayReservations, err := s.repo.ListReservationsByDojoAndDate(ctx, dojoID, today)
	if err != nil {
		return nil, fmt.Errorf("get today reservations: %w", err)
	}

	memberCount, err := s.repo.CountUsersByDojo(ctx, dojoID)
	if err != nil {
		return nil, fmt.Errorf("count members: %w", err)
	}

	if todayReservations == nil {
		todayReservations = []model.Reservation{}
	}

	return &model.DashboardSummary{
		TodayReservationCount: len(todayReservations),
		TotalMemberCount:      memberCount,
		TodayReservations:     todayReservations,
	}, nil
}
