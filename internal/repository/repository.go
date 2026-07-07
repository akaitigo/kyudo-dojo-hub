// Package repository provides PostgreSQL-backed data access.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ryusei/kyudo-dojo-hub/internal/model"
)

// PgxDB is the subset of *pgxpool.Pool used by the repository. Depending on an
// interface (rather than the concrete pool) lets tests inject a mock such as
// pgxmock while production code passes a real *pgxpool.Pool.
type PgxDB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

// Repo wraps a database connection pool for data access.
type Repo struct {
	pool PgxDB
}

// New creates a new Repo from a connection pool.
func New(pool PgxDB) *Repo {
	return &Repo{pool: pool}
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// ListUsers returns all users.
func (r *Repo) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, email, role, dan, shogo, dojo_id, joined_at, created_at, updated_at
		 FROM users ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	return collectUsers(rows)
}

// GetUser returns a user by ID.
func (r *Repo) GetUser(ctx context.Context, id string) (*model.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, email, role, dan, shogo, dojo_id, joined_at, created_at, updated_at
		 FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}
	return u, nil
}

// ListUsersByDojo returns users belonging to a specific dojo.
func (r *Repo) ListUsersByDojo(ctx context.Context, dojoID string) ([]model.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, email, role, dan, shogo, dojo_id, joined_at, created_at, updated_at
		 FROM users WHERE dojo_id = $1 ORDER BY name`, dojoID)
	if err != nil {
		return nil, fmt.Errorf("list users by dojo %s: %w", dojoID, err)
	}
	defer rows.Close()
	return collectUsers(rows)
}

func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	var dan, shogo, dojoID *string
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &dan, &shogo, &dojoID, &u.JoinedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	if dan != nil {
		d := model.DanRank(*dan)
		u.Dan = &d
	}
	if shogo != nil {
		s := model.ShogoTitle(*shogo)
		u.Shogo = &s
	}
	u.DojoID = dojoID
	return &u, nil
}

func collectUsers(rows pgx.Rows) ([]model.User, error) {
	var users []model.User
	for rows.Next() {
		var u model.User
		var dan, shogo, dojoID *string
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &dan, &shogo, &dojoID, &u.JoinedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		if dan != nil {
			d := model.DanRank(*dan)
			u.Dan = &d
		}
		if shogo != nil {
			s := model.ShogoTitle(*shogo)
			u.Shogo = &s
		}
		u.DojoID = dojoID
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

// ListDojos returns all dojos.
func (r *Repo) ListDojos(ctx context.Context) ([]model.Dojo, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, address, target_lanes, open_time, close_time, created_at, updated_at
		 FROM dojos ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list dojos: %w", err)
	}
	defer rows.Close()

	var dojos []model.Dojo
	for rows.Next() {
		var d model.Dojo
		if err := rows.Scan(&d.ID, &d.Name, &d.Address, &d.TargetLanes, &d.OpenTime, &d.CloseTime, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan dojo: %w", err)
		}
		dojos = append(dojos, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dojos: %w", err)
	}
	return dojos, nil
}

// GetDojo returns a dojo by ID.
func (r *Repo) GetDojo(ctx context.Context, id string) (*model.Dojo, error) {
	var d model.Dojo
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, address, target_lanes, open_time, close_time, created_at, updated_at
		 FROM dojos WHERE id = $1`, id).
		Scan(&d.ID, &d.Name, &d.Address, &d.TargetLanes, &d.OpenTime, &d.CloseTime, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get dojo %s: %w", id, err)
	}
	return &d, nil
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

// ListPractices returns practices, optionally filtered by userID.
func (r *Repo) ListPractices(ctx context.Context, userID *string) ([]model.Practice, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if userID != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, dojo_id, date, hit_rate, arrow_count, notes, instructor_comment, created_at, updated_at
			 FROM practices WHERE user_id = $1 ORDER BY date DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, dojo_id, date, hit_rate, arrow_count, notes, instructor_comment, created_at, updated_at
			 FROM practices ORDER BY date DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list practices: %w", err)
	}
	defer rows.Close()

	var practices []model.Practice
	for rows.Next() {
		var p model.Practice
		if err := rows.Scan(&p.ID, &p.UserID, &p.DojoID, &p.Date, &p.HitRate, &p.ArrowCount, &p.Notes, &p.InstructorComment, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan practice: %w", err)
		}
		practices = append(practices, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate practices: %w", err)
	}
	return practices, nil
}

// GetPractice returns a practice by ID.
func (r *Repo) GetPractice(ctx context.Context, id string) (*model.Practice, error) {
	var p model.Practice
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, dojo_id, date, hit_rate, arrow_count, notes, instructor_comment, created_at, updated_at
		 FROM practices WHERE id = $1`, id).
		Scan(&p.ID, &p.UserID, &p.DojoID, &p.Date, &p.HitRate, &p.ArrowCount, &p.Notes, &p.InstructorComment, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get practice %s: %w", id, err)
	}
	return &p, nil
}

// CreatePractice inserts a new practice record.
func (r *Repo) CreatePractice(ctx context.Context, p *model.Practice) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO practices (id, user_id, dojo_id, date, hit_rate, arrow_count, notes, instructor_comment, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		p.ID, p.UserID, p.DojoID, p.Date, p.HitRate, p.ArrowCount, p.Notes, p.InstructorComment, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create practice: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

// ListVideos returns videos, optionally filtered by userID.
func (r *Repo) ListVideos(ctx context.Context, userID *string) ([]model.Video, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if userID != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, practice_id, file_name, file_size, duration, mime_type, status, url, created_at, updated_at
			 FROM videos WHERE user_id = $1 ORDER BY created_at DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, practice_id, file_name, file_size, duration, mime_type, status, url, created_at, updated_at
			 FROM videos ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list videos: %w", err)
	}
	defer rows.Close()

	var videos []model.Video
	for rows.Next() {
		var v model.Video
		if err := rows.Scan(&v.ID, &v.UserID, &v.PracticeID, &v.FileName, &v.FileSize, &v.Duration, &v.MimeType, &v.Status, &v.URL, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan video: %w", err)
		}
		videos = append(videos, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate videos: %w", err)
	}
	return videos, nil
}

// GetVideo returns a video by ID.
func (r *Repo) GetVideo(ctx context.Context, id string) (*model.Video, error) {
	var v model.Video
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, practice_id, file_name, file_size, duration, mime_type, status, url, created_at, updated_at
		 FROM videos WHERE id = $1`, id).
		Scan(&v.ID, &v.UserID, &v.PracticeID, &v.FileName, &v.FileSize, &v.Duration, &v.MimeType, &v.Status, &v.URL, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get video %s: %w", id, err)
	}
	return &v, nil
}

// CreateVideo inserts a new video record.
func (r *Repo) CreateVideo(ctx context.Context, v *model.Video) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO videos (id, user_id, practice_id, file_name, file_size, duration, mime_type, status, url, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		v.ID, v.UserID, v.PracticeID, v.FileName, v.FileSize, v.Duration, v.MimeType, v.Status, v.URL, v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create video: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

// ListAnalyses returns analyses, optionally filtered by userID.
func (r *Repo) ListAnalyses(ctx context.Context, userID *string) ([]model.Analysis, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if userID != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT id, video_id, user_id, scores, phases, overall_score, feedback, created_at
			 FROM analyses WHERE user_id = $1 ORDER BY created_at DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, video_id, user_id, scores, phases, overall_score, feedback, created_at
			 FROM analyses ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list analyses: %w", err)
	}
	defer rows.Close()

	var analyses []model.Analysis
	for rows.Next() {
		var a model.Analysis
		var scoresJSON, phasesJSON []byte
		if err := rows.Scan(&a.ID, &a.VideoID, &a.UserID, &scoresJSON, &phasesJSON, &a.OverallScore, &a.Feedback, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan analysis: %w", err)
		}
		if err := json.Unmarshal(scoresJSON, &a.Scores); err != nil {
			return nil, fmt.Errorf("unmarshal scores: %w", err)
		}
		if err := json.Unmarshal(phasesJSON, &a.Phases); err != nil {
			return nil, fmt.Errorf("unmarshal phases: %w", err)
		}
		analyses = append(analyses, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate analyses: %w", err)
	}
	return analyses, nil
}

// GetAnalysis returns an analysis by ID.
func (r *Repo) GetAnalysis(ctx context.Context, id string) (*model.Analysis, error) {
	var a model.Analysis
	var scoresJSON, phasesJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, video_id, user_id, scores, phases, overall_score, feedback, created_at
		 FROM analyses WHERE id = $1`, id).
		Scan(&a.ID, &a.VideoID, &a.UserID, &scoresJSON, &phasesJSON, &a.OverallScore, &a.Feedback, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get analysis %s: %w", id, err)
	}
	if err := json.Unmarshal(scoresJSON, &a.Scores); err != nil {
		return nil, fmt.Errorf("unmarshal scores: %w", err)
	}
	if err := json.Unmarshal(phasesJSON, &a.Phases); err != nil {
		return nil, fmt.Errorf("unmarshal phases: %w", err)
	}
	return &a, nil
}

// GetAnalysisByVideo returns an analysis for a given video.
func (r *Repo) GetAnalysisByVideo(ctx context.Context, videoID string) (*model.Analysis, error) {
	var a model.Analysis
	var scoresJSON, phasesJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, video_id, user_id, scores, phases, overall_score, feedback, created_at
		 FROM analyses WHERE video_id = $1`, videoID).
		Scan(&a.ID, &a.VideoID, &a.UserID, &scoresJSON, &phasesJSON, &a.OverallScore, &a.Feedback, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get analysis by video %s: %w", videoID, err)
	}
	if err := json.Unmarshal(scoresJSON, &a.Scores); err != nil {
		return nil, fmt.Errorf("unmarshal scores: %w", err)
	}
	if err := json.Unmarshal(phasesJSON, &a.Phases); err != nil {
		return nil, fmt.Errorf("unmarshal phases: %w", err)
	}
	return &a, nil
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

// ListReservations returns reservations with optional dojo and date filters.
func (r *Repo) ListReservations(ctx context.Context, dojoID *string, date *string) ([]model.Reservation, error) {
	query := `SELECT id, dojo_id, user_id, lane_number, date, start_time, end_time, created_at, updated_at FROM reservations WHERE 1=1`
	args := []any{}
	argIdx := 1

	if dojoID != nil {
		query += fmt.Sprintf(" AND dojo_id = $%d", argIdx)
		args = append(args, *dojoID)
		argIdx++
	}
	if date != nil {
		query += fmt.Sprintf(" AND date = $%d", argIdx)
		args = append(args, *date)
		argIdx++
	}
	_ = argIdx // suppress unused warning
	query += " ORDER BY date, start_time"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		if err := rows.Scan(&res.ID, &res.DojoID, &res.UserID, &res.LaneNumber, &res.Date, &res.StartTime, &res.EndTime, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan reservation: %w", err)
		}
		reservations = append(reservations, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reservations: %w", err)
	}
	return reservations, nil
}

// GetReservation returns a reservation by ID.
func (r *Repo) GetReservation(ctx context.Context, id string) (*model.Reservation, error) {
	var res model.Reservation
	err := r.pool.QueryRow(ctx,
		`SELECT id, dojo_id, user_id, lane_number, date, start_time, end_time, created_at, updated_at
		 FROM reservations WHERE id = $1`, id).
		Scan(&res.ID, &res.DojoID, &res.UserID, &res.LaneNumber, &res.Date, &res.StartTime, &res.EndTime, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get reservation %s: %w", id, err)
	}
	return &res, nil
}

// CreateReservation inserts a new reservation after checking for conflicts.
//
// The overlap check (SELECT) and the INSERT run inside a single Serializable
// transaction to close the TOCTOU window that would otherwise allow two
// concurrent requests to both pass the check and double-book the same lane.
// Under Serializable isolation PostgreSQL detects the read/write conflict and
// aborts one transaction with a serialization failure (SQLSTATE 40001); that,
// along with a unique-constraint violation (23505) on an identical start time,
// is surfaced to the caller as ErrReservationConflict.
func (r *Repo) CreateReservation(ctx context.Context, res *model.Reservation) (err error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("begin reservation tx: %w", err)
	}
	defer func() {
		if err != nil {
			// Best-effort rollback; the commit path already released the tx.
			_ = tx.Rollback(ctx)
		}
	}()

	// Check for overlap: start_a < end_b AND start_b < end_a
	var conflictCount int
	if err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM reservations
		 WHERE dojo_id = $1 AND lane_number = $2 AND date = $3
		 AND start_time < $4 AND $5 < end_time`,
		res.DojoID, res.LaneNumber, res.Date, res.EndTime, res.StartTime).Scan(&conflictCount); err != nil {
		return fmt.Errorf("check reservation conflict: %w", err)
	}
	if conflictCount > 0 {
		err = ErrReservationConflict
		return err
	}

	if _, err = tx.Exec(ctx,
		`INSERT INTO reservations (id, dojo_id, user_id, lane_number, date, start_time, end_time, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		res.ID, res.DojoID, res.UserID, res.LaneNumber, res.Date, res.StartTime, res.EndTime, res.CreatedAt, res.UpdatedAt); err != nil {
		if isConflictError(err) {
			err = ErrReservationConflict
			return err
		}
		return fmt.Errorf("create reservation: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		if isConflictError(err) {
			err = ErrReservationConflict
			return err
		}
		return fmt.Errorf("commit reservation: %w", err)
	}
	return nil
}

// isConflictError reports whether err is a PostgreSQL error that, for reservation
// creation, means a competing booking won the race: a unique-constraint violation
// (23505) or a serialization failure (40001).
func isConflictError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" || pgErr.Code == "40001"
	}
	return false
}

// DeleteReservation removes a reservation by ID.
func (r *Repo) DeleteReservation(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM reservations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete reservation %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

// ListExamChecklists returns exam checklists, optionally filtered by userID.
func (r *Repo) ListExamChecklists(ctx context.Context, userID *string) ([]model.ExamChecklist, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if userID != nil {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, target_dan, items, progress_rate, created_at, updated_at
			 FROM exam_checklists WHERE user_id = $1 ORDER BY created_at DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, user_id, target_dan, items, progress_rate, created_at, updated_at
			 FROM exam_checklists ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list exam checklists: %w", err)
	}
	defer rows.Close()

	var checklists []model.ExamChecklist
	for rows.Next() {
		var c model.ExamChecklist
		var itemsJSON []byte
		if err := rows.Scan(&c.ID, &c.UserID, &c.TargetDan, &itemsJSON, &c.ProgressRate, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan exam checklist: %w", err)
		}
		if err := json.Unmarshal(itemsJSON, &c.Items); err != nil {
			return nil, fmt.Errorf("unmarshal checklist items: %w", err)
		}
		checklists = append(checklists, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate exam checklists: %w", err)
	}
	return checklists, nil
}

// GetExamChecklist returns an exam checklist by ID.
func (r *Repo) GetExamChecklist(ctx context.Context, id string) (*model.ExamChecklist, error) {
	var c model.ExamChecklist
	var itemsJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, target_dan, items, progress_rate, created_at, updated_at
		 FROM exam_checklists WHERE id = $1`, id).
		Scan(&c.ID, &c.UserID, &c.TargetDan, &itemsJSON, &c.ProgressRate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get exam checklist %s: %w", id, err)
	}
	if err := json.Unmarshal(itemsJSON, &c.Items); err != nil {
		return nil, fmt.Errorf("unmarshal checklist items: %w", err)
	}
	return &c, nil
}

// UpdateExamChecklist updates an existing exam checklist (items + progressRate).
func (r *Repo) UpdateExamChecklist(ctx context.Context, c *model.ExamChecklist) error {
	itemsJSON, err := json.Marshal(c.Items)
	if err != nil {
		return fmt.Errorf("marshal checklist items: %w", err)
	}
	_, err = r.pool.Exec(ctx,
		`UPDATE exam_checklists SET items = $1, progress_rate = $2, updated_at = $3 WHERE id = $4`,
		itemsJSON, c.ProgressRate, c.UpdatedAt, c.ID)
	if err != nil {
		return fmt.Errorf("update exam checklist: %w", err)
	}
	return nil
}

// CountUsersByDojo returns the count of users belonging to a dojo.
func (r *Repo) CountUsersByDojo(ctx context.Context, dojoID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE dojo_id = $1`, dojoID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users by dojo: %w", err)
	}
	return count, nil
}

// ListReservationsByDojoAndDate returns reservations for a specific dojo and date.
func (r *Repo) ListReservationsByDojoAndDate(ctx context.Context, dojoID string, date string) ([]model.Reservation, error) {
	d := &dojoID
	dt := &date
	return r.ListReservations(ctx, d, dt)
}

// Ping checks database connectivity.
func (r *Repo) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// Now returns the current time (UTC).
func Now() time.Time {
	return time.Now().UTC()
}
