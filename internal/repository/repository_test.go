package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/ryusei/kyudo-dojo-hub/internal/model"
)

func sampleReservation() *model.Reservation {
	now := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	return &model.Reservation{
		ID:         "res-1",
		DojoID:     "dojo-1",
		UserID:     "user-1",
		LaneNumber: 2,
		Date:       "2026-04-10",
		StartTime:  "10:00",
		EndTime:    "11:00",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// countRows returns a single-column "count" result set for the overlap check.
func countRows(mock pgxmock.PgxPoolIface, n int) *pgxmock.Rows {
	return mock.NewRows([]string{"count"}).AddRow(n)
}

func TestCreateReservation_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	defer mock.Close()

	res := sampleReservation()
	mock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.Serializable})
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(res.DojoID, res.LaneNumber, res.Date, res.EndTime, res.StartTime).
		WillReturnRows(countRows(mock, 0))
	mock.ExpectExec("INSERT INTO reservations").
		WithArgs(res.ID, res.DojoID, res.UserID, res.LaneNumber, res.Date, res.StartTime, res.EndTime, res.CreatedAt, res.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	repo := New(mock)
	if err := repo.CreateReservation(context.Background(), res); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateReservation_OverlapConflict(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	defer mock.Close()

	res := sampleReservation()
	// 既存予約と時間帯が重なる → COUNT が 1 を返し、INSERT せずロールバックする。
	mock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.Serializable})
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(res.DojoID, res.LaneNumber, res.Date, res.EndTime, res.StartTime).
		WillReturnRows(countRows(mock, 1))
	mock.ExpectRollback()

	repo := New(mock)
	err = repo.CreateReservation(context.Background(), res)
	if !errors.Is(err, ErrReservationConflict) {
		t.Fatalf("expected ErrReservationConflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateReservation_SerializationFailureIsConflict(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	defer mock.Close()

	res := sampleReservation()
	// 並行トランザクションが競合し、Commit がシリアライゼーション失敗(40001)に
	// なるケース。競合として扱われることを検証する。
	mock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.Serializable})
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(res.DojoID, res.LaneNumber, res.Date, res.EndTime, res.StartTime).
		WillReturnRows(countRows(mock, 0))
	mock.ExpectExec("INSERT INTO reservations").
		WithArgs(res.ID, res.DojoID, res.UserID, res.LaneNumber, res.Date, res.StartTime, res.EndTime, res.CreatedAt, res.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit().WillReturnError(&pgconn.PgError{Code: "40001"})
	mock.ExpectRollback()

	repo := New(mock)
	err = repo.CreateReservation(context.Background(), res)
	if !errors.Is(err, ErrReservationConflict) {
		t.Fatalf("expected ErrReservationConflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateReservation_UniqueViolationIsConflict(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	defer mock.Close()

	res := sampleReservation()
	// 同一 lane/date/start_time の同時 INSERT で一意制約(23505)違反になるケース。
	mock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.Serializable})
	mock.ExpectQuery("SELECT COUNT").
		WithArgs(res.DojoID, res.LaneNumber, res.Date, res.EndTime, res.StartTime).
		WillReturnRows(countRows(mock, 0))
	mock.ExpectExec("INSERT INTO reservations").
		WithArgs(res.ID, res.DojoID, res.UserID, res.LaneNumber, res.Date, res.StartTime, res.EndTime, res.CreatedAt, res.UpdatedAt).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	mock.ExpectRollback()

	repo := New(mock)
	err = repo.CreateReservation(context.Background(), res)
	if !errors.Is(err, ErrReservationConflict) {
		t.Fatalf("expected ErrReservationConflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
