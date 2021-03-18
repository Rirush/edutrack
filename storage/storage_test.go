package storage

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

const TestdataPath = "testdata"

func TestStorage_EnsureCleanEnv(t *testing.T) {
	err := os.RemoveAll(TestdataPath)
	if err != nil {
		t.Errorf("os.RemoveAll() error = %v", err)
	}
	return
}

func TestNewStorage(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "create storage", args: args{rootPath: TestdataPath}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStorage(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_NewSubject(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "create subject", args: args{name: "Test Subject 1", description: "Period 1, Semester 1"}, wantErr: false},
		{name: "create secondary subject", args: args{name: "Test Subject 2", description: "Period 1, Semester 1"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = s.NewSubject(tt.args.name, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_ListSubjects(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	tests := []struct {
		name   string
		wantLen int
	}{
		{name: "check test data", wantLen: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.ListSubjects(); len(got) != tt.wantLen {
				t.Errorf("ListSubjects() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestStorage_SubjectExists(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{name: "check real", args: args{id: s.ListSubjects()[0].ID}, want: true},
		{name: "check fake", args: args{id: uuid.UUID{}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.SubjectExists(tt.args.id); got != tt.want {
				t.Errorf("SubjectExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_RemoveSubject(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "remove real subject", args: args{id: s.ListSubjects()[0].ID}, wantErr: false},
		{name: "remove fake subject", args: args{id: uuid.UUID{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.RemoveSubject(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveSubject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_UpdateSubject(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		id   uuid.UUID
		data Subject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "update existing subject", args: args{id: s.ListSubjects()[0].ID, data: Subject{Name: "Updated Subject", Description: "Period 1, Semester 2"}}, wantErr: false},
		{name: "update fake subject", args: args{id: uuid.UUID{}, data: Subject{Name: "Non-existing subject", Description: "Doesn't exist!!!"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpdateSubject(tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_NewEntry(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}
	
	type args struct {
		subjectID uuid.UUID
		title     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "create new entry", args: args{subjectID: s.ListSubjects()[0].ID, title: "Test entry 1"}, wantErr: false},
		{name: "create another entry", args: args{subjectID: s.ListSubjects()[0].ID, title: "Test entry 2"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.NewEntry(tt.args.subjectID, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_ListEntries(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		subjectID uuid.UUID
	}
	tests := []struct {
		name   string
		args   args
		want   int
	}{
		{name: "list entries", args: args{subjectID: s.ListSubjects()[0].ID}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.ListEntries(tt.args.subjectID); len(got) != tt.want {
				t.Errorf("ListEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_EntryExists(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		subjectID uuid.UUID
		id        uuid.UUID
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{name: "check valid entry", args: args{subjectID: s.ListSubjects()[0].ID, id: s.ListEntries(s.ListSubjects()[0].ID)[0].ID}, want: true},
		{name: "check invalid entry", args: args{subjectID: s.ListSubjects()[0].ID, id: uuid.UUID{}}, want: false},
		{name: "check invalid subject", args: args{subjectID: uuid.UUID{}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.EntryExists(tt.args.subjectID, tt.args.id); got != tt.want {
				t.Errorf("EntryExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_UpdateEntryBody(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		subjectID uuid.UUID
		id        uuid.UUID
		data      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "update existing entry", args: args{subjectID: s.ListSubjects()[0].ID, id: s.ListEntries(s.ListSubjects()[0].ID)[0].ID, data: "# Hello!"}, wantErr: false},
		{name: "update non-existing entry", args: args{subjectID: s.ListSubjects()[0].ID, id: uuid.UUID{}, data: "# Nope!"}, wantErr: true},
		{name: "update non-existing entry in non-existing subject", args: args{subjectID: uuid.UUID{}, id: uuid.UUID{}, data: "absolutely 100% nope"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpdateEntryBody(tt.args.subjectID, tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEntryBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_UpdateEntryMetadata(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}

	type args struct {
		subjectID uuid.UUID
		id        uuid.UUID
		entry     Entry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "update existing entry", args: args{subjectID: s.ListSubjects()[0].ID, id: s.ListEntries(s.ListSubjects()[0].ID)[0].ID, entry: Entry{Title: "updated title!!!"}}, wantErr: false},
		{name: "update fake entry", args: args{subjectID: s.ListSubjects()[0].ID, id: uuid.UUID{}, entry: Entry{Title: "no way"}}, wantErr: true},
		{name: "update fake entry in fake subject", args: args{subjectID: uuid.UUID{}, id: uuid.UUID{}, entry: Entry{Title: "nope"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.UpdateEntryMetadata(tt.args.subjectID, tt.args.id, tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEntryMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteEntry(t *testing.T) {
	s, err := NewStorage(TestdataPath)
	if err != nil {
		t.Errorf("NewStorage() error = %v", err)
		return
	}
	
	type args struct {
		subjectID uuid.UUID
		id        uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "delete existing entry", args: args{subjectID: s.ListSubjects()[0].ID, id: s.ListEntries(s.ListSubjects()[0].ID)[0].ID}, wantErr: false},
		{name: "delete fake entry", args: args{subjectID: s.ListSubjects()[0].ID, id: uuid.UUID{}}, wantErr: true},
		{name: "delete entry in fake subject", args: args{subjectID: uuid.UUID{}, id: uuid.UUID{}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeleteEntry(tt.args.subjectID, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}