// Package model defines the core domain types for kyudo-dojo-hub.
package model

import "time"

// HassetsuPhase represents one of the 8 phases in Japanese archery (射法八節).
type HassetsuPhase string

const (
	Ashibumi   HassetsuPhase = "ashibumi"
	Dozukuri   HassetsuPhase = "dozukuri"
	Yugamae    HassetsuPhase = "yugamae"
	Uchiokoshi HassetsuPhase = "uchiokoshi"
	Hikiwake   HassetsuPhase = "hikiwake"
	Kai        HassetsuPhase = "kai"
	Hanare     HassetsuPhase = "hanare"
	Zanshin    HassetsuPhase = "zanshin"
)

// DanRank represents a kyu/dan rank in kyudo.
type DanRank string

const (
	Shodan  DanRank = "shodan"
	Nidan   DanRank = "nidan"
	Sandan  DanRank = "sandan"
	Yondan  DanRank = "yondan"
	Godan   DanRank = "godan"
	Rokudan DanRank = "rokudan"
	Nanadan DanRank = "nanadan"
	Hachidan DanRank = "hachidan"
	Kudan   DanRank = "kudan"
	Judan   DanRank = "judan"
)

// ShogoTitle represents a shogo (称号) title.
type ShogoTitle string

const (
	Renshi ShogoTitle = "renshi"
	Kyoshi ShogoTitle = "kyoshi"
	Hanshi ShogoTitle = "hanshi"
)

// UserRole represents the role of a user.
type UserRole string

const (
	RolePractitioner UserRole = "practitioner"
	RoleManager      UserRole = "manager"
	RoleAdmin        UserRole = "admin"
)

// User represents a practitioner, manager, or admin.
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

// Dojo represents a kyudo training hall.
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

// Practice represents a practice session record.
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

// VideoStatus represents the processing status of a video.
type VideoStatus string

const (
	VideoUploading  VideoStatus = "uploading"
	VideoProcessing VideoStatus = "processing"
	VideoCompleted  VideoStatus = "completed"
	VideoFailed     VideoStatus = "failed"
)

// Video represents an uploaded archery form video.
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

// HassetsuScores maps each phase to a 0-100 score.
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

// PhaseSegment represents a time segment for a hassetsu phase in a video.
type PhaseSegment struct {
	Phase     HassetsuPhase `json:"phase"`
	StartTime float64       `json:"startTime"`
	EndTime   float64       `json:"endTime"`
}

// Analysis represents the result of a shooting form analysis.
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

// ExamChecklistItem is a single item on a dan exam checklist.
type ExamChecklistItem struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// ExamChecklist is a checklist for a target dan rank exam.
type ExamChecklist struct {
	ID           string              `json:"id"`
	UserID       string              `json:"userId"`
	TargetDan    DanRank             `json:"targetDan"`
	Items        []ExamChecklistItem `json:"items"`
	ProgressRate int                 `json:"progressRate"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

// DashboardSummary provides an overview for a dojo's dashboard.
type DashboardSummary struct {
	TodayReservationCount int           `json:"todayReservationCount"`
	TotalMemberCount      int           `json:"totalMemberCount"`
	TodayReservations     []Reservation `json:"todayReservations"`
}
