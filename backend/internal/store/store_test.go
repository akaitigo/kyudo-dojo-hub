package store_test

import (
	"testing"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/model"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/store"
)

func newTestStore() *store.Store {
	return store.New()
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

func TestGetUsers(t *testing.T) {
	s := newTestStore()
	users := s.GetUsers()
	if len(users) < 10 {
		t.Fatalf("expected at least 10 users, got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	s := newTestStore()
	user, ok := s.GetUser("user-001")
	if !ok {
		t.Fatal("expected to find user-001")
	}
	if user.Name != "田中太郎" {
		t.Fatalf("expected name 田中太郎, got %s", user.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	s := newTestStore()
	_, ok := s.GetUser("nonexistent")
	if ok {
		t.Fatal("expected not to find nonexistent user")
	}
}

func TestGetUsersByDojo(t *testing.T) {
	s := newTestStore()
	users := s.GetUsersByDojo("dojo-001")
	if len(users) == 0 {
		t.Fatal("expected at least 1 user for dojo-001")
	}
	for _, u := range users {
		if u.DojoID == nil || *u.DojoID != "dojo-001" {
			t.Fatalf("user %s has unexpected dojoId %v", u.ID, u.DojoID)
		}
	}
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

func TestGetDojos(t *testing.T) {
	s := newTestStore()
	dojos := s.GetDojos()
	if len(dojos) != 2 {
		t.Fatalf("expected 2 dojos, got %d", len(dojos))
	}
}

func TestGetDojo(t *testing.T) {
	s := newTestStore()
	dojo, ok := s.GetDojo("dojo-001")
	if !ok {
		t.Fatal("expected to find dojo-001")
	}
	if dojo.Name != "東京弓道場" {
		t.Fatalf("expected name 東京弓道場, got %s", dojo.Name)
	}
}

func TestGetDojo_NotFound(t *testing.T) {
	s := newTestStore()
	_, ok := s.GetDojo("nonexistent")
	if ok {
		t.Fatal("expected not to find nonexistent dojo")
	}
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

func TestGetPractices_All(t *testing.T) {
	s := newTestStore()
	practices := s.GetPractices("")
	if len(practices) < 10 {
		t.Fatalf("expected at least 10 practices, got %d", len(practices))
	}
}

func TestGetPractices_ByUser(t *testing.T) {
	s := newTestStore()
	practices := s.GetPractices("user-001")
	if len(practices) == 0 {
		t.Fatal("expected at least 1 practice for user-001")
	}
	for _, p := range practices {
		if p.UserID != "user-001" {
			t.Fatalf("expected userId user-001, got %s", p.UserID)
		}
	}
}

func TestGetPractices_SortedByDateDesc(t *testing.T) {
	s := newTestStore()
	practices := s.GetPractices("user-001")
	if len(practices) < 2 {
		t.Skip("not enough practices to check sort order")
	}
	for i := 0; i < len(practices)-1; i++ {
		if practices[i].Date < practices[i+1].Date {
			t.Fatalf("expected date descending order, got %s before %s", practices[i].Date, practices[i+1].Date)
		}
	}
}

func TestGetPractice(t *testing.T) {
	s := newTestStore()
	practice, ok := s.GetPractice("practice-001")
	if !ok {
		t.Fatal("expected to find practice-001")
	}
	if practice.HitRate != 65 {
		t.Fatalf("expected hitRate 65, got %d", practice.HitRate)
	}
}

func TestCreatePractice_Success(t *testing.T) {
	s := newTestStore()
	dojoID := "dojo-001"
	input := model.CreatePracticeInput{
		UserID:            "user-001",
		DojoID:            &dojoID,
		Date:              "2026-03-30",
		HitRate:           60,
		ArrowCount:        36,
		Notes:             "テスト稽古",
		InstructorComment: "",
	}
	practice, err := s.CreatePractice(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if practice.HitRate != 60 {
		t.Fatalf("expected hitRate 60, got %d", practice.HitRate)
	}
	if practice.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestCreatePractice_HitRateValidation(t *testing.T) {
	s := newTestStore()
	input := model.CreatePracticeInput{
		UserID:     "user-001",
		Date:       "2026-03-30",
		HitRate:    101,
		ArrowCount: 36,
	}
	_, err := s.CreatePractice(input)
	if err == nil {
		t.Fatal("expected validation error for hitRate 101")
	}
}

func TestCreatePractice_HitRateNegative(t *testing.T) {
	s := newTestStore()
	input := model.CreatePracticeInput{
		UserID:     "user-001",
		Date:       "2026-03-30",
		HitRate:    -1,
		ArrowCount: 36,
	}
	_, err := s.CreatePractice(input)
	if err == nil {
		t.Fatal("expected validation error for hitRate -1")
	}
}

func TestCreatePractice_ArrowCountValidation(t *testing.T) {
	s := newTestStore()
	input := model.CreatePracticeInput{
		UserID:     "user-001",
		Date:       "2026-03-30",
		HitRate:    50,
		ArrowCount: 0,
	}
	_, err := s.CreatePractice(input)
	if err == nil {
		t.Fatal("expected validation error for arrowCount 0")
	}
}

func TestCreatePractice_NotesTooLong(t *testing.T) {
	s := newTestStore()
	longNotes := make([]rune, 5001)
	for i := range longNotes {
		longNotes[i] = 'あ'
	}
	input := model.CreatePracticeInput{
		UserID:     "user-001",
		Date:       "2026-03-30",
		HitRate:    50,
		ArrowCount: 10,
		Notes:      string(longNotes),
	}
	_, err := s.CreatePractice(input)
	if err == nil {
		t.Fatal("expected validation error for notes exceeding 5000 chars")
	}
}

func TestCreatePractice_EmptyUserID(t *testing.T) {
	s := newTestStore()
	input := model.CreatePracticeInput{
		Date:       "2026-03-30",
		HitRate:    50,
		ArrowCount: 10,
	}
	_, err := s.CreatePractice(input)
	if err == nil {
		t.Fatal("expected validation error for empty userID")
	}
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

func TestGetVideos(t *testing.T) {
	s := newTestStore()
	videos := s.GetVideos("")
	if len(videos) < 3 {
		t.Fatalf("expected at least 3 videos, got %d", len(videos))
	}
}

func TestGetVideo(t *testing.T) {
	s := newTestStore()
	video, ok := s.GetVideo("video-001")
	if !ok {
		t.Fatal("expected to find video-001")
	}
	if video.FileName != "tanaka_20260328.mp4" {
		t.Fatalf("expected fileName tanaka_20260328.mp4, got %s", video.FileName)
	}
}

func TestCreateVideo_Success(t *testing.T) {
	s := newTestStore()
	input := model.CreateVideoInput{
		UserID:   "user-001",
		FileName: "test.mp4",
		FileSize: 1024 * 1024,
		Duration: 30,
		MimeType: "video/mp4",
		URL:      "blob:test",
	}
	video, err := s.CreateVideo(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.Status != model.StatusCompleted {
		t.Fatalf("expected status completed, got %s", video.Status)
	}
}

func TestCreateVideo_FileSizeTooLarge(t *testing.T) {
	s := newTestStore()
	input := model.CreateVideoInput{
		UserID:   "user-001",
		FileName: "test.mp4",
		FileSize: 501 * 1024 * 1024,
		Duration: 30,
		MimeType: "video/mp4",
	}
	_, err := s.CreateVideo(input)
	if err == nil {
		t.Fatal("expected validation error for file size exceeding 500MB")
	}
}

func TestCreateVideo_DurationTooLong(t *testing.T) {
	s := newTestStore()
	input := model.CreateVideoInput{
		UserID:   "user-001",
		FileName: "test.mp4",
		FileSize: 1024,
		Duration: 301,
		MimeType: "video/mp4",
	}
	_, err := s.CreateVideo(input)
	if err == nil {
		t.Fatal("expected validation error for duration exceeding 300s")
	}
}

func TestCreateVideo_InvalidMimeType(t *testing.T) {
	s := newTestStore()
	input := model.CreateVideoInput{
		UserID:   "user-001",
		FileName: "test.avi",
		FileSize: 1024,
		Duration: 30,
		MimeType: "video/x-msvideo",
	}
	_, err := s.CreateVideo(input)
	if err == nil {
		t.Fatal("expected validation error for unsupported mime type")
	}
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

func TestGetAnalyses(t *testing.T) {
	s := newTestStore()
	analyses := s.GetAnalyses("")
	if len(analyses) < 2 {
		t.Fatalf("expected at least 2 analyses, got %d", len(analyses))
	}
}

func TestGetAnalysisByVideo(t *testing.T) {
	s := newTestStore()
	analysis, ok := s.GetAnalysisByVideo("video-001")
	if !ok {
		t.Fatal("expected to find analysis for video-001")
	}
	if analysis.OverallScore != 73 {
		t.Fatalf("expected overallScore 73, got %d", analysis.OverallScore)
	}
}

func TestGetAnalysisByVideo_NotFound(t *testing.T) {
	s := newTestStore()
	_, ok := s.GetAnalysisByVideo("nonexistent")
	if ok {
		t.Fatal("expected not to find analysis for nonexistent video")
	}
}

func TestGetAnalysis_HasAllPhaseScores(t *testing.T) {
	s := newTestStore()
	analysis, ok := s.GetAnalysisByVideo("video-001")
	if !ok {
		t.Fatal("expected to find analysis for video-001")
	}
	scores := analysis.Scores
	if scores.Ashibumi == 0 && scores.Dozukuri == 0 && scores.Yugamae == 0 &&
		scores.Uchiokoshi == 0 && scores.Hikiwake == 0 && scores.Kai == 0 &&
		scores.Hanare == 0 && scores.Zanshin == 0 {
		t.Fatal("expected non-zero scores for at least one phase")
	}
	if len(analysis.Phases) != 8 {
		t.Fatalf("expected 8 phases, got %d", len(analysis.Phases))
	}
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

func TestGetReservations_All(t *testing.T) {
	s := newTestStore()
	reservations := s.GetReservations("", "")
	if len(reservations) < 10 {
		t.Fatalf("expected at least 10 reservations, got %d", len(reservations))
	}
}

func TestGetReservations_ByDojoAndDate(t *testing.T) {
	s := newTestStore()
	reservations := s.GetReservations("dojo-001", "2026-03-30")
	if len(reservations) == 0 {
		t.Fatal("expected at least 1 reservation for dojo-001 on 2026-03-30")
	}
	for _, r := range reservations {
		if r.DojoID != "dojo-001" {
			t.Fatalf("expected dojoId dojo-001, got %s", r.DojoID)
		}
		if r.Date != "2026-03-30" {
			t.Fatalf("expected date 2026-03-30, got %s", r.Date)
		}
	}
}

func TestCreateReservation_Success(t *testing.T) {
	s := newTestStore()
	input := model.CreateReservationInput{
		DojoID:     "dojo-001",
		UserID:     "user-001",
		LaneNumber: 6,
		Date:       "2099-12-31",
		StartTime:  "09:00",
		EndTime:    "10:00",
	}
	reservation, err := s.CreateReservation(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reservation.LaneNumber != 6 {
		t.Fatalf("expected laneNumber 6, got %d", reservation.LaneNumber)
	}
}

func TestCreateReservation_InvalidTimeFormat(t *testing.T) {
	s := newTestStore()
	input := model.CreateReservationInput{
		DojoID:     "dojo-001",
		UserID:     "user-001",
		LaneNumber: 1,
		Date:       "2099-12-31",
		StartTime:  "9:00",
		EndTime:    "10:00",
	}
	_, err := s.CreateReservation(input)
	if err == nil {
		t.Fatal("expected validation error for invalid time format")
	}
}

func TestCreateReservation_EndBeforeStart(t *testing.T) {
	s := newTestStore()
	input := model.CreateReservationInput{
		DojoID:     "dojo-001",
		UserID:     "user-001",
		LaneNumber: 1,
		Date:       "2099-12-31",
		StartTime:  "10:00",
		EndTime:    "09:00",
	}
	_, err := s.CreateReservation(input)
	if err == nil {
		t.Fatal("expected validation error for endTime before startTime")
	}
}

func TestCreateReservation_OverlapRejected(t *testing.T) {
	s := newTestStore()
	base := model.CreateReservationInput{
		DojoID:     "dojo-overlap-test",
		UserID:     "user-001",
		LaneNumber: 99,
		Date:       "2099-12-31",
		StartTime:  "10:00",
		EndTime:    "11:00",
	}

	_, err := s.CreateReservation(base)
	if err != nil {
		t.Fatalf("unexpected error creating first reservation: %v", err)
	}

	// Partially overlapping reservation
	overlap := base
	overlap.StartTime = "10:30"
	overlap.EndTime = "11:30"
	_, err = s.CreateReservation(overlap)
	if err == nil {
		t.Fatal("expected overlap rejection for 10:30-11:30")
	}
}

func TestCreateReservation_AdjacentAllowed(t *testing.T) {
	s := newTestStore()
	base := model.CreateReservationInput{
		DojoID:     "dojo-adjacent-test",
		UserID:     "user-001",
		LaneNumber: 99,
		Date:       "2099-12-31",
		StartTime:  "10:00",
		EndTime:    "11:00",
	}

	_, err := s.CreateReservation(base)
	if err != nil {
		t.Fatalf("unexpected error creating first reservation: %v", err)
	}

	adjacent := base
	adjacent.StartTime = "11:00"
	adjacent.EndTime = "12:00"
	_, err = s.CreateReservation(adjacent)
	if err != nil {
		t.Fatalf("adjacent reservation should be allowed: %v", err)
	}
}

func TestCreateReservation_DifferentLaneAllowed(t *testing.T) {
	s := newTestStore()
	base := model.CreateReservationInput{
		DojoID:     "dojo-lane-test",
		UserID:     "user-001",
		LaneNumber: 1,
		Date:       "2099-12-31",
		StartTime:  "10:00",
		EndTime:    "11:00",
	}

	_, err := s.CreateReservation(base)
	if err != nil {
		t.Fatalf("unexpected error creating first reservation: %v", err)
	}

	differentLane := base
	differentLane.LaneNumber = 2
	_, err = s.CreateReservation(differentLane)
	if err != nil {
		t.Fatalf("different lane should be allowed: %v", err)
	}
}

func TestDeleteReservation(t *testing.T) {
	s := newTestStore()
	err := s.DeleteReservation("res-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not be found after deletion
	_, ok := s.GetReservation("res-001")
	if ok {
		t.Fatal("expected res-001 to be deleted")
	}
}

func TestDeleteReservation_NotFound(t *testing.T) {
	s := newTestStore()
	err := s.DeleteReservation("nonexistent")
	if err == nil {
		t.Fatal("expected error for deleting nonexistent reservation")
	}
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

func TestGetExamChecklists(t *testing.T) {
	s := newTestStore()
	checklists := s.GetExamChecklists("")
	if len(checklists) < 2 {
		t.Fatalf("expected at least 2 checklists, got %d", len(checklists))
	}
}

func TestGetExamChecklists_ByUser(t *testing.T) {
	s := newTestStore()
	checklists := s.GetExamChecklists("user-001")
	if len(checklists) == 0 {
		t.Fatal("expected at least 1 checklist for user-001")
	}
	for _, c := range checklists {
		if c.UserID != "user-001" {
			t.Fatalf("expected userId user-001, got %s", c.UserID)
		}
	}
}

func TestGetExamChecklist(t *testing.T) {
	s := newTestStore()
	checklist, ok := s.GetExamChecklist("exam-001")
	if !ok {
		t.Fatal("expected to find exam-001")
	}
	if checklist.TargetDan != model.DanYondan {
		t.Fatalf("expected targetDan yondan, got %s", checklist.TargetDan)
	}
}

func TestToggleChecklistItem(t *testing.T) {
	s := newTestStore()

	before, ok := s.GetExamChecklist("exam-002")
	if !ok {
		t.Fatal("expected to find exam-002")
	}
	var itemBefore *model.ExamChecklistItem
	for _, item := range before.Items {
		if item.ID == "item-011" {
			itemCopy := item
			itemBefore = &itemCopy
			break
		}
	}
	if itemBefore == nil {
		t.Fatal("expected to find item-011")
	}

	after, err := s.ToggleChecklistItem("exam-002", "item-011")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var itemAfter *model.ExamChecklistItem
	for _, item := range after.Items {
		if item.ID == "item-011" {
			itemCopy := item
			itemAfter = &itemCopy
			break
		}
	}
	if itemAfter == nil {
		t.Fatal("expected to find item-011 after toggle")
	}

	if itemAfter.Completed == itemBefore.Completed {
		t.Fatal("expected completed state to be toggled")
	}
}

func TestToggleChecklistItem_ChecklistNotFound(t *testing.T) {
	s := newTestStore()
	_, err := s.ToggleChecklistItem("nonexistent", "item-001")
	if err == nil {
		t.Fatal("expected error for nonexistent checklist")
	}
}

func TestToggleChecklistItem_ItemNotFound(t *testing.T) {
	s := newTestStore()
	_, err := s.ToggleChecklistItem("exam-001", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent item")
	}
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

func TestGetDashboardSummary(t *testing.T) {
	s := newTestStore()
	summary := s.GetDashboardSummary("dojo-001")
	if summary.TotalMemberCount == 0 {
		t.Fatal("expected non-zero member count for dojo-001")
	}
	if summary.TodayReservations == nil {
		t.Fatal("expected non-nil todayReservations slice")
	}
}
