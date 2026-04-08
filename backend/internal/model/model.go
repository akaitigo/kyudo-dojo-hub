// Package model defines the core domain types for the kyudo-dojo-hub backend.
package model

import "time"

// HassetsuPhase represents each phase of Shaho Hassetsu (射法八節).
type HassetsuPhase string

const (
	PhaseAshibumi   HassetsuPhase = "ashibumi"
	PhaseDozukuri   HassetsuPhase = "dozukuri"
	PhaseYugamae    HassetsuPhase = "yugamae"
	PhaseUchiokoshi HassetsuPhase = "uchiokoshi"
	PhaseHikiwake   HassetsuPhase = "hikiwake"
	PhaseKai        HassetsuPhase = "kai"
	PhaseHanare     HassetsuPhase = "hanare"
	PhaseZanshin    HassetsuPhase = "zanshin"
)

// AllPhases lists all eight phases in order.
var AllPhases = []HassetsuPhase{
	PhaseAshibumi, PhaseDozukuri, PhaseYugamae, PhaseUchiokoshi,
	PhaseHikiwake, PhaseKai, PhaseHanare, PhaseZanshin,
}

// ValidPhase checks whether a given string is a valid HassetsuPhase.
func ValidPhase(s string) bool {
	for _, p := range AllPhases {
		if string(p) == s {
			return true
		}
	}
	return false
}

// UserRole represents the role a user holds in the system.
type UserRole string

const (
	RolePractitioner UserRole = "practitioner"
	RoleManager      UserRole = "manager"
	RoleAdmin        UserRole = "admin"
)

// DanRank represents a dan grade.
type DanRank string

const (
	DanShodan   DanRank = "shodan"
	DanNidan    DanRank = "nidan"
	DanSandan   DanRank = "sandan"
	DanYondan   DanRank = "yondan"
	DanGodan    DanRank = "godan"
	DanRokudan  DanRank = "rokudan"
	DanNanadan  DanRank = "nanadan"
	DanHachidan DanRank = "hachidan"
	DanKudan    DanRank = "kudan"
	DanJudan    DanRank = "judan"
)

// ShogoTitle represents a shogo title.
type ShogoTitle string

const (
	ShogoRenshi ShogoTitle = "renshi"
	ShogoKyoshi ShogoTitle = "kyoshi"
	ShogoHanshi ShogoTitle = "hanshi"
)

// VideoStatus represents the processing status of a video.
type VideoStatus string

const (
	StatusUploading  VideoStatus = "uploading"
	StatusProcessing VideoStatus = "processing"
	StatusCompleted  VideoStatus = "completed"
	StatusFailed     VideoStatus = "failed"
)

// User represents a registered user.
type User struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      UserRole    `json:"role"`
	Dan       *DanRank    `json:"dan,omitempty"`
	Shogo     *ShogoTitle `json:"shogo,omitempty"`
	DojoID    *string     `json:"dojoId,omitempty"`
	JoinedAt  string      `json:"joinedAt"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// Dojo represents a kyudo dojo.
type Dojo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	TargetLanes int       `json:"targetLanes"`
	OpenTime    string    `json:"openTime"`
	CloseTime   string    `json:"closeTime"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Practice represents a single practice session record.
type Practice struct {
	ID                string    `json:"id"`
	UserID            string    `json:"userId"`
	DojoID            *string   `json:"dojoId,omitempty"`
	Date              string    `json:"date"`
	HitRate           int       `json:"hitRate"`
	ArrowCount        int       `json:"arrowCount"`
	Notes             string    `json:"notes"`
	InstructorComment string    `json:"instructorComment"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// Video represents metadata about an uploaded shooting video.
type Video struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	PracticeID *string     `json:"practiceId,omitempty"`
	FileName   string      `json:"fileName"`
	FileSize   int64       `json:"fileSize"`
	Duration   float64     `json:"duration"`
	MimeType   string      `json:"mimeType"`
	Status     VideoStatus `json:"status"`
	URL        string      `json:"url"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

// HassetsuScores holds scores for each of the eight phases.
type HassetsuScores struct {
	Ashibumi   int `json:"ashibumi"`
	Dozukuri   int `json:"dozukuri"`
	Yugamae    int `json:"yugamae"`
	Uchiokoshi int `json:"uchiokoshi"`
	Hikiwake   int `json:"hikiwake"`
	Kai        int `json:"kai"`
	Hanare     int `json:"hanare"`
	Zanshin    int `json:"zanshin"`
}

// PhaseSegment represents a time segment for a single phase.
type PhaseSegment struct {
	Phase     HassetsuPhase `json:"phase"`
	StartTime float64       `json:"startTime"`
	EndTime   float64       `json:"endTime"`
}

// Analysis represents the results of a shooting form analysis.
type Analysis struct {
	ID           string         `json:"id"`
	VideoID      string         `json:"videoId"`
	UserID       string         `json:"userId"`
	Scores       HassetsuScores `json:"scores"`
	Phases       []PhaseSegment `json:"phases"`
	OverallScore int            `json:"overallScore"`
	Feedback     string         `json:"feedback"`
	CreatedAt    time.Time      `json:"createdAt"`
}

// Reservation represents a target lane reservation.
type Reservation struct {
	ID         string    `json:"id"`
	DojoID     string    `json:"dojoId"`
	UserID     string    `json:"userId"`
	LaneNumber int       `json:"laneNumber"`
	Date       string    `json:"date"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ExamChecklistItem represents a single checklist item.
type ExamChecklistItem struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// ExamChecklist represents a dan exam preparation checklist.
type ExamChecklist struct {
	ID           string              `json:"id"`
	UserID       string              `json:"userId"`
	TargetDan    DanRank             `json:"targetDan"`
	Items        []ExamChecklistItem `json:"items"`
	ProgressRate int                 `json:"progressRate"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

// DashboardSummary holds summary data for the dojo dashboard.
type DashboardSummary struct {
	TodayReservationCount int           `json:"todayReservationCount"`
	TotalMemberCount      int           `json:"totalMemberCount"`
	TodayReservations     []Reservation `json:"todayReservations"`
}

// APIResponse wraps a successful API response.
type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

// APIError wraps an error API response.
type APIError struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// CreatePracticeInput is the input payload for creating a practice record.
type CreatePracticeInput struct {
	UserID            string  `json:"userId"`
	DojoID            *string `json:"dojoId,omitempty"`
	Date              string  `json:"date"`
	HitRate           int     `json:"hitRate"`
	ArrowCount        int     `json:"arrowCount"`
	Notes             string  `json:"notes"`
	InstructorComment string  `json:"instructorComment"`
}

// CreateVideoInput is the input payload for creating a video record.
type CreateVideoInput struct {
	UserID     string  `json:"userId"`
	PracticeID *string `json:"practiceId,omitempty"`
	FileName   string  `json:"fileName"`
	FileSize   int64   `json:"fileSize"`
	Duration   float64 `json:"duration"`
	MimeType   string  `json:"mimeType"`
	URL        string  `json:"url"`
}

// CreateReservationInput is the input payload for creating a reservation.
type CreateReservationInput struct {
	DojoID     string `json:"dojoId"`
	UserID     string `json:"userId"`
	LaneNumber int    `json:"laneNumber"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

// AnalyzeVideoRequest is the input payload for requesting a video analysis.
type AnalyzeVideoRequest struct {
	VideoID string `json:"videoId"`
	UserID  string `json:"userId"`
}
