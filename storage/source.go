package storage

import (
	"errors"
	"time"
)

var (
	errEntryNotFound = errors.New("entry not found")
)

// Source represents a set of records on disk
type Source struct {
	SourceFile string
	Records    []*Record
	SyncedAt   *time.Time
}

// Size returns the number of records it contains
func (s *Source) Size() int {
	return len(s.Records)
}

// WasSynced indicates whether a source has been synced
func (s *Source) WasSynced() bool {
	return s.SyncedAt != nil
}

// FirstRecord returns the first record
func (s *Source) FirstRecord() *Record {
	return s.Records[0]
}

// RemoveEntry finds and removes an entry from a source
func (s *Source) RemoveEntry(e *Entry) error {
	for _, r := range s.Records {
		for i, re := range r.Entries {
			if re.URL == e.URL {
				r.Entries = append(r.Entries[:i], r.Entries[(i+1):]...)
				return nil
			}
		}
	}
	return errEntryNotFound
}
