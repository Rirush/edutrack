package storage

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

var (
	ErrNotDirectory = errors.New("storage path is not a directory")
	ErrNoSubject = errors.New("subject does not exist")
	ErrNoEntry = errors.New("entry does not exist")
)

type Storage struct {
	storageRoot string

	metadataMutex *sync.RWMutex
	subjectsMetadata *SubjectsMetadata
	entriesMetadata *EntriesMetadata
}

func NewStorage(rootPath string) (st *Storage, err error) {
	st = &Storage{
		storageRoot: rootPath,
		metadataMutex: &sync.RWMutex{},
	}

	info, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rootPath, 0755)
		if err != nil {
			return
		}
		err = os.MkdirAll(path.Join(rootPath, "entries"), 0755)
		if err != nil {
			return
		}
		err = os.MkdirAll(path.Join(rootPath, "metadata"), 0755)
		if err != nil {
			return
		}
		info, err = os.Stat(rootPath)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}

	if !info.IsDir() {
		err = ErrNotDirectory
		return
	}

	err = st.reloadMetadata()
	if err != nil {
		return
	}
	err = st.flushMetadata()
	return
}

func (s *Storage) NewSubject(name, description string) (uuid.UUID, error) {
	id := uuid.New()
	s.metadataMutex.Lock()
	s.subjectsMetadata.Subjects[id] = Subject{
		ID: id,
		Name: name,
		Description: description,
	}
	s.entriesMetadata.Entries[id] = make(map[uuid.UUID]Entry)
	s.metadataMutex.Unlock()
	err := s.flushMetadata()
	return id, err
}

func (s *Storage) ListSubjects() []Subject {
	s.metadataMutex.RLock()
	defer s.metadataMutex.RUnlock()
	subjects := make([]Subject, 0, len(s.subjectsMetadata.Subjects))
	for _, v := range s.subjectsMetadata.Subjects {
		subjects = append(subjects, v)
	}
	return subjects
}

func (s *Storage) SubjectExists(id uuid.UUID) bool {
	s.metadataMutex.RLock()
	defer s.metadataMutex.RUnlock()
	_, ok := s.subjectsMetadata.Subjects[id]
	return ok
}

func (s *Storage) RemoveSubject(id uuid.UUID) error {
	if !s.SubjectExists(id) {
		return ErrNoSubject
	}
	s.metadataMutex.Lock()
	delete(s.subjectsMetadata.Subjects, id)
	delete(s.entriesMetadata.Entries, id)
	err := os.RemoveAll(path.Join(s.storageRoot, "entries", id.String()))
	if err != nil {
		s.metadataMutex.Unlock()
		return err
	}
	s.metadataMutex.Unlock()
	return s.flushMetadata()
}

func (s *Storage) UpdateSubject(id uuid.UUID, data Subject) error {
	if !s.SubjectExists(id) {
		return ErrNoSubject
	}
	s.metadataMutex.Lock()
	data.ID = id
	s.subjectsMetadata.Subjects[id] = data
	s.metadataMutex.Unlock()
	return s.flushMetadata()
}

func (s *Storage) NewEntry(subjectID uuid.UUID, title string) (uuid.UUID, error) {
	if !s.SubjectExists(subjectID) {
		return uuid.UUID{}, ErrNoSubject
	}
	id := uuid.New()
	s.metadataMutex.Lock()
	s.entriesMetadata.Entries[subjectID][id] = Entry{
		ID: id,
		Title: title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.metadataMutex.Unlock()
	entryPath := path.Join(s.storageRoot, "entries", subjectID.String(), id.String())
	err := os.MkdirAll(entryPath, 0755)
	if err != nil {
		return id, err
	}
	f, err := os.OpenFile(path.Join(entryPath, "CONTENT.md"), os.O_CREATE, 0644)
	if err != nil {
		return id, err
	}
	f.Close()
	return id, s.flushMetadata()
}

func (s *Storage) ListEntries(subjectID uuid.UUID) []Entry {
	if !s.SubjectExists(subjectID) {
		return nil
	}
	s.metadataMutex.RLock()
	defer s.metadataMutex.RUnlock()
	entries := make([]Entry, 0, len(s.entriesMetadata.Entries[subjectID]))
	for _, v := range s.entriesMetadata.Entries[subjectID] {
		entries = append(entries, v)
	}
	return entries
}

func (s *Storage) EntryExists(subjectID, id uuid.UUID) bool {
	s.metadataMutex.RLock()
	defer s.metadataMutex.RUnlock()
	if !s.SubjectExists(subjectID) {
		return false
	}
	_, ok := s.entriesMetadata.Entries[subjectID][id]
	return ok
}

func (s *Storage) UpdateEntryBody(subjectID, id uuid.UUID, data string) error {
	if !s.SubjectExists(subjectID) {
		return ErrNoSubject
	}
	if !s.EntryExists(subjectID, id) {
		return ErrNoEntry
	}
	s.metadataMutex.Lock()
	entry := s.entriesMetadata.Entries[subjectID][id]
	entry.UpdatedAt = time.Now()
	s.entriesMetadata.Entries[subjectID][id] = entry
	s.metadataMutex.Unlock()
	err := s.flushMetadata()
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(s.storageRoot, "entries", subjectID.String(), id.String(), "CONTENT.md"), []byte(data), 0644)
}

func (s *Storage) UpdateEntryMetadata(subjectID, id uuid.UUID, entry Entry) error {
	if !s.SubjectExists(subjectID) {
		return ErrNoSubject
	}
	if !s.EntryExists(subjectID, id) {
		return ErrNoEntry
	}
	s.metadataMutex.Lock()
	oldEntry := s.entriesMetadata.Entries[subjectID][id]
	oldEntry.Title = entry.Title
	oldEntry.UpdatedAt = time.Now()
	s.entriesMetadata.Entries[subjectID][id] = oldEntry
	s.metadataMutex.Unlock()
	return s.flushMetadata()
}

func (s *Storage) DeleteEntry(subjectID, id uuid.UUID) error {
	if !s.SubjectExists(subjectID) {
		return ErrNoSubject
	}
	if !s.EntryExists(subjectID, id) {
		return ErrNoEntry
	}
	s.metadataMutex.Lock()
	delete(s.entriesMetadata.Entries[subjectID], id)
	s.metadataMutex.Unlock()
	err := s.flushMetadata()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(s.storageRoot, "entries", subjectID.String(), id.String()))
	return err
}

func (s *Storage) reloadMetadata() error {
	s.metadataMutex.Lock()
	defer s.metadataMutex.Unlock()

	err := s.readEntries()
	if err != nil {
		return err
	}
	return s.readSubjects()
}

func (s *Storage) readSubjects() error {
	subjectsPath := path.Join(s.storageRoot, "metadata", "subjects.json")
	f, err := os.Open(subjectsPath)
	switch {
	case os.IsNotExist(err):
		err = os.WriteFile(subjectsPath, []byte("{}"), 0644)
		if err != nil {
			return err
		}
		s.subjectsMetadata = &SubjectsMetadata{
			Subjects: make(map[uuid.UUID]Subject),
		}
	case err != nil:
		return err
	default:
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		s.subjectsMetadata = &SubjectsMetadata{}
		err = json.Unmarshal(data, s.subjectsMetadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) readEntries() error {
	entriesPath := path.Join(s.storageRoot, "metadata", "entries.json")
	f, err := os.Open(entriesPath)
	switch {
	case os.IsNotExist(err):
		err = os.WriteFile(entriesPath, []byte("{}"), 0644)
		if err != nil {
			return err
		}
		s.entriesMetadata = &EntriesMetadata{
			Entries: make(map[uuid.UUID]map[uuid.UUID]Entry),
		}
	case err != nil:
		return err
	default:
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		s.subjectsMetadata = &SubjectsMetadata{}
		err = json.Unmarshal(data, s.subjectsMetadata)
		if err != nil {
			return err
		}
		s.entriesMetadata = &EntriesMetadata{}
		err = json.Unmarshal(data, s.entriesMetadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) flushMetadata() error {
	s.metadataMutex.Lock()
	defer s.metadataMutex.Unlock()

	subjects, err := json.MarshalIndent(s.subjectsMetadata, "", " ")
	if err != nil {
		return err
	}
	subjectsPath := path.Join(s.storageRoot, "metadata", "subjects.json")
	err = os.WriteFile(subjectsPath, subjects, 0755)
	if err != nil {
		return err
	}

	entries, err := json.MarshalIndent(s.entriesMetadata, "", " ")
	if err != nil {
		return err
	}
	entriesPath := path.Join(s.storageRoot, "metadata", "entries.json")
	return os.WriteFile(entriesPath, entries, 0755)
}

