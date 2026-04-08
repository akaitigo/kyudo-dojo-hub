package store

import "github.com/ryusei/kyudo-dojo-hub/backend/internal/model"

// seed populates the store with initial sample data matching the frontend mock data.
func (s *Store) seed() {
	s.users = []model.User{
		{ID: "user-001", Name: "田中太郎", Email: "tanaka@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanSandan), DojoID: strPtr("dojo-001"), JoinedAt: "2023-04-01", CreatedAt: mustTime("2023-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-01-15T10:00:00Z")},
		{ID: "user-002", Name: "佐藤花子", Email: "sato@example.com", Role: model.RoleManager, Dan: danPtr(model.DanGodan), Shogo: shogoPtr(model.ShogoRenshi), DojoID: strPtr("dojo-001"), JoinedAt: "2018-04-01", CreatedAt: mustTime("2018-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-02-20T10:00:00Z")},
		{ID: "user-003", Name: "鈴木一郎", Email: "suzuki@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanShodan), DojoID: strPtr("dojo-001"), JoinedAt: "2025-04-01", CreatedAt: mustTime("2025-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-10T10:00:00Z")},
		{ID: "user-004", Name: "高橋美咲", Email: "takahashi@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanNidan), DojoID: strPtr("dojo-001"), JoinedAt: "2024-04-01", CreatedAt: mustTime("2024-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-01T10:00:00Z")},
		{ID: "user-005", Name: "山本健太", Email: "yamamoto@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanYondan), DojoID: strPtr("dojo-002"), JoinedAt: "2020-04-01", CreatedAt: mustTime("2020-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-02-15T10:00:00Z")},
		{ID: "user-006", Name: "中村雅子", Email: "nakamura@example.com", Role: model.RoleManager, Dan: danPtr(model.DanRokudan), Shogo: shogoPtr(model.ShogoKyoshi), DojoID: strPtr("dojo-002"), JoinedAt: "2015-04-01", CreatedAt: mustTime("2015-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-01-10T10:00:00Z")},
		{ID: "user-007", Name: "小林大輔", Email: "kobayashi@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanSandan), DojoID: strPtr("dojo-001"), JoinedAt: "2022-04-01", CreatedAt: mustTime("2022-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-05T10:00:00Z")},
		{ID: "user-008", Name: "加藤由美", Email: "kato@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanNidan), DojoID: strPtr("dojo-002"), JoinedAt: "2024-10-01", CreatedAt: mustTime("2024-10-01T00:00:00Z"), UpdatedAt: mustTime("2026-02-28T10:00:00Z")},
		{ID: "user-009", Name: "伊藤誠", Email: "ito@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanShodan), DojoID: strPtr("dojo-001"), JoinedAt: "2025-10-01", CreatedAt: mustTime("2025-10-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-15T10:00:00Z")},
		{ID: "user-010", Name: "渡辺真理", Email: "watanabe@example.com", Role: model.RolePractitioner, Dan: danPtr(model.DanYondan), DojoID: strPtr("dojo-001"), JoinedAt: "2019-04-01", CreatedAt: mustTime("2019-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-01-20T10:00:00Z")},
		{ID: "user-011", Name: "松本幸子", Email: "matsumoto@example.com", Role: model.RoleAdmin, Dan: danPtr(model.DanNanadan), Shogo: shogoPtr(model.ShogoHanshi), DojoID: strPtr("dojo-001"), JoinedAt: "2010-04-01", CreatedAt: mustTime("2010-04-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-20T10:00:00Z")},
		{ID: "user-012", Name: "井上翔", Email: "inoue@example.com", Role: model.RolePractitioner, DojoID: strPtr("dojo-002"), JoinedAt: "2026-01-01", CreatedAt: mustTime("2026-01-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-25T10:00:00Z")},
	}

	s.dojos = []model.Dojo{
		{ID: "dojo-001", Name: "東京弓道場", Address: "東京都千代田区丸の内1-1-1", TargetLanes: 6, OpenTime: "09:00", CloseTime: "21:00", CreatedAt: mustTime("2020-01-01T00:00:00Z"), UpdatedAt: mustTime("2026-01-01T00:00:00Z")},
		{ID: "dojo-002", Name: "大阪弓友会弓道場", Address: "大阪府大阪市中央区大手前2-2-2", TargetLanes: 4, OpenTime: "10:00", CloseTime: "20:00", CreatedAt: mustTime("2021-06-01T00:00:00Z"), UpdatedAt: mustTime("2026-02-01T00:00:00Z")},
	}

	s.practices = []model.Practice{
		{ID: "practice-001", UserID: "user-001", DojoID: strPtr("dojo-001"), Date: "2026-03-28", HitRate: 65, ArrowCount: 40, Notes: "会での安定感が向上。離れが少し硬い。", InstructorComment: "会の伸びが良くなっている。離れは肩を意識すること。", CreatedAt: mustTime("2026-03-28T18:00:00Z"), UpdatedAt: mustTime("2026-03-28T18:30:00Z")},
		{ID: "practice-002", UserID: "user-001", DojoID: strPtr("dojo-001"), Date: "2026-03-25", HitRate: 58, ArrowCount: 36, Notes: "胴造りが崩れやすかった。足踏みから意識する。", InstructorComment: "", CreatedAt: mustTime("2026-03-25T18:00:00Z"), UpdatedAt: mustTime("2026-03-25T18:00:00Z")},
		{ID: "practice-003", UserID: "user-001", DojoID: strPtr("dojo-001"), Date: "2026-03-22", HitRate: 70, ArrowCount: 48, Notes: "好調。引分けのバランスが良かった。", InstructorComment: "この調子を維持すること。", CreatedAt: mustTime("2026-03-22T18:00:00Z"), UpdatedAt: mustTime("2026-03-22T18:30:00Z")},
		{ID: "practice-004", UserID: "user-003", DojoID: strPtr("dojo-001"), Date: "2026-03-28", HitRate: 35, ArrowCount: 24, Notes: "弓構えの手の内が安定しない。", InstructorComment: "手の内の形を毎回確認してから引くこと。", CreatedAt: mustTime("2026-03-28T19:00:00Z"), UpdatedAt: mustTime("2026-03-28T19:30:00Z")},
		{ID: "practice-005", UserID: "user-004", DojoID: strPtr("dojo-001"), Date: "2026-03-27", HitRate: 50, ArrowCount: 32, Notes: "打起しの高さを変えてみた。少し改善。", InstructorComment: "", CreatedAt: mustTime("2026-03-27T18:00:00Z"), UpdatedAt: mustTime("2026-03-27T18:00:00Z")},
		{ID: "practice-006", UserID: "user-005", DojoID: strPtr("dojo-002"), Date: "2026-03-28", HitRate: 75, ArrowCount: 60, Notes: "大会前の調整。残心まで意識できている。", InstructorComment: "大会に向けて良い状態。自信を持って臨むこと。", CreatedAt: mustTime("2026-03-28T17:00:00Z"), UpdatedAt: mustTime("2026-03-28T17:30:00Z")},
		{ID: "practice-007", UserID: "user-007", DojoID: strPtr("dojo-001"), Date: "2026-03-26", HitRate: 62, ArrowCount: 44, Notes: "会の時間を長くする練習。3秒以上を目標。", InstructorComment: "", CreatedAt: mustTime("2026-03-26T18:00:00Z"), UpdatedAt: mustTime("2026-03-26T18:00:00Z")},
		{ID: "practice-008", UserID: "user-008", DojoID: strPtr("dojo-002"), Date: "2026-03-27", HitRate: 42, ArrowCount: 28, Notes: "肩の力を抜くことを意識。まだ力みがち。", InstructorComment: "呼吸を大切に。力を抜いて引くこと。", CreatedAt: mustTime("2026-03-27T19:00:00Z"), UpdatedAt: mustTime("2026-03-27T19:30:00Z")},
		{ID: "practice-009", UserID: "user-009", DojoID: strPtr("dojo-001"), Date: "2026-03-29", HitRate: 30, ArrowCount: 20, Notes: "初めての巻藁練習。基本の確認。", InstructorComment: "足踏みの幅と角度を固定すること。", CreatedAt: mustTime("2026-03-29T10:00:00Z"), UpdatedAt: mustTime("2026-03-29T10:30:00Z")},
		{ID: "practice-010", UserID: "user-010", DojoID: strPtr("dojo-001"), Date: "2026-03-28", HitRate: 72, ArrowCount: 52, Notes: "審査に向けて体配の練習も実施。", InstructorComment: "体配は問題ない。射の安定感をさらに高めること。", CreatedAt: mustTime("2026-03-28T16:00:00Z"), UpdatedAt: mustTime("2026-03-28T16:30:00Z")},
		{ID: "practice-011", UserID: "user-001", DojoID: strPtr("dojo-001"), Date: "2026-03-20", HitRate: 55, ArrowCount: 32, Notes: "雨の日。室内での素引き練習中心。", InstructorComment: "", CreatedAt: mustTime("2026-03-20T18:00:00Z"), UpdatedAt: mustTime("2026-03-20T18:00:00Z")},
		{ID: "practice-012", UserID: "user-004", DojoID: strPtr("dojo-001"), Date: "2026-03-29", HitRate: 55, ArrowCount: 36, Notes: "前回より改善。引分けでの肘の位置を修正。", InstructorComment: "良い方向に向かっている。続けること。", CreatedAt: mustTime("2026-03-29T18:00:00Z"), UpdatedAt: mustTime("2026-03-29T18:30:00Z")},
	}

	s.videos = []model.Video{
		{ID: "video-001", UserID: "user-001", PracticeID: strPtr("practice-001"), FileName: "tanaka_20260328.mp4", FileSize: 52_428_800, Duration: 45, MimeType: "video/mp4", Status: model.StatusCompleted, URL: "", CreatedAt: mustTime("2026-03-28T18:10:00Z"), UpdatedAt: mustTime("2026-03-28T18:11:00Z")},
		{ID: "video-002", UserID: "user-005", PracticeID: strPtr("practice-006"), FileName: "yamamoto_20260328.mp4", FileSize: 78_643_200, Duration: 38, MimeType: "video/mp4", Status: model.StatusCompleted, URL: "", CreatedAt: mustTime("2026-03-28T17:10:00Z"), UpdatedAt: mustTime("2026-03-28T17:11:00Z")},
		{ID: "video-003", UserID: "user-001", PracticeID: strPtr("practice-003"), FileName: "tanaka_20260322.mov", FileSize: 104_857_600, Duration: 52, MimeType: "video/quicktime", Status: model.StatusCompleted, URL: "", CreatedAt: mustTime("2026-03-22T18:10:00Z"), UpdatedAt: mustTime("2026-03-22T18:11:00Z")},
	}

	s.analyses = []model.Analysis{
		{
			ID: "analysis-001", VideoID: "video-001", UserID: "user-001",
			Scores: model.HassetsuScores{Ashibumi: 78, Dozukuri: 72, Yugamae: 80, Uchiokoshi: 75, Hikiwake: 68, Kai: 82, Hanare: 60, Zanshin: 70},
			Phases: []model.PhaseSegment{
				{Phase: model.PhaseAshibumi, StartTime: 0, EndTime: 3.5},
				{Phase: model.PhaseDozukuri, StartTime: 3.5, EndTime: 7.0},
				{Phase: model.PhaseYugamae, StartTime: 7.0, EndTime: 12.0},
				{Phase: model.PhaseUchiokoshi, StartTime: 12.0, EndTime: 16.0},
				{Phase: model.PhaseHikiwake, StartTime: 16.0, EndTime: 24.0},
				{Phase: model.PhaseKai, StartTime: 24.0, EndTime: 30.0},
				{Phase: model.PhaseHanare, StartTime: 30.0, EndTime: 31.0},
				{Phase: model.PhaseZanshin, StartTime: 31.0, EndTime: 35.0},
			},
			OverallScore: 73, Feedback: "会での安定感は良好。離れの瞬間に右肩が上がる傾向あり。引分けでの左右均等な力配分を意識すること。",
			CreatedAt: mustTime("2026-03-28T18:15:00Z"),
		},
		{
			ID: "analysis-002", VideoID: "video-002", UserID: "user-005",
			Scores: model.HassetsuScores{Ashibumi: 85, Dozukuri: 88, Yugamae: 82, Uchiokoshi: 80, Hikiwake: 85, Kai: 90, Hanare: 78, Zanshin: 82},
			Phases: []model.PhaseSegment{
				{Phase: model.PhaseAshibumi, StartTime: 0, EndTime: 2.5},
				{Phase: model.PhaseDozukuri, StartTime: 2.5, EndTime: 5.5},
				{Phase: model.PhaseYugamae, StartTime: 5.5, EndTime: 10.0},
				{Phase: model.PhaseUchiokoshi, StartTime: 10.0, EndTime: 13.0},
				{Phase: model.PhaseHikiwake, StartTime: 13.0, EndTime: 20.0},
				{Phase: model.PhaseKai, StartTime: 20.0, EndTime: 26.0},
				{Phase: model.PhaseHanare, StartTime: 26.0, EndTime: 27.0},
				{Phase: model.PhaseZanshin, StartTime: 27.0, EndTime: 30.0},
			},
			OverallScore: 84, Feedback: "全体的に安定した射。会の伸びが特に良い。離れでの残身の意識をさらに高めるとより良くなる。",
			CreatedAt: mustTime("2026-03-28T17:15:00Z"),
		},
	}

	s.reservations = []model.Reservation{
		{ID: "res-001", DojoID: "dojo-001", UserID: "user-001", LaneNumber: 1, Date: "2026-03-30", StartTime: "09:00", EndTime: "10:00", CreatedAt: mustTime("2026-03-28T10:00:00Z"), UpdatedAt: mustTime("2026-03-28T10:00:00Z")},
		{ID: "res-002", DojoID: "dojo-001", UserID: "user-003", LaneNumber: 2, Date: "2026-03-30", StartTime: "09:00", EndTime: "10:00", CreatedAt: mustTime("2026-03-28T11:00:00Z"), UpdatedAt: mustTime("2026-03-28T11:00:00Z")},
		{ID: "res-003", DojoID: "dojo-001", UserID: "user-004", LaneNumber: 3, Date: "2026-03-30", StartTime: "10:00", EndTime: "11:00", CreatedAt: mustTime("2026-03-28T12:00:00Z"), UpdatedAt: mustTime("2026-03-28T12:00:00Z")},
		{ID: "res-004", DojoID: "dojo-001", UserID: "user-007", LaneNumber: 1, Date: "2026-03-30", StartTime: "14:00", EndTime: "15:00", CreatedAt: mustTime("2026-03-29T09:00:00Z"), UpdatedAt: mustTime("2026-03-29T09:00:00Z")},
		{ID: "res-005", DojoID: "dojo-001", UserID: "user-010", LaneNumber: 4, Date: "2026-03-30", StartTime: "15:00", EndTime: "16:00", CreatedAt: mustTime("2026-03-29T10:00:00Z"), UpdatedAt: mustTime("2026-03-29T10:00:00Z")},
		{ID: "res-006", DojoID: "dojo-002", UserID: "user-005", LaneNumber: 1, Date: "2026-03-30", StartTime: "10:00", EndTime: "11:00", CreatedAt: mustTime("2026-03-28T15:00:00Z"), UpdatedAt: mustTime("2026-03-28T15:00:00Z")},
		{ID: "res-007", DojoID: "dojo-002", UserID: "user-008", LaneNumber: 2, Date: "2026-03-30", StartTime: "10:00", EndTime: "11:00", CreatedAt: mustTime("2026-03-28T16:00:00Z"), UpdatedAt: mustTime("2026-03-28T16:00:00Z")},
		{ID: "res-008", DojoID: "dojo-001", UserID: "user-001", LaneNumber: 2, Date: "2026-03-31", StartTime: "18:00", EndTime: "19:00", CreatedAt: mustTime("2026-03-29T20:00:00Z"), UpdatedAt: mustTime("2026-03-29T20:00:00Z")},
		{ID: "res-009", DojoID: "dojo-001", UserID: "user-009", LaneNumber: 5, Date: "2026-03-31", StartTime: "10:00", EndTime: "11:00", CreatedAt: mustTime("2026-03-30T08:00:00Z"), UpdatedAt: mustTime("2026-03-30T08:00:00Z")},
		{ID: "res-010", DojoID: "dojo-001", UserID: "user-004", LaneNumber: 6, Date: "2026-04-01", StartTime: "19:00", EndTime: "20:00", CreatedAt: mustTime("2026-03-30T09:00:00Z"), UpdatedAt: mustTime("2026-03-30T09:00:00Z")},
	}

	s.examChecklists = []model.ExamChecklist{
		{
			ID: "exam-001", UserID: "user-001", TargetDan: model.DanYondan,
			Items: []model.ExamChecklistItem{
				{ID: "item-001", Category: "体配", Description: "入場から退場までの一連の動作が正確に行える", Completed: true},
				{ID: "item-002", Category: "体配", Description: "坐射の所作が正しく行える", Completed: true},
				{ID: "item-003", Category: "射技", Description: "射法八節を正しい手順で行える", Completed: true},
				{ID: "item-004", Category: "射技", Description: "会の伸びが3秒以上ある", Completed: false},
				{ID: "item-005", Category: "射技", Description: "残心が自然に取れている", Completed: false},
				{ID: "item-006", Category: "学科", Description: "射法八節の意義を説明できる", Completed: true},
				{ID: "item-007", Category: "学科", Description: "弓道の歴史について理解している", Completed: false},
				{ID: "item-008", Category: "的中", Description: "審査本番で4本中2本以上的中できる", Completed: false},
			},
			ProgressRate: 50, CreatedAt: mustTime("2026-02-01T00:00:00Z"), UpdatedAt: mustTime("2026-03-28T18:00:00Z"),
		},
		{
			ID: "exam-002", UserID: "user-003", TargetDan: model.DanNidan,
			Items: []model.ExamChecklistItem{
				{ID: "item-009", Category: "体配", Description: "立射の基本動作が正しく行える", Completed: true},
				{ID: "item-010", Category: "射技", Description: "射法八節の基本的な動作ができる", Completed: true},
				{ID: "item-011", Category: "射技", Description: "弓構えで手の内が正しく作れる", Completed: false},
				{ID: "item-012", Category: "学科", Description: "射法八節の名称を全て言える", Completed: true},
				{ID: "item-013", Category: "的中", Description: "審査本番で4本中1本以上的中できる", Completed: true},
			},
			ProgressRate: 80, CreatedAt: mustTime("2026-01-15T00:00:00Z"), UpdatedAt: mustTime("2026-03-29T10:00:00Z"),
		},
	}
}
