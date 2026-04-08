// Package store provides a thread-safe in-memory data store for the kyudo-dojo-hub backend.
package store

import (
	"crypto/rand"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/model"
)

// Store is a thread-safe in-memory data store.
type Store struct {
	mu             sync.RWMutex
	users          []model.User
	dojos          []model.Dojo
	practices      []model.Practice
	videos         []model.Video
	analyses       []model.Analysis
	reservations   []model.Reservation
	examChecklists []model.ExamChecklist
}

// New creates a Store pre-populated with seed data.
func New() *Store {
	s := &Store{}
	s.seed()
	return s
}

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), RandomString(7))
}

// RandomString generates a cryptographically random string of length n
// using lowercase letters and digits. It is safe for ID generation.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read never returns an error on supported platforms,
		// but handle it defensively.
		panic(fmt.Sprintf("crypto/rand.Read failed: %v", err))
	}
	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}

func strPtr(s string) *string { return &s }

func danPtr(d model.DanRank) *model.DanRank       { return &d }
func shogoPtr(s model.ShogoTitle) *model.ShogoTitle { return &s }

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(fmt.Sprintf("invalid time format: %s", s))
	}
	return t
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// GetUsers returns all users.
func (s *Store) GetUsers() []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.User, len(s.users))
	copy(result, s.users)
	return result
}

// GetUser returns a user by ID.
func (s *Store) GetUser(id string) (model.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			return u, true
		}
	}
	return model.User{}, false
}

// GetUsersByDojo returns all users belonging to a specific dojo.
func (s *Store) GetUsersByDojo(dojoID string) []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.User
	for _, u := range s.users {
		if u.DojoID != nil && *u.DojoID == dojoID {
			result = append(result, u)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

// GetDojos returns all dojos.
func (s *Store) GetDojos() []model.Dojo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Dojo, len(s.dojos))
	copy(result, s.dojos)
	return result
}

// GetDojo returns a dojo by ID.
func (s *Store) GetDojo(id string) (model.Dojo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, d := range s.dojos {
		if d.ID == id {
			return d, true
		}
	}
	return model.Dojo{}, false
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

// GetPractices returns all practices, optionally filtered by user ID, sorted by date descending.
func (s *Store) GetPractices(userID string) []model.Practice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Practice
	for _, p := range s.practices {
		if userID == "" || p.UserID == userID {
			result = append(result, p)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})
	return result
}

// GetPractice returns a single practice by ID.
func (s *Store) GetPractice(id string) (model.Practice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.practices {
		if p.ID == id {
			return p, true
		}
	}
	return model.Practice{}, false
}

// CreatePractice validates and creates a new practice record.
func (s *Store) CreatePractice(input model.CreatePracticeInput) (model.Practice, error) {
	if input.HitRate < 0 || input.HitRate > 100 {
		return model.Practice{}, fmt.Errorf("的中率は 0〜100 の範囲で入力してください")
	}
	if input.ArrowCount < 1 || input.ArrowCount > 1000 {
		return model.Practice{}, fmt.Errorf("矢数は 1〜1000 の範囲で入力してください")
	}
	if len([]rune(input.Notes)) > 5000 {
		return model.Practice{}, fmt.Errorf("気づきは 5,000 文字以内で入力してください")
	}
	if len([]rune(input.InstructorComment)) > 5000 {
		return model.Practice{}, fmt.Errorf("師範コメントは 5,000 文字以内で入力してください")
	}
	if input.UserID == "" {
		return model.Practice{}, fmt.Errorf("ユーザーIDは必須です")
	}
	if input.Date == "" {
		return model.Practice{}, fmt.Errorf("日付は必須です")
	}

	now := time.Now()
	practice := model.Practice{
		ID:                fmt.Sprintf("practice-%s", generateID()),
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

	s.mu.Lock()
	defer s.mu.Unlock()
	s.practices = append([]model.Practice{practice}, s.practices...)
	return practice, nil
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

// GetVideos returns all videos, optionally filtered by user ID.
func (s *Store) GetVideos(userID string) []model.Video {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Video
	for _, v := range s.videos {
		if userID == "" || v.UserID == userID {
			result = append(result, v)
		}
	}
	return result
}

// GetVideo returns a single video by ID.
func (s *Store) GetVideo(id string) (model.Video, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.videos {
		if v.ID == id {
			return v, true
		}
	}
	return model.Video{}, false
}

// CreateVideo validates and creates a new video record.
func (s *Store) CreateVideo(input model.CreateVideoInput) (model.Video, error) {
	if input.FileSize > 500*1024*1024 {
		return model.Video{}, fmt.Errorf("ファイルサイズは 500MB 以下にしてください")
	}
	if input.Duration > 300 {
		return model.Video{}, fmt.Errorf("動画長は 5 分以下にしてください")
	}
	allowedTypes := map[string]bool{
		"video/mp4":       true,
		"video/quicktime": true,
		"video/webm":      true,
	}
	if !allowedTypes[input.MimeType] {
		return model.Video{}, fmt.Errorf("mp4, mov, webm 形式のみ対応しています")
	}
	if input.UserID == "" {
		return model.Video{}, fmt.Errorf("ユーザーIDは必須です")
	}

	now := time.Now()
	video := model.Video{
		ID:         fmt.Sprintf("video-%s", generateID()),
		UserID:     input.UserID,
		PracticeID: input.PracticeID,
		FileName:   input.FileName,
		FileSize:   input.FileSize,
		Duration:   input.Duration,
		MimeType:   input.MimeType,
		Status:     model.StatusCompleted,
		URL:        input.URL,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.videos = append([]model.Video{video}, s.videos...)
	return video, nil
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

// GetAnalyses returns all analyses, optionally filtered by user ID.
func (s *Store) GetAnalyses(userID string) []model.Analysis {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Analysis
	for _, a := range s.analyses {
		if userID == "" || a.UserID == userID {
			result = append(result, a)
		}
	}
	return result
}

// GetAnalysis returns a single analysis by ID.
func (s *Store) GetAnalysis(id string) (model.Analysis, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.analyses {
		if a.ID == id {
			return a, true
		}
	}
	return model.Analysis{}, false
}

// GetAnalysisByVideo returns the analysis for a given video.
func (s *Store) GetAnalysisByVideo(videoID string) (model.Analysis, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.analyses {
		if a.VideoID == videoID {
			return a, true
		}
	}
	return model.Analysis{}, false
}

// AddAnalysis adds a new analysis result.
func (s *Store) AddAnalysis(a model.Analysis) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.analyses = append(s.analyses, a)
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

// GetReservations returns all reservations, optionally filtered by dojo ID and date.
func (s *Store) GetReservations(dojoID, date string) []model.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Reservation
	for _, r := range s.reservations {
		if dojoID != "" && r.DojoID != dojoID {
			continue
		}
		if date != "" && r.Date != date {
			continue
		}
		result = append(result, r)
	}
	return result
}

// GetReservation returns a single reservation by ID.
func (s *Store) GetReservation(id string) (model.Reservation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.reservations {
		if r.ID == id {
			return r, true
		}
	}
	return model.Reservation{}, false
}

// CreateReservation validates and creates a new reservation.
func (s *Store) CreateReservation(input model.CreateReservationInput) (model.Reservation, error) {
	hhmmPattern := func(t string) bool {
		if len(t) != 5 {
			return false
		}
		if t[2] != ':' {
			return false
		}
		h := (t[0]-'0')*10 + (t[1] - '0')
		m := (t[3]-'0')*10 + (t[4] - '0')
		return h <= 23 && m <= 59
	}

	if !hhmmPattern(input.StartTime) || !hhmmPattern(input.EndTime) {
		return model.Reservation{}, fmt.Errorf("時間は HH:mm 形式で入力してください")
	}
	if input.EndTime <= input.StartTime {
		return model.Reservation{}, fmt.Errorf("終了時刻は開始時刻より後に設定してください")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for conflicting reservations.
	for _, r := range s.reservations {
		if r.DojoID == input.DojoID &&
			r.LaneNumber == input.LaneNumber &&
			r.Date == input.Date &&
			r.StartTime < input.EndTime &&
			input.StartTime < r.EndTime {
			return model.Reservation{}, fmt.Errorf("同一的場・同一時間帯に既に予約があります")
		}
	}

	now := time.Now()
	reservation := model.Reservation{
		ID:         fmt.Sprintf("res-%s", generateID()),
		DojoID:     input.DojoID,
		UserID:     input.UserID,
		LaneNumber: input.LaneNumber,
		Date:       input.Date,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	s.reservations = append([]model.Reservation{reservation}, s.reservations...)
	return reservation, nil
}

// DeleteReservation removes a reservation by ID.
func (s *Store) DeleteReservation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.reservations {
		if r.ID == id {
			s.reservations = append(s.reservations[:i], s.reservations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("予約が見つかりません")
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

// GetExamChecklists returns all exam checklists, optionally filtered by user ID.
func (s *Store) GetExamChecklists(userID string) []model.ExamChecklist {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.ExamChecklist
	for _, c := range s.examChecklists {
		if userID == "" || c.UserID == userID {
			result = append(result, c)
		}
	}
	return result
}

// GetExamChecklist returns a single exam checklist by ID.
func (s *Store) GetExamChecklist(id string) (model.ExamChecklist, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.examChecklists {
		if c.ID == id {
			return c, true
		}
	}
	return model.ExamChecklist{}, false
}

// ToggleChecklistItem toggles the completed state of a checklist item.
func (s *Store) ToggleChecklistItem(checklistID, itemID string) (model.ExamChecklist, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, c := range s.examChecklists {
		if c.ID != checklistID {
			continue
		}
		found := false
		items := make([]model.ExamChecklistItem, len(c.Items))
		copy(items, c.Items)
		for j := range items {
			if items[j].ID == itemID {
				items[j].Completed = !items[j].Completed
				found = true
			}
		}
		if !found {
			return model.ExamChecklist{}, fmt.Errorf("チェック項目が見つかりません")
		}

		completedCount := 0
		for _, item := range items {
			if item.Completed {
				completedCount++
			}
		}
		progressRate := 0
		if len(items) > 0 {
			progressRate = (completedCount * 100) / len(items)
		}

		s.examChecklists[i].Items = items
		s.examChecklists[i].ProgressRate = progressRate
		s.examChecklists[i].UpdatedAt = time.Now()
		return s.examChecklists[i], nil
	}
	return model.ExamChecklist{}, fmt.Errorf("審査チェックリストが見つかりません")
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// GetDashboardSummary returns a summary for the dojo dashboard.
func (s *Store) GetDashboardSummary(dojoID string) model.DashboardSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	today := time.Now().Format("2006-01-02")
	var todayReservations []model.Reservation
	for _, r := range s.reservations {
		if r.DojoID == dojoID && r.Date == today {
			todayReservations = append(todayReservations, r)
		}
	}

	memberCount := 0
	for _, u := range s.users {
		if u.DojoID != nil && *u.DojoID == dojoID {
			memberCount++
		}
	}

	if todayReservations == nil {
		todayReservations = []model.Reservation{}
	}

	return model.DashboardSummary{
		TodayReservationCount: len(todayReservations),
		TotalMemberCount:      memberCount,
		TodayReservations:     todayReservations,
	}
}
