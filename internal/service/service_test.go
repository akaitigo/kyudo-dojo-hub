package service

import (
	"regexp"
	"testing"
)

func TestValidation_CreatePractice(t *testing.T) {
	tests := []struct {
		name    string
		input   CreatePracticeInput
		wantErr string
	}{
		{
			name:    "hitRate negative",
			input:   CreatePracticeInput{HitRate: -1, ArrowCount: 10, Notes: "", InstructorComment: ""},
			wantErr: "的中率は 0〜100 の範囲で入力してください",
		},
		{
			name:    "hitRate over 100",
			input:   CreatePracticeInput{HitRate: 101, ArrowCount: 10, Notes: "", InstructorComment: ""},
			wantErr: "的中率は 0〜100 の範囲で入力してください",
		},
		{
			name:    "arrowCount zero",
			input:   CreatePracticeInput{HitRate: 50, ArrowCount: 0, Notes: "", InstructorComment: ""},
			wantErr: "矢数は 1〜1000 の範囲で入力してください",
		},
		{
			name:    "arrowCount over 1000",
			input:   CreatePracticeInput{HitRate: 50, ArrowCount: 1001, Notes: "", InstructorComment: ""},
			wantErr: "矢数は 1〜1000 の範囲で入力してください",
		},
		{
			name:    "notes too long",
			input:   CreatePracticeInput{HitRate: 50, ArrowCount: 10, Notes: string(make([]byte, 5001)), InstructorComment: ""},
			wantErr: "気づきは 5,000 文字以内で入力してください",
		},
		{
			name:    "instructorComment too long",
			input:   CreatePracticeInput{HitRate: 50, ArrowCount: 10, Notes: "", InstructorComment: string(make([]byte, 5001))},
			wantErr: "師範コメントは 5,000 文字以内で入力してください",
		},
	}

	svc := &Service{repo: nil} // repo not needed for validation-only tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CreatePractice will fail at repo level for valid inputs (nil repo),
			// but validation errors should be returned before that.
			_, err := svc.CreatePractice(t.Context(), tt.input)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			valErr, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if valErr.Message != tt.wantErr {
				t.Errorf("got message %q, want %q", valErr.Message, tt.wantErr)
			}
		})
	}
}

func TestValidation_CreateVideo(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateVideoInput
		wantErr string
	}{
		{
			name:    "file too large",
			input:   CreateVideoInput{FileSize: 500*1024*1024 + 1, Duration: 60},
			wantErr: "ファイルサイズは 500MB 以下にしてください",
		},
		{
			name:    "duration too long",
			input:   CreateVideoInput{FileSize: 1024, Duration: 301},
			wantErr: "動画長は 5 分以下にしてください",
		},
	}

	svc := &Service{repo: nil}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateVideo(t.Context(), tt.input)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			valErr, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if valErr.Message != tt.wantErr {
				t.Errorf("got message %q, want %q", valErr.Message, tt.wantErr)
			}
		})
	}
}

func TestValidation_CreateReservation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateReservationInput
		wantErr string
	}{
		{
			name:    "invalid start time format",
			input:   CreateReservationInput{StartTime: "9:00", EndTime: "10:00"},
			wantErr: "時間は HH:mm 形式で入力してください",
		},
		{
			name:    "invalid end time format",
			input:   CreateReservationInput{StartTime: "09:00", EndTime: "25:00"},
			wantErr: "時間は HH:mm 形式で入力してください",
		},
		{
			name:    "end before start",
			input:   CreateReservationInput{StartTime: "14:00", EndTime: "13:00"},
			wantErr: "終了時刻は開始時刻より後に設定してください",
		},
		{
			name:    "end equals start",
			input:   CreateReservationInput{StartTime: "10:00", EndTime: "10:00"},
			wantErr: "終了時刻は開始時刻より後に設定してください",
		},
	}

	svc := &Service{repo: nil}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateReservation(t.Context(), tt.input)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			valErr, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if valErr.Message != tt.wantErr {
				t.Errorf("got message %q, want %q", valErr.Message, tt.wantErr)
			}
		})
	}
}

func TestHHMMPattern(t *testing.T) {
	pattern := regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

	valid := []string{"00:00", "09:00", "12:30", "23:59"}
	for _, v := range valid {
		if !pattern.MatchString(v) {
			t.Errorf("expected %q to match HH:mm pattern", v)
		}
	}

	invalid := []string{"9:00", "24:00", "12:60", "1:00", "abc", ""}
	for _, v := range invalid {
		if pattern.MatchString(v) {
			t.Errorf("expected %q NOT to match HH:mm pattern", v)
		}
	}
}

func TestErrorTypes(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		err := &ValidationError{Message: "test error"}
		if err.Error() != "test error" {
			t.Errorf("got %q, want %q", err.Error(), "test error")
		}
	})

	t.Run("NotFoundError", func(t *testing.T) {
		err := &NotFoundError{Resource: "テスト"}
		want := "テスト が見つかりません"
		if err.Error() != want {
			t.Errorf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := &ConflictError{Message: "conflict"}
		if err.Error() != "conflict" {
			t.Errorf("got %q, want %q", err.Error(), "conflict")
		}
	})
}
