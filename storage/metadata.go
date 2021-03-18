package storage

import (
	"github.com/google/uuid"
	"time"
)

type Subject struct {
	ID uuid.UUID
	Name string
	Description string
}

type SubjectsMetadata struct {
	Subjects map[uuid.UUID]Subject
}

type Entry struct {
	ID uuid.UUID
	Title string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntriesMetadata struct {
	Entries map[uuid.UUID]map[uuid.UUID]Entry
}