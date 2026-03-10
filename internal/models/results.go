package models

import "time"

type LessonResults struct {
	Id             int
	User           string
	When           time.Time
	Difficulty     int
	Wpm, ErrorRate float64
	TimeTaken      time.Duration
}

type Storage interface {
	SaveTypingLesson(result LessonResults) error
	GetTypingLessons(user string) ([]LessonResults, error)
}
