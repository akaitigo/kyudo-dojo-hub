package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/ryusei/kyudo-dojo-hub/internal/model"
)

func ptr(s string) *string { return &s }

func newMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	t.Cleanup(mock.Close)
	return mock
}

func TestGetUser_Success(t *testing.T) {
	mock := newMock(t)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := mock.NewRows([]string{
		"id", "name", "email", "role", "dan", "shogo", "dojo_id", "joined_at", "created_at", "updated_at",
	}).AddRow("user-1", "田中太郎", "t@example.com", "practitioner", ptr("sandan"), nil, ptr("dojo-1"), "2020-04-01", now, now)
	mock.ExpectQuery("FROM users WHERE id =").WithArgs("user-1").WillReturnRows(rows)

	repo := New(mock)
	u, err := repo.GetUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if u.Name != "田中太郎" {
		t.Errorf("name = %q, want 田中太郎", u.Name)
	}
	if u.Dan == nil || *u.Dan != model.DanRank("sandan") {
		t.Errorf("dan = %v, want sandan", u.Dan)
	}
	if u.Shogo != nil {
		t.Errorf("shogo = %v, want nil", u.Shogo)
	}
	if u.DojoID == nil || *u.DojoID != "dojo-1" {
		t.Errorf("dojoID = %v, want dojo-1", u.DojoID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mock := newMock(t)
	mock.ExpectQuery("FROM users WHERE id =").WithArgs("missing").WillReturnError(pgx.ErrNoRows)

	repo := New(mock)
	_, err := repo.GetUser(context.Background(), "missing")
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestListUsers_ParsesRows(t *testing.T) {
	mock := newMock(t)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := mock.NewRows([]string{
		"id", "name", "email", "role", "dan", "shogo", "dojo_id", "joined_at", "created_at", "updated_at",
	}).
		AddRow("user-1", "田中", "a@example.com", "practitioner", ptr("sandan"), ptr("renshi"), ptr("dojo-1"), "2020-04-01", now, now).
		AddRow("user-2", "鈴木", "b@example.com", "manager", nil, nil, nil, "2021-04-01", now, now)
	mock.ExpectQuery("FROM users ORDER BY name").WillReturnRows(rows)

	repo := New(mock)
	users, err := repo.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("len = %d, want 2", len(users))
	}
	if users[1].Dan != nil {
		t.Errorf("user-2 dan = %v, want nil", users[1].Dan)
	}
}

func TestCreatePractice_Exec(t *testing.T) {
	mock := newMock(t)
	now := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	p := &model.Practice{
		ID: "p-1", UserID: "user-1", Date: "2026-04-01", HitRate: 70, ArrowCount: 20,
		Notes: "n", InstructorComment: "", CreatedAt: now, UpdatedAt: now,
	}
	mock.ExpectExec("INSERT INTO practices").
		WithArgs(p.ID, p.UserID, p.DojoID, p.Date, p.HitRate, p.ArrowCount, p.Notes, p.InstructorComment, p.CreatedAt, p.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	repo := New(mock)
	if err := repo.CreatePractice(context.Background(), p); err != nil {
		t.Fatalf("CreatePractice: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateVideo_Exec(t *testing.T) {
	mock := newMock(t)
	now := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	v := &model.Video{
		ID: "v-1", UserID: "user-1", FileName: "a.mp4", FileSize: 1024, Duration: 12.5,
		MimeType: "video/mp4", Status: model.VideoCompleted, URL: "", CreatedAt: now, UpdatedAt: now,
	}
	mock.ExpectExec("INSERT INTO videos").
		WithArgs(v.ID, v.UserID, v.PracticeID, v.FileName, v.FileSize, v.Duration, v.MimeType, v.Status, v.URL, v.CreatedAt, v.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	repo := New(mock)
	if err := repo.CreateVideo(context.Background(), v); err != nil {
		t.Fatalf("CreateVideo: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func reservationRows(mock pgxmock.PgxPoolIface) *pgxmock.Rows {
	now := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	return mock.NewRows([]string{
		"id", "dojo_id", "user_id", "lane_number", "date", "start_time", "end_time", "created_at", "updated_at",
	}).AddRow("res-1", "dojo-1", "user-1", 2, "2026-04-10", "10:00", "11:00", now, now)
}

// TestListReservations_WithFilters verifies the dynamic WHERE builder (#29):
// both filters produce sequential placeholders $1/$2 with matching args.
func TestListReservations_WithFilters(t *testing.T) {
	mock := newMock(t)
	mock.ExpectQuery(`FROM reservations WHERE dojo_id = \$1 AND date = \$2 ORDER BY date, start_time`).
		WithArgs("dojo-1", "2026-04-10").
		WillReturnRows(reservationRows(mock))

	repo := New(mock)
	got, err := repo.ListReservations(context.Background(), ptr("dojo-1"), ptr("2026-04-10"))
	if err != nil {
		t.Fatalf("ListReservations: %v", err)
	}
	if len(got) != 1 || got[0].ID != "res-1" {
		t.Fatalf("unexpected result: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// TestListReservations_DateOnly verifies a single filter uses $1 (not $2).
func TestListReservations_DateOnly(t *testing.T) {
	mock := newMock(t)
	mock.ExpectQuery(`FROM reservations WHERE date = \$1 ORDER BY date, start_time`).
		WithArgs("2026-04-10").
		WillReturnRows(reservationRows(mock))

	repo := New(mock)
	if _, err := repo.ListReservations(context.Background(), nil, ptr("2026-04-10")); err != nil {
		t.Fatalf("ListReservations: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// TestListReservations_NoFilters verifies no WHERE clause is emitted.
func TestListReservations_NoFilters(t *testing.T) {
	mock := newMock(t)
	mock.ExpectQuery(`FROM reservations ORDER BY date, start_time`).
		WithArgs().
		WillReturnRows(reservationRows(mock))

	repo := New(mock)
	if _, err := repo.ListReservations(context.Background(), nil, nil); err != nil {
		t.Fatalf("ListReservations: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteReservation_Success(t *testing.T) {
	mock := newMock(t)
	mock.ExpectExec("DELETE FROM reservations").
		WithArgs("res-1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	repo := New(mock)
	if err := repo.DeleteReservation(context.Background(), "res-1"); err != nil {
		t.Fatalf("DeleteReservation: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteReservation_NotFound(t *testing.T) {
	mock := newMock(t)
	mock.ExpectExec("DELETE FROM reservations").
		WithArgs("missing").
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	repo := New(mock)
	err := repo.DeleteReservation(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCountUsersByDojo(t *testing.T) {
	mock := newMock(t)
	mock.ExpectQuery("SELECT COUNT").
		WithArgs("dojo-1").
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(5))

	repo := New(mock)
	n, err := repo.CountUsersByDojo(context.Background(), "dojo-1")
	if err != nil {
		t.Fatalf("CountUsersByDojo: %v", err)
	}
	if n != 5 {
		t.Errorf("count = %d, want 5", n)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
