// Package models defines the database models for the Cyber Q&A application.
package models

import (
	"gorm.io/gorm"
)

// UserA represents User A's answers in the database.
type UserA struct {
	gorm.Model
	Token        string `gorm:"uniqueIndex;size:255"`
	Answers      string `gorm:"type:text"` // JSON string of answers
	ShareAnswers bool
}

// UserB represents User B's answers in the database.
type UserB struct {
	gorm.Model
	Token        string `gorm:"uniqueIndex;size:255"`
	Answers      string `gorm:"type:text"` // JSON string of answers
	ShareAnswers bool
}

// Session represents a Q&A session between two users.
type Session struct {
	gorm.Model
	Token         string `gorm:"uniqueIndex;size:255"`
	UserAID       uint
	UserA         UserA
	UserBID       *uint // Nullable since UserB might not have submitted yet
	UserB         *UserB
	Compatibility int    // Compatibility score (0-100)
	Summary       string `gorm:"type:text"` // AI-generated summary
}

// Question represents a question in the Q&A application.
type Question struct {
	gorm.Model
	ID               uint   `json:"id" gorm:"primaryKey"`
	QuestionText     string `json:"question" gorm:"type:text"`
	IsMultipleChoice bool   `json:"isMultipleChoice"`
	Options          string `json:"options" gorm:"type:text"` // JSON string of options
}
